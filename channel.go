package discutils

import (
	"github.com/bwmarrin/discordgo"
)

func PurgeChannel(s *discordgo.Session, chanID string) error {
	before := ""
	for {
		msgs, err := s.ChannelMessages(chanID, 100, before, "", "")
		if err != nil {
			return err
		}

		msgIDs := []string{}

		for i, msg := range msgs {
			msgIDs = append(msgIDs, msg.ID)
			if i == 99 {
				before = msg.ID
			}
		}

		if len(msgIDs) == 0 {
			break
		}

		err = s.ChannelMessagesBulkDelete(chanID, msgIDs)
		if err != nil {
			return err
		}
		if len(msgs) != 100 {
			break
		}
	}

	return nil
}
