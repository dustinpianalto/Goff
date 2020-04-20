package exts

import (
	"djpianalto.com/goff/djpianalto.com/goff/utils"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dustinpianalto/disgoman"
	"sort"
	"strconv"
	"strings"
	"time"
)

func pingCommand(ctx disgoman.Context, _ []string) {
	timeBefore := time.Now()
	msg, _ := ctx.Send("Pong!")
	took := time.Now().Sub(timeBefore)
	_, err := ctx.Session.ChannelMessageEdit(ctx.Message.ChannelID, msg.ID, fmt.Sprintf("Pong!\nPing Took **%s**", took.String()))
	if err != nil {
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Ping Failed",
			Error:   err,
		}
	}
}

func inviteCommand(ctx disgoman.Context, args []string) {
	var ids []string
	if len(args) == 0 {
		ids = []string{ctx.Session.State.User.ID}
	} else {
		for _, id := range args {
			ids = append(ids, id)
		}
	}
	for _, id := range ids {
		url := fmt.Sprintf("<https://discordapp.com/oauth2/authorize?client_id=%v&scope=bot>", id)
		_, err := ctx.Send(url)
		if err != nil {
			ctx.ErrorChannel <- disgoman.CommandError{
				Context: ctx,
				Message: "Couldn't send the invite link.",
				Error:   err,
			}
		}
	}
}

func gitCommand(ctx disgoman.Context, _ []string) {
	embed := &discordgo.MessageEmbed{
		Title: "Hi there, My code is on Github",
		Color: 0,
		URL:   "https://github.com/dustinpianalto/Goff",
	}
	_, err := ctx.Session.ChannelMessageSendEmbed(ctx.Channel.ID, embed)
	if err != nil {
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Git failed",
			Error:   err,
		}
	}
}

func sayCommand(ctx disgoman.Context, args []string) {
	resp := strings.Join(args, " ")
	resp = strings.ReplaceAll(resp, "@everyone", "@\ufff0everyone")
	resp = strings.ReplaceAll(resp, "@here", "@\ufff0here")
	_, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, resp)
	if err != nil {
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Say Failed",
			Error:   err,
		}
	}
}

func userCommand(ctx disgoman.Context, args []string) {
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
			ctx.ErrorChannel <- disgoman.CommandError{
				Context: ctx,
				Message: "Couldn't get that member",
				Error:   err,
			}
			return
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
	if err != nil {
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Couldn't send the user embed",
			Error:   err,
		}
	}
}
