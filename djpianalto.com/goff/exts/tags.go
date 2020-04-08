package exts

import (
	"djpianalto.com/goff/djpianalto.com/goff/utils"
	"fmt"
	"github.com/dustinpianalto/disgoman"
	"github.com/kballard/go-shellquote"
	"log"
	"strings"
)

func addTagCommand(ctx disgoman.Context, args []string) error {
	if len(args) >= 1 {
		args, err := shellquote.Split(strings.Join(args, " "))
		if err != nil {
			ctx.Send(err.Error())
			return err
		}
		queryString := `SELECT tags.id, tags.tag, tags.content from tags
		WHERE tags.guild_id = $1
		AND tags.tag = $2;`
		row := utils.Database.QueryRow(queryString, ctx.Guild.ID, args[0])
		var dest string
		if err := row.Scan(&dest); err != nil {
			tag := args[0]
			if tag == "" {
				ctx.Send("That is not a valid tag name")
				return nil
			}
			if len(args) <= 1 {
				ctx.Send("I got a name but no value.")
				return nil
			}
			value := args[1]
			if value == "" {
				ctx.Send("You have to include a content for the tag")
				return nil
			}
			queryString = `INSERT INTO tags (tag, content, creator, guild_id) VALUES ($1, $2, $3, $4);`
			_, err := utils.Database.Exec(queryString, tag, value, ctx.Message.Author.ID, ctx.Guild.ID)
			if err != nil {
				ctx.Send(err.Error())
				return err
			}
			ctx.Send(fmt.Sprintf("Tag %v added successfully.", tag))
			return nil
		} else {
			ctx.Send("That tag already exists.")
			return nil
		}
	} else {
		ctx.Send("You need to tell me what tag you want to add...")
		return nil
	}
}

func tagCommand(ctx disgoman.Context, args []string) error {
	if len(args) >= 1 {
		tagString := strings.Join(args, " ")
		queryString := `SELECT tags.id, tags.tag, tags.content from tags
		WHERE tags.guild_id = $1;`
		rows, err := utils.Database.Query(queryString, ctx.Guild.ID)
		if err != nil {
			ctx.Send(err.Error())
			return err
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
					return nil
				}
			}
			ctx.Send(fmt.Sprintf("Tag %v not found", args[0]))
			return nil
		}
	} else {
		ctx.Send("I need a tag to check fetch...")
		return nil
	}
}
