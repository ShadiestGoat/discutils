package events

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
)

// Parses a command's arguments into a map.
// The returning type will be {argument -> *opt}
// For sub commands (or sub command groups), the key will be "cmd-{layer}" (where layer starts at 1), and the option will be a string
// Optional arguments that aren't specified by the user will not be in the map
func ParseCommand(d discordgo.ApplicationCommandInteractionData) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	return parseCmd(d.Options, 1)
}

func parseCmd(d []*discordgo.ApplicationCommandInteractionDataOption, layer int) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	r := map[string]*discordgo.ApplicationCommandInteractionDataOption{}

	for _, opt := range d {
		switch opt.Type {
		case discordgo.ApplicationCommandOptionSubCommandGroup, discordgo.ApplicationCommandOptionSubCommand:
			tmp := parseCmd(opt.Options, layer+1)
			for key, val := range tmp {
				if r[key] != nil {
					panic(key + " is not unique!")
				}
				r[key] = val
			}
			// so that you can do .StringValue()
			opt.Type = discordgo.ApplicationCommandOptionString
			opt.Value = opt.Name

			r["cmd-"+strconv.Itoa(layer)] = opt
		default:
			r[opt.Name] = opt
		}
	}

	return r
}

// Same as ParseCommand, but for modals
// The returning map will be of type {CustomID -> Value}
func ParseModal(d discordgo.ModalSubmitInteractionData) map[string]string {
	m := map[string]string{}

	for _, row := range d.Components {
		r, ok := row.(*discordgo.ActionsRow)
		if !ok {
			continue
		}

		for _, inp := range r.Components {
			if inp.Type() != discordgo.TextInputComponent {
				continue
			}
			c := inp.(*discordgo.TextInput)
			m[c.CustomID] = c.Value
		}
	}

	return m
}
