package exts

import (
	"djpianalto.com/goff/djpianalto.com/goff/utils"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dustinpianalto/disgoman"
	"strings"
	"time"
)

func kickUser(ctx disgoman.Context, args []string) {
	var member *discordgo.Member
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

	if !disgoman.HasHigherRole(ctx.Session, ctx.Guild.ID, ctx.Message.Author.ID, member.User.ID) {
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "You must have a higher role than the person you are trying to kick",
			Error:   errors.New("need higher role"),
		}
		return
	}

	if !disgoman.HasHigherRole(ctx.Session, ctx.Guild.ID, ctx.Session.State.User.ID, member.User.ID) {
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "I don't have a high enough role to kick that person",
			Error:   errors.New("need higher role"),
		}
		return
	}

	var reason string
	if len(args) > 1 {
		reason = strings.Join(args[1:], " ")
	} else {
		reason = "No Reason Given"
	}
	auditReason := fmt.Sprintf("%v#%v: %v", ctx.User.Username, ctx.User.Discriminator, reason)
	err = ctx.Session.GuildMemberDeleteWithReason(ctx.Guild.ID, member.User.ID, auditReason)

	if err != nil {
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: fmt.Sprintf("Something went wrong kicking %v", member.User.Username),
			Error:   err,
		}
		return
	}

	event := &utils.LogEvent{
		Embed: discordgo.MessageEmbed{
			Title: "User Kicked",
			Description: fmt.Sprintf(
				"User %v#%v was kicked by %v.\nReason: %v",
				member.User.Username,
				member.User.Discriminator,
				ctx.Message.Author.Username,
				reason),
			Timestamp: time.Now().Format("2006-01-02 15:04:05"),
			Color:     0xff8c00,
		},
		GuildID: ctx.Guild.ID,
		Session: *ctx.Session,
	}
	utils.LoggingChannel <- event
	_, _ = ctx.Send(fmt.Sprintf("User %v#%v has been kicked.", member.User.Username, member.User.Discriminator))
}
