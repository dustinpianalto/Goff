package guild_management

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/dustinpianalto/goff/internal/postgres"
)

func OnMessageUpdate(session *discordgo.Session, m *discordgo.MessageUpdate) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic in OnMessageUpdate", r)
		}
	}()
	msg := m.BeforeUpdate
	if msg.Author.Bot {
		return
	}
	var channelID string
	row := postgres.DB.QueryRow("SELECT logging_channel FROM guilds where id=$1", msg.GuildID)
	err := row.Scan(&channelID)
	if err != nil || channelID == "" {
		return
	}
	channel, err := session.State.Channel(msg.ChannelID)
	if err != nil {
		log.Println(err)
		return
	}
	if m.Content == "" {
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
			log.Println("Recovered from panic in OnMessageDelete", r)
		}
	}()
	msg := m.BeforeDelete
	if msg == nil {
		log.Printf("Message Deleted but the original message was not in my cache so we are ignoring it.\nMessage ID: %v\nGuild ID: %v\nChannel ID: %v\n", m.ID, m.GuildID, m.ChannelID)
		return
	}
	if msg.Author.Bot {
		return
	}
	var channelID string
	row := postgres.DB.QueryRow("SELECT logging_channel FROM guilds where id=$1", msg.GuildID)
	err := row.Scan(&channelID)
	if err != nil || channelID == "" {
		return
	}
	channel, err := session.State.Channel(msg.ChannelID)
	if err != nil {
		log.Println(err)
		return
	}
	desc := ""
	al, err := session.GuildAuditLog(msg.GuildID, "", "", 72, 1)
	if err != nil {
		log.Println(err)
	} else {
		for _, log := range al.AuditLogEntries {
			if log.TargetID == msg.Author.ID && log.Options.ChannelID == msg.ChannelID {
				user, err := session.User(log.UserID)
				if err == nil {
					desc = fmt.Sprintf("**Content:** %v\nIn Channel: %v\nDeleted By: %v", msg.Content, channel.Mention(), user.Mention())
				}
				break
			}
		}
	}
	if desc == "" {
		desc = fmt.Sprintf("**Content:** %v\nIn Channel: %v", msg.Content, channel.Mention())
	}
	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("Message Deleted: %v", msg.ID),
		Description: desc,
		Color:       session.State.UserColor(msg.Author.ID, channelID),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("Author: %v", msg.Author.String()),
			IconURL: msg.Author.AvatarURL(""),
		},
	}
	session.ChannelMessageSendEmbed(channelID, embed)
}
