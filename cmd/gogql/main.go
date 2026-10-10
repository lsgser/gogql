/*
|--------------------------------------------------------------------------
| gogql CLI
|--------------------------------------------------------------------------
|
| Entry point for the gogql command-line tool. Dispatches init (full project
| scaffold), module subcommands (add, typedefs, resolvers, schema), version,
| and help. Parses flags per subcommand and delegates file generation to
| internal/scaffold. Install via go install github.com/lsgser/gogql/cmd/gogql
| or pin with go get -tool in your app's go.mod.
|
| Subcommands are documented in internal/cli.Commands.
|
*/

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lsgser/gogql/internal/cli"
	"github.com/lsgser/gogql/internal/scaffold"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		runInit(os.Args[2:])
	case "module":
		runModule(os.Args[2:])
	case "version", "-v", "--version":
		printVersion()
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printVersion() {
	fmt.Println("gogql", cli.Version)
	fmt.Println()
	fmt.Println("Commands:")
	for _, c := range cli.Commands {
		fmt.Printf("  %-18s %s\n", c.Name, c.Summary)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "gogql — modular GraphQL server toolkit for Go (v%s)\n\n", cli.Version)
	fmt.Fprintf(os.Stderr, "Usage:\n")
	for _, c := range cli.Commands {
		fmt.Fprintf(os.Stderr, "  %s\n", c.Usage)
	}
	fmt.Fprintf(os.Stderr, "\nRun from your app root (directory with go.mod) for module commands.\n")
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
	fmt.Println("  go mod tidy")
	fmt.Println("  go run .")
	fmt.Println("  go tool gogql module add posts   # add another module")
}

func runModule(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: gogql module <add|typedefs|resolvers|schema> <name>")
		os.Exit(1)
	}
	sub := args[0]
	fs := flag.NewFlagSet("module "+sub, flag.ExitOnError)
	force := fs.Bool("force", false, "overwrite existing generated files")
	project := fs.String("C", "", "project directory (default: search upward for go.mod)")
	_ = fs.Parse(args[1:])
	if fs.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "module name required\n\n")
		for _, c := range cli.Commands {
			if strings.HasPrefix(c.Name, "module ") && strings.Contains(c.Name, sub) {
				fmt.Fprintf(os.Stderr, "  %s\n", c.Usage)
			}
		}
		os.Exit(1)
	}
	name := fs.Arg(0)
	opts := scaffold.ModuleGenOptions{
		ProjectDir: *project,
		Name:       name,
		Force:      *force,
	}
	var err error
	switch sub {
	case "add":
		err = scaffold.AddModule(opts)
	case "typedefs":
		err = scaffold.AddTypedefs(opts)
	case "resolvers":
		err = scaffold.AddResolvers(opts)
	case "schema":
		err = scaffold.AddSchema(opts)
	default:
		fmt.Fprintf(os.Stderr, "unknown module subcommand %q (try add, typedefs, resolvers, schema)\n", sub)
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "module %s: %v\n", sub, err)
		os.Exit(1)
	}
	fmt.Printf("OK module %s (%s)\n", name, sub)
}
