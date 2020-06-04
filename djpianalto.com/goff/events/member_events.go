package events

import (
	"djpianalto.com/goff/djpianalto.com/goff/utils"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
	"time"
)

func OnGuildMemberAddLogging(s *discordgo.Session, member *discordgo.GuildMemberAdd) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic in OnGuildMemberAddLogging", r)
		}
	}()
	var channelID string
	row := utils.Database.QueryRow("SELECT logging_channel FROM guilds where id=$1", member.GuildID)
	err := row.Scan(&channelID)
	if err != nil || channelID == "" {
		return
	}
	guild, err := s.Guild(member.GuildID)
	if err != nil {
		log.Println(err)
		return
	}

	var title string
	if member.User.Bot {
		title = "Bot Joined"
	} else {
		title = "Member Joined"
	}

	thumb := &discordgo.MessageEmbedThumbnail{
		URL: member.User.AvatarURL(""),
	}

	int64ID, _ := strconv.ParseInt(member.User.ID, 10, 64)
	snow := utils.ParseSnowflake(int64ID)

	field := &discordgo.MessageEmbedField{
		Name:   "User was created:",
		Value:  utils.ParseDateString(snow.CreationTime),
		Inline: false,
	}

	joinTime, _ := member.JoinedAt.Parse()

	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: fmt.Sprintf("%v (%v) Has Joined the Server", member.User.Mention(), member.User.ID),
		Color:       0x0cc56a,
		Thumbnail:   thumb,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("Current Member Count: %v", guild.MemberCount),
			IconURL: guild.IconURL(),
		},
		Timestamp: joinTime.Format(time.RFC3339),
		Fields:    []*discordgo.MessageEmbedField{field},
	}
	s.ChannelMessageSendEmbed(channelID, embed)
}

func OnGuildMemberRemoveLogging(s *discordgo.Session, member *discordgo.GuildMemberRemove) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic in OnGuildMemberAddLogging", r)
		}
	}()
	timeNow := time.Now()
	var channelID string
	row := utils.Database.QueryRow("SELECT logging_channel FROM guilds where id=$1", member.GuildID)
	err := row.Scan(&channelID)
	if err != nil || channelID == "" {
		return
	}
	guild, err := s.Guild(member.GuildID)
	if err != nil {
		log.Println(err)
		return
	}

	var title string
	if member.User.Bot {
		title = "Bot Left"
	} else {
		title = "Member Left"
	}

	thumb := &discordgo.MessageEmbedThumbnail{
		URL: member.User.AvatarURL(""),
	}

	joinTime, _ := member.JoinedAt.Parse()
	duration := utils.ParseDurationString(timeNow.Sub(joinTime))

	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: fmt.Sprintf("%v (%v) Has Left the Server\nThey were here for %v", member.User.Mention(), member.User.ID, duration),
		Color:       0xff9431,
		Thumbnail:   thumb,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("Current Member Count: %v", guild.MemberCount),
			IconURL: guild.IconURL(),
		},
		Timestamp: timeNow.Format(time.RFC3339),
	}
	s.ChannelMessageSendEmbed(channelID, embed)
}
