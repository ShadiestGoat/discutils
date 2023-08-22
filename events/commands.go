package events

import "github.com/bwmarrin/discordgo"

var commands = map[string]*discordgo.ApplicationCommand{}

// Register a command and it's handler
func RegisterCommand(cmd *discordgo.ApplicationCommand, handler HandlerCommand) {
	commands[cmd.Name] = cmd
	commandHandlers[cmd.Name] = handler
}

func CommandPerms(perms ...int) *int64 {
	var permsGiven int64

	for _, p := range perms {
		permsGiven |= int64(p)
	}

	return &permsGiven
}

// Returns a command ID, or an empty string if the command is not found/not registered
func CommandID(cmdName string) string {
	if cmd, ok := commands[cmdName]; ok {
		return cmd.ID
	}

	return ""
}

// Creates a command mention. subCommands should be a space separated sub command list (or empty).
// If the command does not exist, this will return an empty string.
func CommandMention(cmdName string, subCommands string) string {
	// </NAME SUBCOMMAND_GROUP SUBCOMMAND:ID>
	id := CommandID(cmdName)
	if id == "" {
		return ""
	}

	inside := cmdName
	if subCommands != "" {
		inside += " " + subCommands
	}

	return "</" + inside + ":" + id + ">"
}

func Command(name string) *discordgo.ApplicationCommand {
	return commands[name]
}
