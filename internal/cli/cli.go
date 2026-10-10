/*
|--------------------------------------------------------------------------
| CLI metadata
|--------------------------------------------------------------------------
|
| Single source of truth for gogql command names, one-line summaries, usage
| strings, and Version. cmd/gogql/main.go prints Commands on version/help;
| keep this list updated when adding subcommands.
|
| Key vars: Version, Commands.
|
*/

package cli

const Version = "0.3.3"

// Command describes a gogql subcommand for help and version output.
type Command struct {
	Name    string
	Summary string
	Usage   string
}

// Commands is the canonical list of gogql CLI commands.
var Commands = []Command{
	{
		Name:    "init",
		Summary: "Scaffold a new server (default src/ layout: common, modules, schema, config, utils)",
		Usage:   "gogql init [directory] [-module path] [-gogql replace]",
	},
	{
		Name:    "module add",
		Summary: "Create a domain module (*.graphql, *.resolvers.go, *.model.go, *.service.go)",
		Usage:   "gogql module add <name> [-force]",
	},
	{
		Name:    "module typedefs",
		Summary: "Add <name>.graphql SDL for a domain module",
		Usage:   "gogql module typedefs <name> [-force]",
	},
	{
		Name:    "module resolvers",
		Summary: "Add resolver, model, and service files for a domain module",
		Usage:   "gogql module resolvers <name> [-force]",
	},
	{
		Name:    "module schema",
		Summary: "Add domain .graphql SDL only (alias for module typedefs)",
		Usage:   "gogql module schema <name> [-force]",
	},
	{
		Name:    "version",
		Summary: "Print CLI version and list commands",
		Usage:   "gogql version",
	},
}
