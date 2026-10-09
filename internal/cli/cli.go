/*
|--------------------------------------------------------------------------
| CLI metadata
|--------------------------------------------------------------------------
|
| Command names, summaries, and version string for the gogql binary.
|
*/

package cli

const Version = "0.2.0"

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
		Summary: "Scaffold a new GraphQL server project",
		Usage:   "gogql init [directory] [-module path] [-gogql replace]",
	},
	{
		Name:    "module add",
		Summary: "Create a module (module.go, typedefs.go, resolvers.go, schema/)",
		Usage:   "gogql module add <name> [-force]",
	},
	{
		Name:    "module typedefs",
		Summary: "Add typedefs.go and schema/*.graphql to a module",
		Usage:   "gogql module typedefs <name> [-force]",
	},
	{
		Name:    "module resolvers",
		Summary: "Add resolvers.go (and module.go if missing)",
		Usage:   "gogql module resolvers <name> [-force]",
	},
	{
		Name:    "module schema",
		Summary: "Add schema/*.graphql SDL files only",
		Usage:   "gogql module schema <name> [-force]",
	},
	{
		Name:    "version",
		Summary: "Print CLI version and list commands",
		Usage:   "gogql version",
	},
}
