/*
Copyright 2019 The Vitess Authors.

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

// You can modify this file to hook up a different logging library instead of glog.
// If you adapt to a different logging framework, you may need to use that
// framework's equivalent of *Depth() functions so the file and line number printed
// point to the real caller instead of your adapter function.

package log

import (
	"github.com/golang/glog"
	"github.com/spf13/pflag"
)

// Level is used with V() to test log verbosity.
type Level = glog.Level

const DebugVerbosity = 5
const TraceVerbosity = 10

var (
	// V quickly checks if the logging verbosity meets a threshold.
	V = glog.V

	// Flush ensures any pending I/O is written.
	Flush = glog.Flush

	Debug      = glog.V(DebugVerbosity).Info
	Debugf     = glog.V(DebugVerbosity).Infof
	DebugDepth = glog.V(DebugVerbosity).InfoDepth

	Trace      = glog.V(TraceVerbosity).Info
	Tracef     = glog.V(TraceVerbosity).Infof
	TraceDepth = glog.V(TraceVerbosity).InfoDepth

	// Info formats arguments like fmt.Print.
	Info = glog.Info
	// Infof formats arguments like fmt.Printf.
	Infof = glog.Infof
	// InfoDepth formats arguments like fmt.Print and uses depth to choose which call frame to log.
	InfoDepth = glog.InfoDepth

	// Warning formats arguments like fmt.Print.
	Warning = glog.Warning
	// Warningf formats arguments like fmt.Printf.
	Warningf = glog.Warningf
	// WarningDepth formats arguments like fmt.Print and uses depth to choose which call frame to log.
	WarningDepth = glog.WarningDepth

	// Error formats arguments like fmt.Print.
	Error = glog.Error
	// Errorf formats arguments like fmt.Printf.
	Errorf = glog.Errorf
	// ErrorDepth formats arguments like fmt.Print and uses depth to choose which call frame to log.
	ErrorDepth = glog.ErrorDepth

	// Exit formats arguments like fmt.Print.
	Exit = glog.Exit
	// Exitf formats arguments like fmt.Printf.
	Exitf = glog.Exitf
	// ExitDepth formats arguments like fmt.Print and uses depth to choose which call frame to log.
	ExitDepth = glog.ExitDepth

	// Fatal formats arguments like fmt.Print.
	Fatal = glog.Fatal
	// Fatalf formats arguments like fmt.Printf
	Fatalf = glog.Fatalf
	// FatalDepth formats arguments like fmt.Print and uses depth to choose which call frame to log.
	FatalDepth = glog.FatalDepth
)

// RegisterFlags installs log flags on the given FlagSet.
//
// `go/cmd/*` entrypoints should either use servenv.ParseFlags(WithArgs)? which
// calls this function, or call this function directly before parsing
// command-line arguments.
func RegisterFlags(fs *pflag.FlagSet) {
	fs.Uint64Var(&glog.MaxSize, "log_rotate_max_size", glog.MaxSize, "size in bytes at which logs are rotated (glog.MaxSize)")
}
