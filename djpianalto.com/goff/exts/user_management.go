package exts

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dustinpianalto/disgoman"
	"strings"
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

	if !disgoman.HasHigherRole(ctx.Session, ctx.Guild.ID, ctx.User.ID, member.User.ID) {
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "I don't have a high enough role to kick that person",
			Error:   errors.New("need higher role"),
		}
		return
	}

	if len(args) > 1 {
		err = ctx.Session.GuildMemberDeleteWithReason(ctx.Guild.ID, member.User.ID, strings.Join(args[1:], " "))
	} else {
		err = ctx.Session.GuildMemberDelete(ctx.Guild.ID, member.User.ID)
	}
	if err != nil {
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: fmt.Sprintf("Something went wrong kicking %v", member.User.Username),
			Error:   err,
		}
	}
}
