package events

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Register all the event handlers.
// It is best to run this after all the events & commands were registered, but before you open the connection
func Init(s *discordgo.Session) {
	s.AddHandler(handleInteraction)
	s.AddHandler(MessageReactionAdd.handle)
	s.AddHandler(MemberJoin.handle)
	s.AddHandler(MessageReactionRemove.handle)
	s.AddHandler(MessageCreate.handle)
	s.AddHandler(MessageRemove.handle)
	s.AddHandlerOnce(Ready.handle)
}

func commandIsEqual(cmd1, cmd2 *discordgo.ApplicationCommand) bool {
	if cmd1.DefaultMemberPermissions == nil {
		cmd1.DefaultMemberPermissions = CommandPerms()
	}
	if cmd2.DefaultMemberPermissions == nil {
		cmd2.DefaultMemberPermissions = CommandPerms()
	}

	if cmd1.Description != cmd2.Description || len(cmd1.Options) != len(cmd2.Options) {
		return false
	}

	if *cmd1.DefaultMemberPermissions != *cmd2.DefaultMemberPermissions {
		return false
	}

	for i := range cmd1.Options {
		if !optIsEqual(cmd1.Options[i], cmd2.Options[i]) {
			return false
		}
	}

	return true
}

func optIsEqual(opt1, opt2 *discordgo.ApplicationCommandOption) bool {
	if opt1.Name != opt2.Name || opt1.Description != opt2.Description || opt1.Type != opt2.Type || opt1.Required != opt2.Required {
		return false
	}

	if len(opt1.Options) != len(opt2.Options) || len(opt1.Choices) != len(opt2.Choices) || len(opt1.ChannelTypes) != len(opt2.ChannelTypes) {
		return false
	}

	for i := range opt1.ChannelTypes {
		if opt1.ChannelTypes[i] != opt2.ChannelTypes[i] {
			return false
		}
	}

	for i := range opt1.Choices {
		c1 := opt1.Choices[i]
		c2 := opt2.Choices[i]

		if c1.Name != c2.Name || fmt.Sprint(c1.Value) != fmt.Sprint(c2.Value) {
			return false
		}
	}

	for i := range opt1.Options {
		if !optIsEqual(opt1.Options[i], opt2.Options[i]) {
			return false
		}
	}

	return true
}

// Register all the commands onto discord. This works as a 'patch' - only register commands that you need to register.
//
// This should be run after the connection is open! Easiest way to get the appID is s.State.User.ID
func InitCommands(s *discordgo.Session, appID string) error {
	curCommands, err := s.ApplicationCommands(appID, "")
	if err != nil {
		return err
	}

	oldCommands := map[string]*discordgo.ApplicationCommand{}

	for i, cmd := range curCommands {
		oldCommands[cmd.Name] = curCommands[i]

		if _, ok := commands[cmd.Name]; !ok {
			err := s.ApplicationCommandDelete(appID, "", cmd.ID)
			if err != nil {
				return err
			}
		}
	}

	for _, v := range commands {
		if oldCommands[v.Name] == nil || !commandIsEqual(oldCommands[v.Name], v) {
			v, err = s.ApplicationCommandCreate(s.State.User.ID, "", v)
			if err != nil {
				return err
			}
		} else {
			v = oldCommands[v.Name]
		}

		// Update with an ID, etc
		commands[v.Name] = v
	}

	return nil
}
