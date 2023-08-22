package discutils

import (
	"github.com/bwmarrin/discordgo"
)

var (
	COLOR_PRIMARY = 0xad6bee
	COLOR_SUCCESS = 0x08dd7e
	COLOR_DANGER  = 0xA51D2A
)

var BaseEmbed = discordgo.MessageEmbed{
	Title: CHAR_ZWS,
	Color: COLOR_PRIMARY,
}

var ErrEmbed func(err string) *discordgo.MessageEmbed

// Generate an embed using the base embed
func GEmbed(title, desc string) *discordgo.MessageEmbed {
	emb := BaseEmbed

	emb.Title = title

	if title == "" {
		emb.Title = CHAR_ZWS
	}

	emb.Description = desc

	return &emb
}
