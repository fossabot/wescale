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

package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"vitess.io/vitess/go/cmd/vtctldclient/cli"
	"vitess.io/vitess/go/vt/topo/topoproto"

	vtctldatapb "vitess.io/vitess/go/vt/proto/vtctldata"
)

var (
	// GetShard makes a GetShard gRPC request to a vtctld.
	GetShard = &cobra.Command{
		Use:  "GetShard <keyspace/shard>",
		Args: cobra.ExactArgs(1),
		RunE: commandGetShard,
	}
	// RemoveShardCell makes a RemoveShardCell gRPC request to a vtctld.
	RemoveShardCell = &cobra.Command{
		Use:  "RemoveShardCell <keyspace/shard> <cell>",
		Args: cobra.ExactArgs(2),
		RunE: commandRemoveShardCell,
	}
)

func commandGetShard(cmd *cobra.Command, args []string) error {
	keyspace, shard, err := topoproto.ParseKeyspaceShard(cmd.Flags().Arg(0))
	if err != nil {
		return err
	}

	resp, err := client.GetShard(commandCtx, &vtctldatapb.GetShardRequest{
		Keyspace:  keyspace,
		ShardName: shard,
	})
	if err != nil {
		return err
	}

	data, err := cli.MarshalJSON(resp.Shard)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", data)

	return nil
}

var removeShardCellOptions = struct {
	Force     bool
	Recursive bool
}{}

func commandRemoveShardCell(cmd *cobra.Command, args []string) error {
	keyspace, shard, err := topoproto.ParseKeyspaceShard(cmd.Flags().Arg(0))
	if err != nil {
		return err
	}

	cell := cmd.Flags().Arg(1)

	_, err = client.RemoveShardCell(commandCtx, &vtctldatapb.RemoveShardCellRequest{
		Keyspace:  keyspace,
		ShardName: shard,
		Cell:      cell,
		Force:     removeShardCellOptions.Force,
		Recursive: removeShardCellOptions.Recursive,
	})

	if err != nil {
		return err
	}

	fmt.Printf("Successfully removed cell %v from shard %s/%s\n", cell, keyspace, shard)

	return nil
}

func init() {
	Root.AddCommand(GetShard)

	RemoveShardCell.Flags().BoolVarP(&removeShardCellOptions.Force, "force", "f", false, "Proceed even if the cell's topology server cannot be reached. The assumption is that you turned down the entire cell, and just need to update the global topo data.")
	RemoveShardCell.Flags().BoolVarP(&removeShardCellOptions.Recursive, "recursive", "r", false, "Also delete all tablets in that cell beloning to the specified shard.")
	Root.AddCommand(RemoveShardCell)
}
