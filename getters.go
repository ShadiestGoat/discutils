package discutils

import "github.com/bwmarrin/discordgo"

func GetChannel(s *discordgo.Session, channelID string) *discordgo.Channel {
	c, err := s.State.Channel(channelID)
	if err == nil && c != nil {
		return c
	}

	c, err = s.Channel(channelID)

	if err != nil {
		c = nil
	}

	return c
}

func GetMember(s *discordgo.Session, guildID, userID string) *discordgo.Member {
	m, err := s.State.Member(guildID, userID)
	if err == nil && m != nil {
		return m
	}
	m, err = s.GuildMember(guildID, userID)
	if err != nil {
		m = nil
	}

	return m
}

func GetMessage(s *discordgo.Session, channelID, msgID string) *discordgo.Message {
	m, err := s.State.Message(channelID, msgID)
	if err == nil && m != nil {
		return m
	}

	m, err = s.ChannelMessage(channelID, msgID)
	if err != nil {
		m = nil
	}

	return m
}
