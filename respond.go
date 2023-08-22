package discutils

import (
	"github.com/bwmarrin/discordgo"
)

type InteractionOpt int

const (
	// For first-time responses, updates a message. Usually used for buttons
	I_UPDATE InteractionOpt = 1 << iota
	// For editing after a response (eg. for a defer)
	I_EDIT_RESP

	// Only applicable to some operations. Should be used in conjunction with I_EDIT to remove components
	OPT_REMOVE_COMPONENTS
	// Make the message ephemeral. Should be used with I_UPDATE
	OPT_EPHEMERAL
)

func makeBitMask(opts []InteractionOpt) InteractionOpt {
	i := InteractionOpt(0)

	for _, v := range opts {
		i |= v
	}

	return i
}

// Respond to an interaction by deferring it
func IDefer(s *discordgo.Session, i *discordgo.Interaction, opts ...InteractionOpt) error {
	t := discordgo.InteractionResponseDeferredChannelMessageWithSource
	var flags discordgo.MessageFlags

	o := makeBitMask(opts)

	if BitMask(o, I_UPDATE) {
		t = discordgo.InteractionResponseDeferredMessageUpdate
	}
	if BitMask(o, OPT_EPHEMERAL) {
		flags = discordgo.MessageFlagsEphemeral
	}

	return s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: t,
		Data: &discordgo.InteractionResponseData{
			Flags: flags,
		},
	})
}

// Respond to an interaction with a single embed. Supports OPT_REMOVE_COMPONENTS - include it to remove components
func IEmbed(s *discordgo.Session, i *discordgo.Interaction, emb *discordgo.MessageEmbed, opts ...InteractionOpt) error {
	var flags discordgo.MessageFlags
	t := discordgo.InteractionResponseChannelMessageWithSource
	mask := makeBitMask(opts)

	if BitMask(mask, I_EDIT_RESP) {
		comps := &[]discordgo.MessageComponent{}
		if BitMask(mask, OPT_REMOVE_COMPONENTS) {
			comps = nil
		}

		_, err := s.InteractionResponseEdit(i, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				emb,
			},
			Components: comps,
		})
		return err
	}

	if BitMask(mask, I_UPDATE) {
		t = discordgo.InteractionResponseUpdateMessage
	}
	if BitMask(mask, OPT_EPHEMERAL) {
		flags = discordgo.MessageFlagsEphemeral
	}

	return s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: t,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				emb,
			},
			Flags: flags,
		},
	})
}

// Respond to an interaction with an error embed.
// This method calls IEmbed under the hood, so OPT_REMOVE_COMPONENTS works here too.
// The embed is formatted using ErrEmbed
//
// WARNING: IS_EPHEMERAL is reversed in here! That is, by default (no opts) it has IS_EPHEMERAL, but, if IS_EPHEMERAL is present, it will make it **not** ephemeral!
func IError(s *discordgo.Session, i *discordgo.Interaction, errMsg string, opts ...InteractionOpt) error {
	return IEmbed(s, i, ErrEmbed(errMsg), makeBitMask(opts)^OPT_EPHEMERAL)
}

// Respond with an error embed and components.
// WARNING: IS_EPHEMERAL is reversed in here! That is, by default (no opts) it has IS_EPHEMERAL, but, if IS_EPHEMERAL is present, it will make it **not** ephemeral!
func IErrorComponents(s *discordgo.Session, i *discordgo.Interaction, errMsg string, components []discordgo.MessageComponent, opts ...InteractionOpt) error {
	return IResp(s, i, &IRespOpts{
		Embeds: []*discordgo.MessageEmbed{
			ErrEmbed(errMsg),
		},
		Comps: components,
	}, makeBitMask(opts)^OPT_EPHEMERAL)
}

type IRespOpts struct {
	Embeds  []*discordgo.MessageEmbed
	Comps   []discordgo.MessageComponent
	Content *string
}

// Respond to an interaction, the 'lowest level' of the utils that this package provides.
// Can auto resolve components into Action Row (so you can directly pass buttons, without the need to nest them)
// This supports OPT_REMOVE_COMPONENTS. if len(conf.Comps) == 0 && OPT_REMOVE_COMPONENTS && I_EDIT_RESP, then components will be removed.
// Argument opts is a bit mask - eg. I_EDIT_RESP | OPT_REMOVE_COMPONENTS
func IResp(s *discordgo.Session, i *discordgo.Interaction, conf *IRespOpts, opts InteractionOpt) error {
	var flags discordgo.MessageFlags
	t := discordgo.InteractionResponseChannelMessageWithSource

	components := conf.Comps

	if len(components) > 0 {
		_, ok1 := components[0].(discordgo.ActionsRow)
		_, ok2 := components[0].(*discordgo.ActionsRow)

		if !ok1 && !ok2 {
			components = []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: components,
				},
			}
		}
	} else if BitMask(opts, OPT_REMOVE_COMPONENTS) {
		components = nil
	}

	if BitMask(opts, I_EDIT_RESP) {
		var compsToSend *[]discordgo.MessageComponent

		if components != nil {
			compsToSend = &components
		}

		_, err := s.InteractionResponseEdit(i, &discordgo.WebhookEdit{
			Embeds:     &conf.Embeds,
			Components: compsToSend,
			Content:    conf.Content,
		})

		return err
	}

	if BitMask(opts, I_UPDATE) {
		t = discordgo.InteractionResponseUpdateMessage
	}

	if BitMask(opts, OPT_EPHEMERAL) {
		flags = discordgo.MessageFlagsEphemeral
	}

	content := ""
	if conf.Content != nil {
		content = *conf.Content
	}

	return s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: t,
		Data: &discordgo.InteractionResponseData{
			Embeds:     conf.Embeds,
			Content:    content,
			Components: components,
			Flags:      flags,
		},
	})
}

// Shortcut for autocomplete response
func IAutocomplete(s *discordgo.Session, i *discordgo.Interaction, choices []*discordgo.ApplicationCommandOptionChoice) error {
	return s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
}

// Shortcut for modal response
func IModal(s *discordgo.Session, i *discordgo.Interaction, customID, title string, components []discordgo.MessageComponent) error {
	return s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Components: components,
			CustomID:   customID,
			Title:      title,
		},
	})
}
