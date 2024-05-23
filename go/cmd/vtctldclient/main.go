/*
Copyright 2020 The Vitess Authors.

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

package main

import (
	"flag"

	"github.com/wesql/wescale/go/acl"
	"github.com/wesql/wescale/go/cmd/vtctldclient/command"
	"github.com/wesql/wescale/go/exit"
	"github.com/wesql/wescale/go/vt/grpcclient"
	"github.com/wesql/wescale/go/vt/grpccommon"
	"github.com/wesql/wescale/go/vt/log"
	"github.com/wesql/wescale/go/vt/logutil"
	"github.com/wesql/wescale/go/vt/servenv"
	"github.com/wesql/wescale/go/vt/vtctl/grpcclientcommon"
	"github.com/wesql/wescale/go/vt/vtctl/vtctlclient"

	_flag "github.com/wesql/wescale/go/internal/flag"
)

func main() {
	defer exit.Recover()

	// Grab all those global flags across the codebase and shove 'em on in.
	// (TODO|andrew) remove this line after the migration to pflag is complete.
	command.Root.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	log.RegisterFlags(command.Root.PersistentFlags())
	logutil.RegisterFlags(command.Root.PersistentFlags())
	grpcclient.RegisterFlags(command.Root.PersistentFlags())
	grpccommon.RegisterFlags(command.Root.PersistentFlags())
	grpcclientcommon.RegisterFlags(command.Root.PersistentFlags())
	servenv.RegisterMySQLServerFlags(command.Root.PersistentFlags())
	vtctlclient.RegisterFlags(command.Root.PersistentFlags())
	acl.RegisterFlags(command.Root.PersistentFlags())

	// hack to get rid of an "ERROR: logging before flag.Parse"
	_flag.TrickGlog()

	// back to your regularly scheduled cobra programming
	if err := command.Root.Execute(); err != nil {
		log.Error(err)
		exit.Return(1)
	}
}
