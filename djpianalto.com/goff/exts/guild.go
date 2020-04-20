package exts

import (
	"djpianalto.com/goff/djpianalto.com/goff/utils"
	"fmt"
	"github.com/dustinpianalto/disgoman"
	"strings"
)

// Guild management commands

func loggingChannel(ctx disgoman.Context, args []string) {
	var idString string
	if len(args) > 0 {
		idString = args[0]
		if strings.HasPrefix(idString, "<#") && strings.HasSuffix(idString, ">") {
			idString = idString[2 : len(idString)-1]
		}
	} else {
		idString = ""
	}
	fmt.Println(idString)
	if idString == "" {
		_, err := utils.Database.Exec("UPDATE guilds SET logging_channel='' WHERE id=$1;", ctx.Guild.ID)
		if err != nil {
			ctx.ErrorChannel <- disgoman.CommandError{
				Context: ctx,
				Message: "Error Updating Database",
				Error:   err,
			}
			return
		}
		_, _ = ctx.Send("Logging Channel Updated.")
		return
	}
	channel, err := ctx.Session.State.Channel(idString)
	if err != nil {
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Can't find that channel.",
			Error:   err,
		}
		return
	}
	if channel.GuildID != ctx.Guild.ID {
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "The channel passed is not in this guild.",
			Error:   err,
		}
		return
	}
	_, err = utils.Database.Exec("UPDATE guilds SET logging_channel=$1 WHERE id=$2;", idString, ctx.Guild.ID)
	if err != nil {
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Error Updating Database",
			Error:   err,
		}
		return
	}
	_, _ = ctx.Send("Logging Channel Updated.")
}

func getLoggingChannel(ctx disgoman.Context, _ []string) {
	var channelID string
	row := utils.Database.QueryRow("SELECT logging_channel FROM guilds where id=$1", ctx.Guild.ID)
	err := row.Scan(&channelID)
	if err != nil {
		fmt.Println(err)
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Error getting data from the database.",
			Error:   err,
		}
		return
	}
	if channelID == "" {
		_, _ = ctx.Send("The logging channel is not set.")
		return
	}
	channel, err := ctx.Session.State.GuildChannel(ctx.Guild.ID, channelID)
	if err != nil {
		fmt.Println(err)
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "I got the channel ID but it does not appear to be a valid channel in this guild.",
			Error:   err,
		}
		return
	}
	_, _ = ctx.Send(fmt.Sprintf("The logging channel is currently %s", channel.Mention()))
	return
}

func welcomeChannel(ctx disgoman.Context, args []string) {
	var idString string
	if len(args) > 0 {
		idString = args[0]
		if strings.HasPrefix(idString, "<#") && strings.HasSuffix(idString, ">") {
			idString = idString[2 : len(idString)-1]
		}
	} else {
		idString = ""
	}
	fmt.Println(idString)
	if idString == "" {
		_, err := utils.Database.Exec("UPDATE guilds SET welcome_channel='' WHERE id=$1;", ctx.Guild.ID)
		if err != nil {
			ctx.ErrorChannel <- disgoman.CommandError{
				Context: ctx,
				Message: "Error Updating Database",
				Error:   err,
			}
			return
		}
		_, _ = ctx.Send("Welcomer Disabled.")
		return
	}
	channel, err := ctx.Session.State.Channel(idString)
	if err != nil {
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Can't find that channel.",
			Error:   err,
		}
		return
	}
	if channel.GuildID != ctx.Guild.ID {
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "The channel passed is not in this guild.",
			Error:   err,
		}
		return
	}
	_, err = utils.Database.Exec("UPDATE guilds SET welcome_channel=$1 WHERE id=$2;", idString, ctx.Guild.ID)
	if err != nil {
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Error Updating Database",
			Error:   err,
		}
		return
	}
	_, _ = ctx.Send("Welcome Channel Updated.")
	return
}

func getWelcomeChannel(ctx disgoman.Context, _ []string) {
	var channelID string
	row := utils.Database.QueryRow("SELECT welcome_channel FROM guilds where id=$1", ctx.Guild.ID)
	err := row.Scan(&channelID)
	if err != nil {
		fmt.Println(err)
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Error getting data from the database.",
			Error:   err,
		}
		return
	}
	if channelID == "" {
		_, _ = ctx.Send("The welcomer is disabled.")
		return
	}
	channel, err := ctx.Session.State.GuildChannel(ctx.Guild.ID, channelID)
	if err != nil {
		fmt.Println(err)
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "I got the channel ID but it does not appear to be a valid channel in this guild.",
			Error:   err,
		}
		return
	}
	_, _ = ctx.Send(fmt.Sprintf("The welcome channel is currently %s", channel.Mention()))
}
