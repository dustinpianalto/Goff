package exts

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/dustinpianalto/disgoman"
	"github.com/dustinpianalto/goff/utils"
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

func addGuildCommand(ctx disgoman.Context, args []string) {
	var guildID string
	row := utils.Database.QueryRow("SELECT id FROM guilds where id=$1", ctx.Guild.ID)
	err := row.Scan(&guildID)
	if err == nil {
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "This guild is already in my database",
			Error:   err,
		}
		return
	}

	_, err = utils.Database.Query("INSERT INTO guilds (id) VALUES ($1)", ctx.Guild.ID)
	if err != nil {
		fmt.Println(err)
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "There was a problem inserting this guild into the database",
			Error:   err,
		}
		return
	}
	_, _ = ctx.Send("This guild has been added.")

}

func puzzleChannel(ctx disgoman.Context, args []string) {
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
		_, err := utils.Database.Exec("UPDATE guilds SET puzzle_channel='' WHERE id=$1;", ctx.Guild.ID)
		if err != nil {
			ctx.ErrorChannel <- disgoman.CommandError{
				Context: ctx,
				Message: "Error Updating Database",
				Error:   err,
			}
			return
		}
		_, _ = ctx.Send("Puzzle Channel Updated.")
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
	_, err = utils.Database.Exec("UPDATE guilds SET puzzle_channel=$1 WHERE id=$2;", idString, ctx.Guild.ID)
	if err != nil {
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Error Updating Database",
			Error:   err,
		}
		return
	}
	_, _ = ctx.Send("Puzzle Channel Updated.")
}

func getPuzzleChannel(ctx disgoman.Context, _ []string) {
	var channelID string
	row := utils.Database.QueryRow("SELECT puzzle_channel FROM guilds where id=$1", ctx.Guild.ID)
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
		_, _ = ctx.Send("The puzzle channel is not set.")
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
	_, _ = ctx.Send(fmt.Sprintf("The puzzle channel is currently %s", channel.Mention()))
	return
}

func puzzleRole(ctx disgoman.Context, args []string) {
	var idString string
	if len(args) > 0 {
		idString = args[0]
		if strings.HasPrefix(idString, "<@&") && strings.HasSuffix(idString, ">") {
			idString = idString[3 : len(idString)-1]
		}
	} else {
		idString = ""
	}
	fmt.Println(idString)
	if idString == "" {
		_, err := utils.Database.Exec("UPDATE guilds SET puzzle_role=NULL WHERE id=$1;", ctx.Guild.ID)
		if err != nil {
			ctx.ErrorChannel <- disgoman.CommandError{
				Context: ctx,
				Message: "Error Updating Database",
				Error:   err,
			}
			return
		}
		_, _ = ctx.Send("Puzzle Role Cleared.")
		return
	}
	role, err := ctx.Session.State.Role(ctx.Guild.ID, idString)
	if err != nil {
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Can't find that Role.",
			Error:   err,
		}
		return
	}
	_, err = utils.Database.Exec("INSERT INTO roles (id, guild_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", role.ID, ctx.Guild.ID)
	if err != nil {
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Error Updating Database",
			Error:   err,
		}
		return
	}
	_, err = utils.Database.Exec("UPDATE guilds SET puzzle_role=$1 WHERE id=$2;", role.ID, ctx.Guild.ID)
	if err != nil {
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Error Updating Database",
			Error:   err,
		}
		return
	}
	_, _ = ctx.Send("Puzzle Role Updated.")
}

func getPuzzleRole(ctx disgoman.Context, _ []string) {
	var roleID sql.NullString
	row := utils.Database.QueryRow("SELECT puzzle_role FROM guilds where id=$1", ctx.Guild.ID)
	err := row.Scan(&roleID)
	if err != nil {
		fmt.Println(err)
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Error getting data from the database.",
			Error:   err,
		}
		return
	}
	if !roleID.Valid {
		_, _ = ctx.Send("The puzzle role is not set.")
		return
	}
	role, err := ctx.Session.State.Role(ctx.Guild.ID, roleID.String)
	if err != nil {
		fmt.Println(err)
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "I got the role ID but it does not appear to be a valid role in this guild.",
			Error:   err,
		}
		return
	}
	_, _ = ctx.Send(fmt.Sprintf("The puzzle role is currently %s", role.Mention()))
	return
}
