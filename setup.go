package discutils

import "github.com/bwmarrin/discordgo"

type Options struct {
	// This function will be called after setting all the colors.
	// This is a function so that you can use that to your advantage.
	Base       func() discordgo.MessageEmbed
	ErrorEmbed func(err string) *discordgo.MessageEmbed

	ColorPrimary *int
	ColorSuccess *int
	ColorDanger  *int
}

// Configure global embed options (and colors)
func InitEmbed(opts *Options) {
	if opts == nil {
		opts = &Options{
			ErrorEmbed: func(err string) *discordgo.MessageEmbed {
				return &discordgo.MessageEmbed{
					Title:       "Error!",
					Description: err,
					Color:       COLOR_DANGER,
				}
			},
		}
	}

	if opts.ColorPrimary != nil {
		COLOR_PRIMARY = *opts.ColorPrimary
	}
	if opts.ColorSuccess != nil {
		COLOR_SUCCESS = *opts.ColorSuccess
	}
	if opts.ColorDanger != nil {
		COLOR_DANGER = *opts.ColorDanger
	}

	if opts.Base == nil {
		BaseEmbed = discordgo.MessageEmbed{
			Title: CHAR_ZWS,
			Color: COLOR_PRIMARY,
		}
	} else {
		BaseEmbed = opts.Base()
	}

	ErrEmbed = opts.ErrorEmbed
}
