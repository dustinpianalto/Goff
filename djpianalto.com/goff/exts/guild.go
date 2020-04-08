package exts

import (
	"djpianalto.com/goff/djpianalto.com/goff/utils"
	"fmt"
	"github.com/dustinpianalto/disgoman"
	"strings"
)

// Guild management commands

func loggingChannel(ctx disgoman.Context, args []string) error {
	var idString string
	if len(args) > 0 {
		idString = args[0]
		if strings.HasPrefix(idString, "<#") && strings.HasSuffix(idString, ">") {
			idString = idString[2 : len(idString)-1]
		}
	} else {
		idString = "0"
	}
	if idString == "" {
		_, err := utils.Database.Exec("UPDATE guilds SET logging_channel=NULL WHERE id=$1;", ctx.Guild.ID)
		if err != nil {
			_, _ = ctx.Send("Error Updating Database")
		}
	}
	channel, err := ctx.Session.State.Channel(idString)
	if err != nil {
		_, _ = ctx.Send("Can't find that channel.")
		return err
	}
	if channel.GuildID != ctx.Guild.ID {
		_, _ = ctx.Send("The channel passed is not in this guild.")
		return err
	}
	_, err = utils.Database.Exec("UPDATE guilds SET logging_channel=$1 WHERE id=$2;", idString, ctx.Guild.ID)
	if err != nil {
		_, _ = ctx.Send("Error Updating Database")
		return err
	}
	_, _ = ctx.Send("Logging Channel Updated.")
	return nil
}

func getLoggingChannel(ctx disgoman.Context, _ []string) error {
	var channelID string
	row := utils.Database.QueryRow("SELECT logging_channel FROM guilds where id=$1", ctx.Guild.ID)
	err := row.Scan(channelID)
	if err != nil {
		fmt.Println(err)
		_, _ = ctx.Send("Error getting data from the database.")
		return err
	}
	if channelID == "" {
		_, _ = ctx.Send("The logging channel is not set.")
		return nil
	}
	channel, err := ctx.Session.State.GuildChannel(ctx.Guild.ID, channelID)
	if err != nil {
		fmt.Println(err)
		_, _ = ctx.Send("I got the channel ID but it does not appear to be a valid channel in this guild.")
		return err
	}
	_, _ = ctx.Send(fmt.Sprintf("The logging channel is currently %s", channel.Mention()))
	return nil
}
