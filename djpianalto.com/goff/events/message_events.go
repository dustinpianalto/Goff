package events

import (
	"fmt"

	"djpianalto.com/goff/djpianalto.com/goff/utils"
	"github.com/bwmarrin/discordgo"
)

func OnMessageUpdate(session *discordgo.Session, m *discordgo.MessageUpdate) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic in OnMessageUpdate", r)
		}
	}()
	msg := m.BeforeUpdate
	if msg.Author.Bot {
		return
	}
	var channelID string
	row := utils.Database.QueryRow("SELECT logging_channel FROM guilds where id=$1", msg.GuildID)
	err := row.Scan(&channelID)
	if err != nil || channelID == "" {
		return
	}
	channel, err := session.State.Channel(msg.ChannelID)
	if err != nil {
		return
	}
	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("Message Edited: %v", msg.ID),
		Description: fmt.Sprintf("**Before:** %v\n**After:** %v\nIn Channel: %v", msg.Content, m.Content, channel.Mention()),
		Color:       session.State.UserColor(msg.Author.ID, channelID),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("Author: %v", msg.Author.String()),
			IconURL: msg.Author.AvatarURL(""),
		},
	}
	session.ChannelMessageSendEmbed(channelID, embed)
}

func OnMessageDelete(session *discordgo.Session, m *discordgo.MessageDelete) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic in OnMessageDelete", r)
		}
	}()
	msg := m.BeforeDelete
	if msg.Author.Bot {
		return
	}
	var channelID string
	row := utils.Database.QueryRow("SELECT logging_channel FROM guilds where id=$1", msg.GuildID)
	err := row.Scan(&channelID)
	if err != nil || channelID == "" {
		return
	}
	channel, err := session.State.Channel(msg.ChannelID)
	if err != nil {
		return
	}
	al, err := session.GuildAuditLog(msg.GuildID, "", "", 72, 100)
	if err != nil {
		fmt.Println(err)
	}
	for log := range al.AuditLogEntries {
		fmt.Println(log.TargetID, log.UserID, log.ID)
	}
	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("Message Deleted: %v", msg.ID),
		Description: fmt.Sprintf("**Content:** %v\nIn Channel: %v", msg.Content, channel.Mention()),
		Color:       session.State.UserColor(msg.Author.ID, channelID),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("Author: %v", msg.Author.String()),
			IconURL: msg.Author.AvatarURL(""),
		},
	}
	session.ChannelMessageSendEmbed(channelID, embed)
}
