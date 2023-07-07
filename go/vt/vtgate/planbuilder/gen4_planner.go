/*
Copyright ApeCloud, Inc.
Licensed under the Apache v2(found in the LICENSE file in the root directory).
*/

/*
Copyright 2021 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package planbuilder

import (
	"fmt"

	"vitess.io/vitess/go/internal/global"

	querypb "vitess.io/vitess/go/vt/proto/query"
	"vitess.io/vitess/go/vt/sqlparser"
	"vitess.io/vitess/go/vt/vterrors"
	"vitess.io/vitess/go/vt/vtgate/engine"
	"vitess.io/vitess/go/vt/vtgate/planbuilder/operators"
	"vitess.io/vitess/go/vt/vtgate/planbuilder/plancontext"
	"vitess.io/vitess/go/vt/vtgate/semantics"
	"vitess.io/vitess/go/vt/vtgate/vindexes"
)

func gen4Planner(query string, plannerVersion querypb.ExecuteOptions_PlannerVersion) stmtPlanner {
	return func(stmt sqlparser.Statement, reservedVars *sqlparser.ReservedVars, vschema plancontext.VSchema) (*planResult, error) {
		switch stmt := stmt.(type) {
		case sqlparser.SelectStatement:
			return gen4SelectStmtPlanner(plannerVersion, stmt, reservedVars, vschema)
		case *sqlparser.Update:
			return gen4UpdateStmtPlanner(plannerVersion, stmt, reservedVars, vschema)
		case *sqlparser.Delete:
			return gen4DeleteStmtPlanner(plannerVersion, stmt, reservedVars, vschema)
		default:
			return nil, vterrors.VT12001(fmt.Sprintf("%T", stmt))
		}
	}
}

func gen4SelectStmtPlanner(
	plannerVersion querypb.ExecuteOptions_PlannerVersion,
	stmt sqlparser.SelectStatement,
	reservedVars *sqlparser.ReservedVars,
	vschema plancontext.VSchema,
) (*planResult, error) {
	plan, err := handleSelectLock(stmt, reservedVars, vschema)
	if err != nil {
		return nil, err
	}
	if plan != nil {
		return plan, nil
	}
	plan, err = buildPlanForBypass(stmt, reservedVars, vschema)
	if err != nil {
		return nil, err
	}
	return plan, nil
}

func gen4planSQLCalcFoundRows(vschema plancontext.VSchema, sel *sqlparser.Select, query string, reservedVars *sqlparser.ReservedVars) (*planResult, error) {
	ksName := ""
	if ks, _ := vschema.DefaultKeyspace(); ks != nil {
		ksName = ks.Name
	}
	semTable, err := semantics.Analyze(sel, ksName, vschema)
	if err != nil {
		return nil, err
	}
	// record any warning as planner warning.
	vschema.PlannerWarning(semTable.Warning)

	plan, tablesUsed, err := buildSQLCalcFoundRowsPlan(query, sel, reservedVars, vschema, planSelectGen4)
	if err != nil {
		return nil, err
	}
	return newPlanResult(plan.Primitive(), tablesUsed...), nil
}

func planSelectGen4(reservedVars *sqlparser.ReservedVars, vschema plancontext.VSchema, sel *sqlparser.Select) (*jointab, logicalPlan, []string, error) {
	plan, _, tablesUsed, err := newBuildSelectPlan(sel, reservedVars, vschema, 0)
	if err != nil {
		return nil, nil, nil, err
	}
	return nil, plan, tablesUsed, nil
}

func gen4PredicateRewrite(stmt sqlparser.Statement, getPlan func(selStatement sqlparser.SelectStatement) (logicalPlan, *semantics.SemTable, []string, error)) (logicalPlan, *semantics.SemTable, []string) {
	rewritten, isSel := sqlparser.RewritePredicate(stmt).(sqlparser.SelectStatement)
	if !isSel {
		// Fail-safe code, should never happen
		return nil, nil, nil
	}
	plan2, st, op, err := getPlan(rewritten)
	if err == nil && !shouldRetryAfterPredicateRewriting(plan2) {
		// we only use this new plan if it's better than the old one we got
		return plan2, st, op
	}
	return nil, nil, nil
}

func newBuildSelectPlan(
	selStmt sqlparser.SelectStatement,
	reservedVars *sqlparser.ReservedVars,
	vschema plancontext.VSchema,
	version querypb.ExecuteOptions_PlannerVersion,
) (plan logicalPlan, semTable *semantics.SemTable, tablesUsed []string, err error) {
	ksName := ""
	usingKs, err := GetUsingKs(selStmt, vschema)
	if err != nil {
		return nil, nil, nil, vterrors.VT09005()
	}
	ksName = usingKs.Name
	semTable, err = semantics.Analyze(selStmt, ksName, vschema)
	if err != nil {
		return nil, nil, nil, err
	}
	// record any warning as planner warning.
	vschema.PlannerWarning(semTable.Warning)

	ctx := plancontext.NewPlanningContext(reservedVars, semTable, vschema, version)

	plan, tablesUsed, err = pushdownShortcut(ctx, selStmt, usingKs)
	if err != nil {
		return nil, nil, nil, err
	}
	plan, err = pushCommentDirectivesOnPlan(plan, selStmt)
	if err != nil {
		return nil, nil, nil, err
	}
	return plan, semTable, tablesUsed, err
}

// GetUsingKs returns the keyspace to use for the query.
func GetUsingKs(node sqlparser.SQLNode, vschema plancontext.VSchema) (*vindexes.Keyspace, error) {
	// If the query has a USING clause, it returns the keyspace specified in the clause.
	usingKs, err := vschema.DefaultKeyspace()
	if err == nil {
		return usingKs, nil
	}

	// If the query doesn't lookup any tables, it returns the default keyspace.
	allTables := sqlparser.GetAllTableNames(node)
	if len(allTables) == 1 && allTables[0].Name.String() == "dual" {
		return vschema.FindKeyspace(global.DefaultKeyspace)
	}

	// Otherwise, it returns the default keyspace, the default keyspace is the first keyspace in the query.
	firstKsName := ""
	for _, table := range allTables {
		if table.Qualifier.IsEmpty() {
			continue
		}
		firstKsName = table.Qualifier.String()
	}
	if firstKsName == "" {
		return nil, vterrors.VT09005()
	}
	usingKs, err = vschema.FindKeyspace(firstKsName)
	if err != nil {
		return nil, err
	}
	return usingKs, nil
}

// optimizePlan removes unnecessary simpleProjections that have been created while planning
func optimizePlan(plan logicalPlan) logicalPlan {
	newPlan, _ := visit(plan, func(plan logicalPlan) (bool, logicalPlan, error) {
		this, ok := plan.(*simpleProjection)
		if !ok {
			return true, plan, nil
		}

		input, ok := this.input.(*simpleProjection)
		if !ok {
			return true, plan, nil
		}

		for i, col := range this.eSimpleProj.Cols {
			this.eSimpleProj.Cols[i] = input.eSimpleProj.Cols[col]
		}
		this.input = input.input
		return true, this, nil
	})
	return newPlan
}

func gen4UpdateStmtPlanner(
	version querypb.ExecuteOptions_PlannerVersion,
	updStmt *sqlparser.Update,
	reservedVars *sqlparser.ReservedVars,
	vschema plancontext.VSchema,
) (*planResult, error) {
	if updStmt.With != nil {
		return nil, vterrors.VT12001("WITH expression in UPDATE statement")
	}

	ksName := ""
	usingKs, err := GetUsingKs(updStmt, vschema)
	if err != nil {
		return nil, vterrors.VT09005()
	}
	ksName = usingKs.Name
	semTable, err := semantics.Analyze(updStmt, ksName, vschema)
	if err != nil {
		return nil, err
	}
	// record any warning as planner warning.
	vschema.PlannerWarning(semTable.Warning)

	err = rewriteRoutedTables(updStmt, vschema)
	if err != nil {
		return nil, err
	}

	tables := semTable.GetVindexTable()
	edml := engine.NewDML()
	edml.Keyspace = usingKs
	edml.Table = tables
	edml.Opcode = engine.Unsharded
	edml.Query = generateQuery(updStmt)
	upd := &engine.Update{DML: edml}
	return newPlanResult(upd, operators.QualifiedTables(usingKs, tables)...), nil
}

func gen4DeleteStmtPlanner(
	version querypb.ExecuteOptions_PlannerVersion,
	deleteStmt *sqlparser.Delete,
	reservedVars *sqlparser.ReservedVars,
	vschema plancontext.VSchema,
) (*planResult, error) {
	if deleteStmt.With != nil {
		return nil, vterrors.VT12001("WITH expression in DELETE statement")
	}

	var err error
	if len(deleteStmt.TableExprs) == 1 && len(deleteStmt.Targets) == 1 {
		deleteStmt, err = rewriteSingleTbl(deleteStmt)
		if err != nil {
			return nil, err
		}
	}

	ksName := ""
	usingKs, err := GetUsingKs(deleteStmt, vschema)
	if err != nil {
		return nil, vterrors.VT09005()
	}
	ksName = usingKs.Name
	semTable, err := semantics.Analyze(deleteStmt, ksName, vschema)
	if err != nil {
		return nil, err
	}

	// record any warning as planner warning.
	vschema.PlannerWarning(semTable.Warning)
	err = rewriteRoutedTables(deleteStmt, vschema)
	if err != nil {
		return nil, err
	}

	tables := semTable.GetVindexTable()
	edml := engine.NewDML()
	edml.Keyspace = usingKs
	edml.Table = tables
	edml.Opcode = engine.Unsharded
	edml.Query = generateQuery(deleteStmt)
	del := &engine.Delete{DML: edml}
	return newPlanResult(del, operators.QualifiedTables(usingKs, tables)...), nil
}

func rewriteRoutedTables(stmt sqlparser.Statement, vschema plancontext.VSchema) error {
	// Rewrite routed tables
	return sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
		aliasTbl, isAlias := node.(*sqlparser.AliasedTableExpr)
		if !isAlias {
			return true, nil
		}
		tableName, ok := aliasTbl.Expr.(sqlparser.TableName)
		if !ok {
			return true, nil
		}
		var vschemaTable *vindexes.Table
		vschemaTable, _, _, _, _, err = vschema.FindTableOrVindex(tableName)
		if err != nil {
			return false, err
		}

		if vschemaTable.Name.String() != tableName.Name.String() {
			name := tableName.Name
			if aliasTbl.As.IsEmpty() {
				// if the user hasn't specified an alias, we'll insert one here so the old table name still works
				aliasTbl.As = sqlparser.NewIdentifierCS(name.String())
			}
			tableName.Name = sqlparser.NewIdentifierCS(vschemaTable.Name.String())
			aliasTbl.Expr = tableName
		}

		return true, nil
	}, stmt)
}

func setLockOnAllSelect(plan logicalPlan) {
	_, _ = visit(plan, func(plan logicalPlan) (bool, logicalPlan, error) {
		switch node := plan.(type) {
		case *routeGen4:
			node.Select.SetLock(sqlparser.ShareModeLock)
			return true, node, nil
		}
		return true, plan, nil
	})
}

func planLimit(limit *sqlparser.Limit, plan logicalPlan) (logicalPlan, error) {
	if limit == nil {
		return plan, nil
	}
	rb, ok := plan.(*routeGen4)
	if ok && rb.isSingleShard() {
		rb.SetLimit(limit)
		return plan, nil
	}

	lPlan, err := createLimit(plan, limit)
	if err != nil {
		return nil, err
	}

	// visit does not modify the plan.
	_, err = visit(lPlan, setUpperLimit)
	if err != nil {
		return nil, err
	}
	return lPlan, nil
}

func planHorizon(ctx *plancontext.PlanningContext, plan logicalPlan, in sqlparser.SelectStatement, truncateColumns bool) (logicalPlan, error) {
	switch node := in.(type) {
	case *sqlparser.Select:
		hp := horizonPlanning{
			sel: node,
		}

		replaceSubQuery(ctx, node)
		var err error
		plan, err = hp.planHorizon(ctx, plan, truncateColumns)
		if err != nil {
			return nil, err
		}
		plan, err = planLimit(node.Limit, plan)
		if err != nil {
			return nil, err
		}
	case *sqlparser.Union:
		var err error
		rb, isRoute := plan.(*routeGen4)
		if !isRoute && ctx.SemTable.NotSingleRouteErr != nil {
			return nil, ctx.SemTable.NotSingleRouteErr
		}
		if isRoute && rb.isSingleShard() {
			err = planSingleShardRoutePlan(node, rb)
		} else {
			plan, err = planOrderByOnUnion(ctx, plan, node)
		}
		if err != nil {
			return nil, err
		}

		plan, err = planLimit(node.Limit, plan)
		if err != nil {
			return nil, err
		}
	}
	return plan, nil

}

func planOrderByOnUnion(ctx *plancontext.PlanningContext, plan logicalPlan, union *sqlparser.Union) (logicalPlan, error) {
	qp, err := operators.CreateQPFromUnion(union)
	if err != nil {
		return nil, err
	}
	hp := horizonPlanning{
		qp: qp,
	}
	if len(qp.OrderExprs) > 0 {
		plan, err = hp.planOrderBy(ctx, qp.OrderExprs, plan)
		if err != nil {
			return nil, err
		}
	}
	return plan, nil
}

func pushCommentDirectivesOnPlan(plan logicalPlan, stmt sqlparser.Statement) (logicalPlan, error) {
	var directives *sqlparser.CommentDirectives
	cmt, ok := stmt.(sqlparser.Commented)
	if ok {
		directives = cmt.GetParsedComments().Directives()
		scatterAsWarns := directives.IsSet(sqlparser.DirectiveScatterErrorsAsWarnings)
		timeout := queryTimeout(directives)

		if scatterAsWarns || timeout > 0 {
			_, _ = visit(plan, func(logicalPlan logicalPlan) (bool, logicalPlan, error) {
				switch plan := logicalPlan.(type) {
				case *routeGen4:
					plan.eroute.ScatterErrorsAsWarnings = scatterAsWarns
					plan.eroute.QueryTimeout = timeout
				}
				return true, logicalPlan, nil
			})
		}
	}

	return plan, nil
}

// checkIfDeleteSupported checks if the delete query is supported or we must return an error.
func checkIfDeleteSupported(del *sqlparser.Delete, semTable *semantics.SemTable) error {
	if semTable.NotUnshardedErr != nil {
		return semTable.NotUnshardedErr
	}

	// Delete is only supported for a single TableExpr which is supposed to be an aliased expression
	multiShardErr := vterrors.VT12001("multi-shard or vindex write statement")
	if len(del.TableExprs) != 1 {
		return multiShardErr
	}
	_, isAliasedExpr := del.TableExprs[0].(*sqlparser.AliasedTableExpr)
	if !isAliasedExpr {
		return multiShardErr
	}

	if len(del.Targets) > 1 {
		return vterrors.VT12001("multi-table DELETE statement in a sharded keyspace")
	}

	err := sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
		switch node.(type) {
		case *sqlparser.Subquery, *sqlparser.DerivedTable:
			// We have a subquery, so we must fail the planning.
			// If this subquery and the table expression were all belonging to the same unsharded keyspace,
			// we would have already created a plan for them before doing these checks.
			return false, vterrors.VT12001("subqueries in DML")
		}
		return true, nil
	}, del)
	if err != nil {
		return err
	}

	return nil
}
