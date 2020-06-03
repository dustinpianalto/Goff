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
		Title:       fmt.Sprintf("Message Edited: %v", m.ID),
		Description: fmt.Sprintf("**Before:** %v\n**After:** %v\nIn Channel: %v", m.BeforeUpdate.Content, m.Content, channel.Mention()),
		Color:       session.State.UserColor(m.Author.ID, channelID),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("Author: %v", m.Author.String()),
			IconURL: m.Author.AvatarURL(""),
		},
	}
	session.ChannelMessageSendEmbed(channelID, embed)
}

func OnMessageDelete(session *discordgo.Session, m *discordgo.MessageDelete) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic in OnMessageUpdate", r)
		}
	}()
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
		Title:       fmt.Sprintf("Message Deleted: %v", m.ID),
		Description: fmt.Sprintf("**Content:** %v\nIn Channel: %v", m.Content, channel.Mention()),
		Color:       session.State.UserColor(m.Author.ID, channelID),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("Author: %v", m.Author.String()),
			IconURL: m.Author.AvatarURL(""),
		},
	}
	session.ChannelMessageSendEmbed(channelID, embed)
}
