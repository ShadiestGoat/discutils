package discutils

import (
	"github.com/bwmarrin/discordgo"
)

func IsMemberPremium(mem *discordgo.Member) bool {
	return mem != nil && mem.PremiumSince != nil && !mem.PremiumSince.IsZero()
}

// Returns true if the member has the role with id roleNeeded
func HasRole(m *discordgo.Member, roleNeeded string) bool {
	for _, r := range m.Roles {
		if r == roleNeeded {
			return true
		}
	}

	return false
}

// Returns true if the member has all the roles mentioned in argument "roles"
func HasAllRoles(m *discordgo.Member, roles []string) bool {
	rolesNeeded := map[string]bool{}

	for _, r := range m.Roles {
		rolesNeeded[r] = true
	}

	for _, r := range m.Roles {
		delete(rolesNeeded, r)
	}

	return len(rolesNeeded) == 0
}

// Returns true if the member has any of the roles mentioned in argument "roles"
func HasAnyRoles(m *discordgo.Member, roles []string) bool {
	rolesNeeded := map[string]bool{}
	for _, r := range m.Roles {
		rolesNeeded[r] = true
	}

	for _, r := range m.Roles {
		if rolesNeeded[r] {
			return true
		}
	}

	return false
}

func EMessageURL(msg *discordgo.Message) string {
	g := msg.GuildID

	if g == "" {
		g = "@me"
	}

	return MessageURL(g, msg.ChannelID, msg.ID)
}

// Use Guild == "@me" for DM
func MessageURL(Guild, Channel, Message string) string {
	return "https://discord.com/channels/" + Guild + "/" + Channel + "/" + Message
}

// Return the display name for the user, with a fallback to the username
func Username(u *discordgo.User) string {
	if u.GlobalName == "" {
		return u.Username
	}
	return u.GlobalName
}

// Returns the name that should be used for a member
//
// If a nickname exists, it returns that. If not, it returns the
func MemberName(mem *discordgo.Member) string {
	if mem.Nick == "" {
		return Username(mem.User)
	}

	return mem.Nick
}

// Updates a message with itself. Used to fake-acknowledge an interaction.
func UpdateMessageWithSelf(s *discordgo.Session, i *discordgo.Interaction) {
	s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:    i.Message.Content,
			Embeds:     i.Message.Embeds,
			Components: i.Message.Components,
		},
	})
}

// Get the author of an interaction
func IAuthor(i *discordgo.Interaction) *discordgo.User {
	if i.Member != nil {
		return i.Member.User
	}

	return i.User
}

// Returns the author of the original interaction message
// Use only for buttons!
func IMessageAuthor(i *discordgo.Interaction) *discordgo.User {
	if i == nil || i.Message == nil {
		return nil
	}

	j := i.Message.Interaction

	if j.Member != nil {
		return j.Member.User
	}

	return j.User
}

// Returns true if userID is the author of the message's interaction
// This can be used with regular messages or, if its a button interaction, i.Message
func IOwnInteraction(msg *discordgo.Message, userID string) bool {
	if msg.Interaction == nil {
		return false
	}
	id := ""

	if msg.Interaction.Member != nil {
		id = msg.Interaction.Member.User.ID
	} else if msg.Interaction.User != nil {
		id = msg.Interaction.User.ID
	}

	return id == userID
}
