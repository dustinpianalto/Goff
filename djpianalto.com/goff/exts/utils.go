package exts

import (
	"djpianalto.com/goff/djpianalto.com/goff/utils"
	"fmt"
	"github.com/MikeModder/anpan"
	"github.com/bwmarrin/discordgo"
	"sort"
	"strconv"
	"strings"
	"time"
)

func pingCommand(ctx anpan.Context, _ []string) error {
	timeBefore := time.Now()
	msg, _ := ctx.Reply("Pong!")
	took := time.Now().Sub(timeBefore)
	_, err := ctx.Session.ChannelMessageEdit(ctx.Message.ChannelID, msg.ID, fmt.Sprintf("Pong!\nPing Took **%s**", took.String()))
	return err
}

func sayCommand(ctx anpan.Context, args []string) error {
	resp := strings.Join(args, " ")
	resp = strings.ReplaceAll(resp, "@everyone", "@\ufff0everyone")
	resp = strings.ReplaceAll(resp, "@here", "@\ufff0here")
	_, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, resp)
	return err
}

func userCommand(ctx anpan.Context, args []string) error {
	var member *discordgo.Member
	if len(args) == 0 {
		member, _ = ctx.Session.GuildMember(ctx.Guild.ID, ctx.Message.Author.ID)
	} else {
		var err error
		if len(ctx.Message.Mentions) > 0 {
			member, err = ctx.Session.GuildMember(ctx.Guild.ID, ctx.Message.Mentions[0].ID)
		} else {
			member, err = ctx.Session.GuildMember(ctx.Guild.ID, args[0])
		}
		if err != nil {
			return err
		}
	}
	thumb := &discordgo.MessageEmbedThumbnail{
		URL: member.User.AvatarURL(""),
	}

	var botString string
	if member.User.Bot {
		botString = "BOT"
	} else {
		botString = ""
	}

	var roles []*discordgo.Role
	for _, roleID := range member.Roles {
		role, _ := ctx.Session.State.Role(ctx.Guild.ID, roleID)
		roles = append(roles, role)
	}
	sort.Slice(roles, func(i, j int) bool { return roles[i].Position > roles[j].Position })
	var roleMentions []string
	for _, role := range roles {
		roleMentions = append(roleMentions, role.Mention())
	}
	rolesString := strings.Join(roleMentions, " ")

	rolesField := &discordgo.MessageEmbedField{
		Name:   "Roles:",
		Value:  rolesString,
		Inline: false,
	}

	guildJoinTime, _ := member.JoinedAt.Parse()
	guildJoinedField := &discordgo.MessageEmbedField{
		Name:   "Joined Guild:",
		Value:  utils.ParseDateString(guildJoinTime),
		Inline: false,
	}

	int64ID, _ := strconv.ParseInt(member.User.ID, 10, 64)
	s := utils.ParseSnowflake(int64ID)
	discordJoinedField := &discordgo.MessageEmbedField{
		Name:   "Joined Discord:",
		Value:  utils.ParseDateString(s.CreationTime),
		Inline: false,
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("%v#%v  %v", member.User.Username, member.User.Discriminator, botString),
		Description: fmt.Sprintf("**%v** (%v)", member.Nick, member.User.ID),
		Color:       0,
		Thumbnail:   thumb,
		Fields: []*discordgo.MessageEmbedField{
			guildJoinedField,
			discordJoinedField,
			rolesField,
		},
	}
	_, err := ctx.Session.ChannelMessageSendEmbed(ctx.Channel.ID, embed)
	return err
}
