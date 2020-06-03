package events

import (
	"fmt"

	"djpianalto.com/goff/djpianalto.com/goff/utils"
	"github.com/bwmarrin/discordgo"
)

func OnMessageUpdate(session *discordgo.Session, m *discordgo.MessageUpdate) {
	if m.Author.Bot {
		return
	}
	var channelID string
	row := utils.Database.QueryRow("SELECT logging_channel FROM guilds where id=$1", m.GuildID)
	err := row.Scan(&channelID)
	if err != nil || channelID == "" {
		return
	}
	channel, err := session.State.Channel(m.ChannelID)
	if err != nil {
		return
	}
	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("Edited Message: %v in %v", m.ID, channel.Mention()),
		Description: fmt.Sprintf("**Before:** %v\n**After:** %v", m.BeforeUpdate.Content, m.Content),
		Color:       session.State.UserColor(m.Author.ID, channelID),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("Author: %v", m.Author.String()),
			IconURL: m.Author.AvatarURL(""),
		},
	}
	session.ChannelMessageSendEmbed(channelID, embed)
}
