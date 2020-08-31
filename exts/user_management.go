package exts

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dustinpianalto/disgoman"
	"github.com/dustinpianalto/goff/utils"
)

func kickUserCommand(ctx disgoman.Context, args []string) {
	var member *discordgo.Member
	var err error
	if len(ctx.Message.Mentions) > 0 {
		member, err = ctx.Session.GuildMember(ctx.Guild.ID, ctx.Message.Mentions[0].ID)
	} else if len(args) >= 1 {
		idString := args[0]
		if strings.HasPrefix(idString, "<@!") && strings.HasSuffix(idString, ">") {
			idString = idString[3 : len(idString)-1]
		}
		member, err = ctx.Session.GuildMember(ctx.Guild.ID, idString)
	} else {
		err = errors.New("that is not a valid id")
	}
	if err != nil {
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Couldn't get that member",
			Error:   err,
		}
		return
	}

	if higher, _ := disgoman.HasHigherRole(ctx.Session, ctx.Guild.ID, ctx.Message.Author.ID, member.User.ID); !higher {
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "You must have a higher role than the person you are trying to kick",
			Error:   errors.New("need higher role"),
		}
		return
	}

	if higher, _ := disgoman.HasHigherRole(ctx.Session, ctx.Guild.ID, ctx.Session.State.User.ID, member.User.ID); !higher {
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
		Session: ctx.Session,
	}
	utils.LoggingChannel <- event
	_, _ = ctx.Send(fmt.Sprintf("User %v#%v has been kicked.", member.User.Username, member.User.Discriminator))
}

func banUserCommand(ctx disgoman.Context, args []string) {
	var user *discordgo.User
	var err error
	if len(ctx.Message.Mentions) > 0 {
		user, err = ctx.Session.User(ctx.Message.Mentions[0].ID)
	} else if len(args) >= 1 {
		idString := args[0]
		if strings.HasPrefix(idString, "<@!") && strings.HasSuffix(idString, ">") {
			idString = idString[3 : len(idString)-1]
		}
		user, err = ctx.Session.User(idString)
	} else {
		err = errors.New("that is not a valid id")
	}
	if err != nil {
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Couldn't get that user",
			Error:   err,
		}
		return
	}

	if higher, err := disgoman.HasHigherRole(ctx.Session, ctx.Guild.ID, ctx.Message.Author.ID, user.ID); err != nil {
		if err.Error() == "can't find caller member" {
			ctx.ErrorChannel <- disgoman.CommandError{
				Context: ctx,
				Message: "Who are you?",
				Error:   err,
			}
			return
		}
	} else if !higher {
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "You must have a higher role than the person you are trying to ban",
			Error:   errors.New("need higher role"),
		}
		return
	}

	if higher, err := disgoman.HasHigherRole(ctx.Session, ctx.Guild.ID, ctx.Session.State.User.ID, user.ID); err != nil {
		if err.Error() == "can't find caller member" {
			ctx.ErrorChannel <- disgoman.CommandError{
				Context: ctx,
				Message: "Who am I?",
				Error:   err,
			}
			return
		}
	} else if !higher {
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "I don't have a high enough role to ban that person",
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
	days := 7
	if ctx.Invoked == "ban-no-delete" {
		days = 0
	}
	err = ctx.Session.GuildBanCreateWithReason(ctx.Guild.ID, user.ID, auditReason, days)

	if err != nil {
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: fmt.Sprintf("Something went wrong banning %v", user.Username),
			Error:   err,
		}
		return
	}

	event := &utils.LogEvent{
		Embed: discordgo.MessageEmbed{
			Title: "User Banned",
			Description: fmt.Sprintf(
				"User %v#%v was banned by %v.\nReason: %v",
				user.Username,
				user.Discriminator,
				ctx.Message.Author.Username,
				reason),
			Timestamp: time.Now().Format("2006-01-02 15:04:05"),
			Color:     0xff0000,
		},
		GuildID: ctx.Guild.ID,
		Session: ctx.Session,
	}
	utils.LoggingChannel <- event
	_, _ = ctx.Send(fmt.Sprintf("User %v#%v has been banned.", user.Username, user.Discriminator))
}

func unbanUserCommand(ctx disgoman.Context, args []string) {
	var user *discordgo.User
	var err error
	if len(ctx.Message.Mentions) > 0 {
		user, err = ctx.Session.User(ctx.Message.Mentions[0].ID)
	} else if len(args) >= 1 {
		idString := args[0]
		if strings.HasPrefix(idString, "<@!") && strings.HasSuffix(idString, ">") {
			idString = idString[3 : len(idString)-1]
		}
		user, err = ctx.Session.User(idString)
	} else {
		err = errors.New("that is not a valid id")
	}
	if err != nil {
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Couldn't get that user",
			Error:   err,
		}
		return
	}

	bans, err := ctx.Session.GuildBans(ctx.Guild.ID)
	if err != nil {
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Error processing current bans",
			Error:   err,
		}
		return
	}
	for _, ban := range bans {
		if ban.User.ID == user.ID {
			err = ctx.Session.GuildBanDelete(ctx.Guild.ID, user.ID)
			if err != nil {
				ctx.ErrorChannel <- disgoman.CommandError{
					Context: ctx,
					Message: fmt.Sprintf("Something went wrong unbanning %v", user.Username),
					Error:   err,
				}
				return
			}
			event := &utils.LogEvent{
				Embed: discordgo.MessageEmbed{
					Title: "User Banned",
					Description: fmt.Sprintf(
						"User %v#%v was unbanned by %v.\nOrignal Ban Reason: %v",
						user.Username,
						user.Discriminator,
						ctx.Message.Author.Username,
						ban.Reason),
					Timestamp: time.Now().Format("2006-01-02 15:04:05"),
					Color:     0x00ff00,
				},
				GuildID: ctx.Guild.ID,
				Session: ctx.Session,
			}
			utils.LoggingChannel <- event
			_, _ = ctx.Send(fmt.Sprintf("User %v#%v has been unbanned.", user.Username, user.Discriminator))
			return
		}
	}
	_, _ = ctx.Send(fmt.Sprintf("%v#%v is not banned in this guild.", user.Username, user.Discriminator))
}
