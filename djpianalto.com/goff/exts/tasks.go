package exts

import (
	"djpianalto.com/goff/djpianalto.com/goff/utils"
	"errors"
	"github.com/dustinpianalto/disgoman"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
	"strings"
	"time"
)

func addReminderCommand(ctx disgoman.Context, args []string) {
	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)

	text := strings.Join(args, " ")
	r, err := w.Parse(text, time.Now())
	if err != nil {
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Error parsing time",
			Error:   err,
		}
		return
	}
	if r == nil {
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "You need to include a valid time",
			Error:   errors.New("no time found"),
		}
		return
	}
	content := strings.Replace(text, r.Text, "", 1)
	query := "INSERT INTO tasks (type, content, guild_id, channel_id, user_id, trigger_time) " +
		"VALUES ('Reminder', $1, $2, $3, $4, $5)"
	utils.Database.Exec(query, content, ctx.Guild.ID, ctx.Channel.ID, ctx.User.ID, r.Time)
}
