package exts

import (
	"djpianalto.com/goff/djpianalto.com/goff/utils"
	"github.com/dustinpianalto/disgoman"
	"strconv"
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
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		_, _ = ctx.Send("An invalid ID was passed.")
		return err
	}
	if id == 0 {
		_, err = utils.Database.Exec("UPDATE guilds SET logging_channel=NULL WHERE id=$1;", ctx.Guild.ID)
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
	_, err = utils.Database.Exec("UPDATE guilds SET logging_channel=$1 WHERE id=$2;", id, ctx.Guild.ID)
	if err != nil {
		_, _ = ctx.Send("Error Updating Database")
		return err
	}
	_, _ = ctx.Send("Logging Channel Updated.")
	return nil
}
