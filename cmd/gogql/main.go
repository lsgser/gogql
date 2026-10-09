/*
|--------------------------------------------------------------------------
| gogql CLI
|--------------------------------------------------------------------------
|
| Command-line entrypoint: gogql init (scaffold project) and gogql version.
|
*/

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lsgser/gogql/internal/scaffold"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		runInit(os.Args[2:])
	case "version", "-v", "--version":
		fmt.Println("gogql", version)
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n", os.Args[1])
		usage()
		os.Exit(1)
	}
}

func runInit(args []string) {
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	module := fs.String("module", "", "Go module path (default: directory name)")
	gogqlPath := fs.String("gogql", "../gogql", "local replace path for github.com/lsgser/gogql (use -gogql= for published module only)")
	_ = fs.Parse(args)

	dir := "."
	if fs.NArg() > 0 {
		dir = fs.Arg(0)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "mkdir: %v\n", err)
		os.Exit(1)
	}
	abs, _ := filepath.Abs(dir)
	opts := scaffold.InitOptions{Dir: abs, ModulePath: *module}
	if *gogqlPath != "" {
		opts.GogqlReplace = *gogqlPath
	}
	if err := scaffold.Init(opts); err != nil {
		fmt.Fprintf(os.Stderr, "init: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created gogql project in %s\n", abs)
	fmt.Println("Next:")
	fmt.Println("  cd", dir)
	fmt.Println("  go run .")
}

func usage() {
	fmt.Fprintf(os.Stderr, `gogql — modular GraphQL server toolkit for Go

Usage:
  gogql init [directory]   scaffold a new server project
  gogql version            print version

Flags for init:
  -module string   Go module import path
  -gogql string    local replace path for gogql (default ../gogql)

`)
}
