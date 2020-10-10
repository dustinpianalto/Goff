package tags

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/dustinpianalto/disgoman"
	"github.com/dustinpianalto/goff"
	"github.com/dustinpianalto/goff/internal/postgres"
	"github.com/dustinpianalto/goff/internal/services"
)

var AddTagCommand = &disgoman.Command{
	Name:                "addtag",
	Aliases:             nil,
	Description:         "Add a tag",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: 0,
	SanitizeEveryone:    true,
	Invoke:              addTagCommandFunc,
}

func addTagCommandFunc(ctx disgoman.Context, input []string) {
	if len(input) >= 1 {
		queryString := `SELECT tags.id, tags.tag, tags.content from tags
		WHERE tags.guild_id = $1
		AND tags.tag = $2;`
		row := postgres.DB.QueryRow(queryString, ctx.Guild.ID, input[0])
		var dest string
		if err := row.Scan(&dest); err != nil {
			tag := input[0]
			if tag == "" {
				ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
					Context: ctx,
					Message: "That is not a valid tag name",
					Error:   err,
				}
				return
			}
			if len(input) <= 1 {
				ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
					Context: ctx,
					Message: "I got a name but no value",
					Error:   err,
				}
				return
			}
			value := strings.Join(input[1:], " ")
			if value == "" {
				ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
					Context: ctx,
					Message: "You have to include a content for the tag",
					Error:   err,
				}
				return
			}
			err = services.UserService.CreateUser(&goff.User{ID: ctx.User.ID})
			if err != nil {
				log.Printf("Error creating user %s: %s", ctx.User.ID, err.Error())
			}
			queryString = `INSERT INTO tags (tag, content, creator, guild_id) VALUES ($1, $2, $3, $4);`
			_, err := postgres.DB.Exec(queryString, tag, value, ctx.Message.Author.ID, ctx.Guild.ID)
			if err != nil {
				ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
					Context: ctx,
					Message: "",
					Error:   err,
				}
				return
			}
			ctx.Send(fmt.Sprintf("Tag %v added successfully.", tag))
			return
		} else {
			ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
				Context: ctx,
				Message: "That tag already exists",
				Error:   err,
			}
			return
		}
	} else {
		ctx.Send("You need to tell me what tag you want to add...")
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "You need to tell me what tag you want to add...",
			Error:   errors.New("nothing to do"),
		}
		return
	}
}

var TagCommand = &disgoman.Command{
	Name:                "tag",
	Aliases:             nil,
	Description:         "Get a tag",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: 0,
	Invoke:              tagCommandFunc,
}

func tagCommandFunc(ctx disgoman.Context, args []string) {
	if len(args) >= 1 {
		tagString := strings.Join(args, " ")
		queryString := `SELECT tags.id, tags.tag, tags.content from tags
		WHERE tags.guild_id = $1;`
		rows, err := postgres.DB.Query(queryString, ctx.Guild.ID)
		if err != nil {
			ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
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
					log.Println(err)
				}
				if tagString == tag {
					ctx.Send(content)
					return
				}
			}
			ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
				Context: ctx,
				Message: fmt.Sprintf("Tag %v not found", args[0]),
				Error:   err,
			}
			return
		}
	} else {
		ctx.Send("I need a tag to fetch...")
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "I need a tag to fetch...",
			Error:   errors.New("nothing to do"),
		}
		return
	}
}
