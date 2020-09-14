package events

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dustinpianalto/goff/utils"
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
	guild, err := s.State.Guild(member.GuildID)
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
	guild, err := s.State.Guild(member.GuildID)
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

	desc := ""
	al, err := s.GuildAuditLog(member.GuildID, "", "", 20, 1)
	if err != nil {
		log.Println(err)
	} else {
		for _, log := range al.AuditLogEntries {
			if log.TargetID == member.User.ID {
				int64ID, _ := strconv.ParseInt(log.ID, 10, 64)
				logSnow := utils.ParseSnowflake(int64ID)
				if timeNow.Sub(logSnow.CreationTime).Seconds() <= 10 {
					user, err := s.User(log.UserID)
					if err == nil {
						desc = fmt.Sprintf("%v (%v) was Kicked by: %v\nReason: %v", member.User.String(), member.User.ID, user.String(), log.Reason)
					} else {
						desc = fmt.Sprintf("%v (%v) was Kicked by: %v\nReason: %v", member.User.String(), member.User.ID, log.UserID, log.Reason)
					}
					break
				}
			}
		}
	}
	if desc == "" {
		desc = fmt.Sprintf("%v (%v) Has Left the Server", member.User.String(), member.User.ID)
	}

	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: desc,
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

func AddMemberToDatabase(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic in AddMemberToDatabase", r)
		}
	}()
	queryString := `INSERT INTO users (id, banned, logging, is_active, is_staff, is_admin) 
						VALUES ($1, false, false, true, false, false)`
	_, err := utils.Database.Exec(queryString, m.User.ID)
	if err != nil {
		log.Println(fmt.Errorf("error inserting %s into database: %w", m.User.ID, err))
		return
	}
	log.Printf("New User: %s\n", m.User.ID)
}

func MarkMemberInactive(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic in MarkMemberInactive")
		}
	}()
	queryString := `UPDATE users SET is_active = false WHERE id = $1`
	_, err := utils.Database.Exec(queryString, m.User.ID)
	if err != nil {
		log.Println(fmt.Errorf("error marking %s as inactive: %w", m.User.ID, err))
		return
	}
	log.Println("User left: %s")
}
