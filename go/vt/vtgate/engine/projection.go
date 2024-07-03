/*
Copyright 2022 The Vitess Authors.

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

package engine

import (
	"context"
	"errors"
	"strings"
	"sync"
	"unicode"
	"vitess.io/vitess/go/sqltypes"
	querypb "vitess.io/vitess/go/vt/proto/query"
	"vitess.io/vitess/go/vt/sqlparser"
	"vitess.io/vitess/go/vt/vtgate/evalengine"
)

var _ Primitive = (*Projection)(nil)

// Projection can evaluate expressions and project the results
type Projection struct {
	Cols  []string
	Exprs []evalengine.Expr
	Input Primitive
	noTxNeeded
	IsCustomFunctionProjection bool
	Meta                       *CustomFunctionProjectionMeta
}

type CustomFunctionProjectionMeta struct {
	Origin sqlparser.SelectExprs
}

// RouteType implements the Primitive interface
func (p *Projection) RouteType() string {
	return p.Input.RouteType()
}

// GetKeyspaceName implements the Primitive interface
func (p *Projection) GetKeyspaceName() string {
	return p.Input.GetKeyspaceName()
}

// GetTableName implements the Primitive interface
func (p *Projection) GetTableName() string {
	return p.Input.GetTableName()
}

// TryExecute implements the Primitive interface
func (p *Projection) TryExecute(ctx context.Context, vcursor VCursor, bindVars map[string]*querypb.BindVariable, wantfields bool) (*sqltypes.Result, error) {
	if p.IsCustomFunctionProjection {
		return p.executeCustomFunctionProjection(ctx, vcursor, bindVars, wantfields)
	}

	result, err := vcursor.ExecutePrimitive(ctx, p.Input, bindVars, wantfields)
	if err != nil {
		return nil, err
	}

	env := evalengine.EnvWithBindVars(bindVars, vcursor.ConnCollation())
	env.Fields = result.Fields
	var resultRows []sqltypes.Row
	for _, row := range result.Rows {
		resultRow := make(sqltypes.Row, 0, len(p.Exprs))
		env.Row = row
		for _, exp := range p.Exprs {
			result, err := env.Evaluate(exp)
			if err != nil {
				return nil, err
			}
			resultRow = append(resultRow, result.Value())
		}
		resultRows = append(resultRows, resultRow)
	}
	if wantfields {
		err := p.addFields(env, result)
		if err != nil {
			return nil, err
		}
	}
	result.Rows = resultRows
	return result, nil
}

func compareStrings(s1, s2 string) bool {
	normalized1 := normalizeString(s1)
	normalized2 := normalizeString(s2)
	return normalized1 == normalized2
}

func normalizeString(s string) string {
	var builder strings.Builder
	for _, r := range s {
		if !unicode.IsSpace(r) {
			builder.WriteRune(unicode.ToLower(r))
		}
	}
	return builder.String()
}

// todo newborn22 是否在这里对函数支持的类型做检查？
func GetColNameFromFuncExpr(funcExpr *sqlparser.FuncExpr) ([]sqlparser.SelectExpr, error) {
	rst := make([]sqlparser.SelectExpr, 0)
	for _, expr := range funcExpr.Exprs {
		switch expr.(type) {
		case *sqlparser.AliasedExpr:
			alias, _ := expr.(*sqlparser.AliasedExpr)
			_, ok := alias.Expr.(*sqlparser.ColName)
			if ok {
				rst = append(rst, alias)
			} else {
				subFunc, ok := alias.Expr.(*sqlparser.FuncExpr)
				if ok {
					subExpr, err := GetColNameFromFuncExpr(subFunc)
					if err != nil {
						return nil, err
					}
					rst = append(rst, subExpr...)
				}
			}

		default:
			return nil, errors.New("not support")
		}
	}
	return rst, nil
}

func (p *Projection) executeCustomFunctionProjection(ctx context.Context, vcursor VCursor, bindVars map[string]*querypb.BindVariable, wantfields bool) (*sqltypes.Result, error) {
	qr, err := vcursor.ExecutePrimitive(ctx, p.Input, bindVars, wantfields)
	if err != nil {
		return nil, err
	}

	newFieldNames := make([]string, 0)
	idx := 0
	// build final field
	colNumsForStar := 0
	for _, expr := range p.Meta.Origin {
		if alias, ok := expr.(*sqlparser.AliasedExpr); ok {
			if funcExpr, ok := alias.Expr.(*sqlparser.FuncExpr); ok {
				newFieldNames = append(newFieldNames, sqlparser.String(funcExpr))

				colNames, err := GetColNameFromFuncExpr(funcExpr)
				if err != nil {
					return nil, err
				}
				idx += len(colNames)
				continue
			}
		}

		// todo newborn22, 这样子直接拿名字ok?
		colNameInOrigin := sqlparser.String(expr)
		if colNameInOrigin == "*" {
			set := make(map[string]bool)
			for {
				if _, exist := set[qr.Fields[idx].Name]; exist {
					break
				}
				set[qr.Fields[idx].Name] = true
				newFieldNames = append(newFieldNames, qr.Fields[idx].Name)
				idx++
			}
			colNumsForStar = len(newFieldNames)
		} else {
			newFieldNames = append(newFieldNames, qr.Fields[idx].Name)
			idx++
		}
	}

	// build final result
	rows := [][]sqltypes.Value{}
	for _, gotRow := range qr.Named().Rows {
		rowValues := make([]string, 0, len(newFieldNames))
		idx := 0
		for _, colExpr := range p.Meta.Origin {
			if alias, ok := colExpr.(*sqlparser.AliasedExpr); ok {
				if funcExpr, ok := alias.Expr.(*sqlparser.FuncExpr); ok {
					funcRst, err := CalFuncExpr(funcExpr, gotRow)
					if err != nil {
						return nil, err
					}
					rowValues = append(rowValues, funcRst)
					idx++
					continue
				}
			}
			if _, ok := colExpr.(*sqlparser.StarExpr); ok {
				for i := 0; i < colNumsForStar; i++ {
					rowValues = append(rowValues, gotRow[newFieldNames[idx]].ToString())
					idx++
				}
				continue
			}
			rowValues = append(rowValues, gotRow[newFieldNames[idx]].ToString())
			idx++
		}
		rows = append(rows, BuildVarCharRow(rowValues...))
	}

	return &sqltypes.Result{
		Fields: BuildVarCharFields(newFieldNames...),
		Rows:   rows,
	}, nil
	//return qr, nil
}

// TryStreamExecute implements the Primitive interface
func (p *Projection) TryStreamExecute(ctx context.Context, vcursor VCursor, bindVars map[string]*querypb.BindVariable, wantfields bool, callback func(*sqltypes.Result) error) error {
	env := evalengine.EnvWithBindVars(bindVars, vcursor.ConnCollation())
	var once sync.Once
	var fields []*querypb.Field
	return vcursor.StreamExecutePrimitive(ctx, p.Input, bindVars, wantfields, func(qr *sqltypes.Result) error {
		var err error
		if wantfields {
			once.Do(func() {
				env.Fields = qr.Fields
				fieldRes := &sqltypes.Result{}
				err = p.addFields(env, fieldRes)
				if err != nil {
					return
				}
				fields = fieldRes.Fields
				err = callback(fieldRes)
				if err != nil {
					return
				}
			})
			qr.Fields = fields
		}
		if err != nil {
			return err
		}
		resultRows := make([]sqltypes.Row, 0, len(qr.Rows))
		for _, r := range qr.Rows {
			resultRow := make(sqltypes.Row, 0, len(p.Exprs))
			env.Row = r
			for _, exp := range p.Exprs {
				c, err := env.Evaluate(exp)
				if err != nil {
					return err
				}
				resultRow = append(resultRow, c.Value())
			}
			resultRows = append(resultRows, resultRow)
		}
		qr.Rows = resultRows
		return callback(qr)
	})
}

// GetFields implements the Primitive interface
func (p *Projection) GetFields(ctx context.Context, vcursor VCursor, bindVars map[string]*querypb.BindVariable) (*sqltypes.Result, error) {
	qr, err := p.Input.GetFields(ctx, vcursor, bindVars)
	if err != nil {
		return nil, err
	}
	env := evalengine.EnvWithBindVars(bindVars, vcursor.ConnCollation())
	err = p.addFields(env, qr)
	if err != nil {
		return nil, err
	}
	return qr, nil
}

func (p *Projection) addFields(env *evalengine.ExpressionEnv, qr *sqltypes.Result) error {
	qr.Fields = nil
	for i, col := range p.Cols {
		q, err := env.TypeOf(p.Exprs[i])
		if err != nil {
			return err
		}
		qr.Fields = append(qr.Fields, &querypb.Field{
			Name: col,
			Type: q,
		})
	}
	return nil
}

// Inputs implements the Primitive interface
func (p *Projection) Inputs() []Primitive {
	return []Primitive{p.Input}
}

// description implements the Primitive interface
func (p *Projection) description() PrimitiveDescription {
	var exprs []string
	for idx, e := range p.Exprs {
		expr := evalengine.FormatExpr(e)
		alias := p.Cols[idx]
		if alias != "" {
			expr += " as " + alias
		}
		exprs = append(exprs, expr)
	}
	return PrimitiveDescription{
		OperatorType: "Projection",
		Other: map[string]any{
			"Expressions": exprs,
		},
	}
}

func CalFuncExpr(funcExpr *sqlparser.FuncExpr, rowValues sqltypes.RowNamedValues) (string, error) {
	// get function paras
	params := make([]string, 0, len(funcExpr.Exprs))
	// todo newborn22， 简单地支持了 literal, colname, funcExpr作为参数
	for _, para := range funcExpr.Exprs {
		alias, ok := para.(*sqlparser.AliasedExpr)
		if !ok {
			return "", errors.New("only support literal, colname and funcExpr as parameter")
		}
		switch alias.Expr.(type) {
		case *sqlparser.ColName:
			colName := alias.Expr.(*sqlparser.ColName).Name.String()
			val := rowValues[colName].ToString()
			params = append(params, val)
		case *sqlparser.Literal:
			val := sqlparser.String(alias.Expr.(*sqlparser.Literal))
			params = append(params, val)
		case *sqlparser.FuncExpr:
			// todo newborn22, 递归调用
			rst, err := CalFuncExpr(alias.Expr.(*sqlparser.FuncExpr), rowValues)
			if err != nil {
				return "", err
			}
			params = append(params, rst)
		default:
			return "", errors.New("only support literal, colname and funcExpr as parameter")
		}
	}
	// get the function
	function, _ := CUSTOM_FUNCTIONS[funcExpr.Name.String()]
	return function(params)
}
