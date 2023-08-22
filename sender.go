package discutils

import "github.com/bwmarrin/discordgo"

// Send a message to a channel. If the channel is a news channel, auto-crosspost it.
func SendMessage(s *discordgo.Session, channelID string, msg *discordgo.MessageSend) (resp *discordgo.Message, err error) {
	if s == nil {
		return nil, nil
	}

	resp, err = s.ChannelMessageSendComplex(channelID, msg)
	if err != nil {
		return
	}

	channel := GetChannel(s, channelID)

	if channel == nil || (channel.Type != discordgo.ChannelTypeGuildNews && channel.Type != discordgo.ChannelTypeGuildNewsThread) {
		return
	}

	_, err = s.ChannelMessageCrosspost(channelID, resp.ID)
	return
}

// Shortcut to DMing a user (by also calling UserChannelCreate)
func DM(s *discordgo.Session, userID string, msg *discordgo.MessageSend) (resp *discordgo.Message, err error) {
	c, err := s.UserChannelCreate(userID)
	if err != nil {
		return
	}

	resp, err = s.ChannelMessageSendComplex(c.ID, msg)

	return
}
