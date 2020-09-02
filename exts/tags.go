package exts

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/dustinpianalto/disgoman"
	"github.com/dustinpianalto/goff/utils"
)

func addTagCommand(ctx disgoman.Context, input []string) {
	if len(input) >= 1 {
		queryString := `SELECT tags.id, tags.tag, tags.content from tags
		WHERE tags.guild_id = $1
		AND tags.tag = $2;`
		row := utils.Database.QueryRow(queryString, ctx.Guild.ID, input[0])
		var dest string
		if err := row.Scan(&dest); err != nil {
			tag := input[0]
			if tag == "" {
				ctx.ErrorChannel <- disgoman.CommandError{
					Context: ctx,
					Message: "That is not a valid tag name",
					Error:   err,
				}
				return
			}
			if len(input) <= 1 {
				ctx.ErrorChannel <- disgoman.CommandError{
					Context: ctx,
					Message: "I got a name but no value",
					Error:   err,
				}
				return
			}
			value := strings.Join(input[1:], " ")
			if value == "" {
				ctx.ErrorChannel <- disgoman.CommandError{
					Context: ctx,
					Message: "You have to include a content for the tag",
					Error:   err,
				}
				return
			}
			queryString = `INSERT INTO tags (tag, content, creator, guild_id) VALUES ($1, $2, $3, $4);`
			_, err := utils.Database.Exec(queryString, tag, value, ctx.Message.Author.ID, ctx.Guild.ID)
			if err != nil {
				ctx.ErrorChannel <- disgoman.CommandError{
					Context: ctx,
					Message: "",
					Error:   err,
				}
				return
			}
			ctx.Send(fmt.Sprintf("Tag %v added successfully.", tag))
			return
		} else {
			ctx.ErrorChannel <- disgoman.CommandError{
				Context: ctx,
				Message: "That tag already exists",
				Error:   err,
			}
			return
		}
	} else {
		ctx.Send("You need to tell me what tag you want to add...")
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "You need to tell me what tag you want to add...",
			Error:   errors.New("nothing to do"),
		}
		return
	}
}

func tagCommand(ctx disgoman.Context, args []string) {
	if len(args) >= 1 {
		tagString := strings.Join(args, " ")
		queryString := `SELECT tags.id, tags.tag, tags.content from tags
		WHERE tags.guild_id = $1;`
		rows, err := utils.Database.Query(queryString, ctx.Guild.ID)
		if err != nil {
			ctx.ErrorChannel <- disgoman.CommandError{
				Context: ctx,
				Message: "",
				Error:   err,
			}
			return
		} else {
			for rows.Next() {
				var (
					id      int
					tag     string
					content string
				)
				if err := rows.Scan(&id, &tag, &content); err != nil {
					log.Fatal(err)
				}
				if tagString == tag {
					ctx.Send(content)
					return
				}
			}
			ctx.ErrorChannel <- disgoman.CommandError{
				Context: ctx,
				Message: fmt.Sprintf("Tag %v not found", args[0]),
				Error:   err,
			}
			return
		}
	} else {
		ctx.Send("I need a tag to fetch...")
		ctx.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "I need a tag to fetch...",
			Error:   errors.New("nothing to do"),
		}
		return
	}
}
