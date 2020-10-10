package tasks

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dustinpianalto/disgoman"
	"github.com/dustinpianalto/goff/internal/postgres"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
)

var AddReminderCommand = &disgoman.Command{
	Name:                "remind",
	Aliases:             nil,
	Description:         "Remind me at a later time",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: 0,
	Invoke:              addReminderFunc,
}

func addReminderFunc(ctx disgoman.Context, args []string) {
	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)

	text := strings.Join(args, " ")
	r, err := w.Parse(text, time.Now())
	if err != nil {
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Error parsing time",
			Error:   err,
		}
		return
	}
	if r == nil {
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "You need to include a valid time",
			Error:   errors.New("no time found"),
		}
		return
	}
	content := strings.Replace(text, r.Text+" ", "", 1)
	query := "INSERT INTO tasks (type, content, guild_id, channel_id, user_id, trigger_time) " +
		"VALUES ('Reminder', $1, $2, $3, $4, $5)"
	_, err = postgres.DB.Exec(query, content, ctx.Guild.ID, ctx.Channel.ID, ctx.User.ID, r.Time)
	if err != nil {
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Error adding task to database",
			Error:   err,
		}
		return
	}
	_ = ctx.Session.MessageReactionAdd(ctx.Channel.ID, ctx.Message.ID, "✅")
	_, _ = ctx.Session.ChannelMessageSend(
		ctx.Channel.ID,
		fmt.Sprintf("I will remind you at %v, with `%v`", r.Time.Format(time.RFC1123), content),
	)
}
