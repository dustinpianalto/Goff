package guild_management

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/dustinpianalto/disgoman"
	"github.com/dustinpianalto/goff/internal/postgres"
	"github.com/dustinpianalto/goff/internal/services"
)

// Guild management commands
var SetLoggingChannelCommand = &disgoman.Command{
	Name:                "set-logging-channel",
	Aliases:             []string{"slc"},
	Description:         "Set the channel logging messages will be sent to.",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: disgoman.PermissionManageServer,
	Invoke:              setLoggingChannelFunc,
}

func setLoggingChannelFunc(ctx disgoman.Context, args []string) {
	var idString string
	if len(args) > 0 {
		idString = args[0]
		if strings.HasPrefix(idString, "<#") && strings.HasSuffix(idString, ">") {
			idString = idString[2 : len(idString)-1]
		}
	} else {
		idString = ""
	}
	guild, err := services.GuildService.Guild(ctx.Guild.ID)
	if err != nil {
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Error Updating Database",
			Error:   err,
		}
		return
	}
	if idString == "" {
		guild.LoggingChannel = idString
		err = services.GuildService.UpdateGuild(guild)
		if err != nil {
			ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
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
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Can't find that channel.",
			Error:   err,
		}
		return
	}
	if channel.GuildID != ctx.Guild.ID {
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "The channel passed is not in this guild.",
			Error:   err,
		}
		return
	}
	guild.LoggingChannel = channel.ID
	err = services.GuildService.UpdateGuild(guild)
	if err != nil {
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Error Updating Database",
			Error:   err,
		}
		return
	}
	_, _ = ctx.Send("Logging Channel Updated.")
}

var GetLoggingChannelCommand = &disgoman.Command{
	Name:                "get-logging-channel",
	Aliases:             []string{"glc"},
	Description:         "Gets the channel logging messages will be sent to.",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: disgoman.PermissionManageServer,
	Invoke:              getLoggingChannelFunc,
}

func getLoggingChannelFunc(ctx disgoman.Context, _ []string) {
	var channelID string
	row := postgres.DB.QueryRow("SELECT logging_channel FROM guilds where id=$1", ctx.Guild.ID)
	err := row.Scan(&channelID)
	if err != nil {
		fmt.Println(err)
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
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
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "I got the channel ID but it does not appear to be a valid channel in this guild.",
			Error:   err,
		}
		return
	}
	_, _ = ctx.Send(fmt.Sprintf("The logging channel is currently %s", channel.Mention()))
	return
}

var SetWelcomeChannelCommand = &disgoman.Command{
	Name:                "set-welcome-channel",
	Aliases:             []string{"swc"},
	Description:         "Set the channel welcome messages will be sent to.",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: disgoman.PermissionManageServer,
	Invoke:              setWelcomeChannelFunc,
}

func setWelcomeChannelFunc(ctx disgoman.Context, args []string) {
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
		_, err := postgres.DB.Exec("UPDATE guilds SET welcome_channel='' WHERE id=$1;", ctx.Guild.ID)
		if err != nil {
			ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
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
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Can't find that channel.",
			Error:   err,
		}
		return
	}
	if channel.GuildID != ctx.Guild.ID {
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "The channel passed is not in this guild.",
			Error:   err,
		}
		return
	}
	_, err = postgres.DB.Exec("UPDATE guilds SET welcome_channel=$1 WHERE id=$2;", idString, ctx.Guild.ID)
	if err != nil {
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Error Updating Database",
			Error:   err,
		}
		return
	}
	_, _ = ctx.Send("Welcome Channel Updated.")
	return
}

var GetWelcomeChannelCommand = &disgoman.Command{
	Name:                "get-welcome-channel",
	Aliases:             []string{"gwc"},
	Description:         "Gets the channel welcome messages will be sent to.",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: disgoman.PermissionManageServer,
	Invoke:              getWelcomeChannelFunc,
}

func getWelcomeChannelFunc(ctx disgoman.Context, _ []string) {
	var channelID string
	row := postgres.DB.QueryRow("SELECT welcome_channel FROM guilds where id=$1", ctx.Guild.ID)
	err := row.Scan(&channelID)
	if err != nil {
		fmt.Println(err)
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
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
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "I got the channel ID but it does not appear to be a valid channel in this guild.",
			Error:   err,
		}
		return
	}
	_, _ = ctx.Send(fmt.Sprintf("The welcome channel is currently %s", channel.Mention()))
}

var AddGuildCommand = &disgoman.Command{
	Name:                "addGuild",
	Aliases:             nil,
	Description:         "Adds the current guild to the database",
	OwnerOnly:           true,
	Hidden:              false,
	RequiredPermissions: 0,
	Invoke:              addGuildCommandFunc,
}

func addGuildCommandFunc(ctx disgoman.Context, args []string) {
	var guildID string
	row := postgres.DB.QueryRow("SELECT id FROM guilds where id=$1", ctx.Guild.ID)
	err := row.Scan(&guildID)
	if err == nil {
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "This guild is already in my database",
			Error:   err,
		}
		return
	}

	_, err = postgres.DB.Query("INSERT INTO guilds (id) VALUES ($1)", ctx.Guild.ID)
	if err != nil {
		fmt.Println(err)
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "There was a problem inserting this guild into the database",
			Error:   err,
		}
		return
	}
	_, _ = ctx.Send("This guild has been added.")

}

var SetPuzzleChannelCommand = &disgoman.Command{
	Name:                "set-puzzle-channel",
	Aliases:             []string{"spc"},
	Description:         "Set the channel puzzle messages will be sent to.",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: disgoman.PermissionManageServer,
	Invoke:              setPuzzleChannelFunc,
}

func setPuzzleChannelFunc(ctx disgoman.Context, args []string) {
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
		_, err := postgres.DB.Exec("UPDATE guilds SET puzzle_channel='' WHERE id=$1;", ctx.Guild.ID)
		if err != nil {
			ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
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
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Can't find that channel.",
			Error:   err,
		}
		return
	}
	if channel.GuildID != ctx.Guild.ID {
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "The channel passed is not in this guild.",
			Error:   err,
		}
		return
	}
	_, err = postgres.DB.Exec("UPDATE guilds SET puzzle_channel=$1 WHERE id=$2;", idString, ctx.Guild.ID)
	if err != nil {
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Error Updating Database",
			Error:   err,
		}
		return
	}
	_, _ = ctx.Send("Puzzle Channel Updated.")
}

var GetPuzzleChannelCommand = &disgoman.Command{
	Name:                "get-puzzle-channel",
	Aliases:             []string{"gpc"},
	Description:         "Gets the channel puzzle messages will be sent to.",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: disgoman.PermissionManageServer,
	Invoke:              getPuzzleChannelFunc,
}

func getPuzzleChannelFunc(ctx disgoman.Context, _ []string) {
	var channelID string
	row := postgres.DB.QueryRow("SELECT puzzle_channel FROM guilds where id=$1", ctx.Guild.ID)
	err := row.Scan(&channelID)
	if err != nil {
		fmt.Println(err)
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
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
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "I got the channel ID but it does not appear to be a valid channel in this guild.",
			Error:   err,
		}
		return
	}
	_, _ = ctx.Send(fmt.Sprintf("The puzzle channel is currently %s", channel.Mention()))
	return
}

var SetPuzzleRoleCommand = &disgoman.Command{
	Name:                "set-puzzle-role",
	Aliases:             []string{"spr"},
	Description:         "Set the role to be pinged when there is a new puzzle",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: disgoman.PermissionManageServer,
	Invoke:              setPuzzleRoleFunc,
}

func setPuzzleRoleFunc(ctx disgoman.Context, args []string) {
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
		_, err := postgres.DB.Exec("UPDATE guilds SET puzzle_role=NULL WHERE id=$1;", ctx.Guild.ID)
		if err != nil {
			ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
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
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Can't find that Role.",
			Error:   err,
		}
		return
	}
	_, err = postgres.DB.Exec("INSERT INTO roles (id, guild_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", role.ID, ctx.Guild.ID)
	if err != nil {
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Error Updating Database",
			Error:   err,
		}
		return
	}
	_, err = postgres.DB.Exec("UPDATE guilds SET puzzle_role=$1 WHERE id=$2;", role.ID, ctx.Guild.ID)
	if err != nil {
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Error Updating Database",
			Error:   err,
		}
		return
	}
	_, _ = ctx.Send("Puzzle Role Updated.")
}

var GetPuzzleRoleCommand = &disgoman.Command{
	Name:                "get-puzzle-role",
	Aliases:             []string{"gpr"},
	Description:         "Get the role that will be pinged when there is a new puzzle",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: disgoman.PermissionManageServer,
	Invoke:              getPuzzleRoleFunc,
}

func getPuzzleRoleFunc(ctx disgoman.Context, _ []string) {
	var roleID sql.NullString
	row := postgres.DB.QueryRow("SELECT puzzle_role FROM guilds where id=$1", ctx.Guild.ID)
	err := row.Scan(&roleID)
	if err != nil {
		fmt.Println(err)
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
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
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "I got the role ID but it does not appear to be a valid role in this guild.",
			Error:   err,
		}
		return
	}
	_, _ = ctx.Send(fmt.Sprintf("The puzzle role is currently %s", role.Mention()))
	return
}
