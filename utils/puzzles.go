package utils

import (
	"io"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/emersion/go-message/mail"
)

func ProcessPuzzleEmail(mr *mail.Reader, dg *discordgo.Session) {
	var body []byte
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println(err)
			break
		}

		switch h := p.Header.(type) {
		case *mail.InlineHeader:
			// This is the message's text (can be plain-text or HTML)
			if t, _, _ := h.ContentType(); t == "text/plain" {
				body, _ = ioutil.ReadAll(p.Body)
				break
			}
		}
	}
	if len(body) > 0 {
		s := string(body)
		puzzle := strings.Split(s, "----------")[0]
		date, err := mr.Header.Date()
		if err != nil {
			log.Println(err)
			return
		}
		e := discordgo.MessageEmbed{
			Title:       "Daily Coding Problem",
			URL:         "https://dailycodingproblem.com/",
			Description: "```" + puzzle + "```",
			Timestamp:   date.Format(time.RFC3339),
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Daily Coding Problem",
			},
		}
		var guilds []Guild
		queryString := `SELECT id, puzzle_channel, puzzle_role from guilds`
		rows, err := Database.Query(queryString)
		if err != nil {
			log.Println(err)
		}
		for rows.Next() {
			var guild Guild
			err := rows.Scan(&guild.ID, &guild.PuzzleChannel, &guild.PuzzleRole)
			if err != nil {
				log.Println(err)
				continue
			}
			guilds = append(guilds, guild)
		}
		var puzzleID int64
		queryString = "INSERT INTO puzzles (text, time) VALUES ($1, $2) RETURNING id"
		err = Database.QueryRow(queryString, puzzle, date).Scan(&puzzleID)
		if err != nil {
			log.Println(err)
			return
		}
		for _, g := range guilds {
			if g.PuzzleChannel == "" {
				continue
			}
			var msg discordgo.MessageSend
			role, err := dg.State.Role(g.ID, g.PuzzleRole.String)
			if err != nil {
				msg = discordgo.MessageSend{
					Embed: &e,
				}
			} else {
				msg = discordgo.MessageSend{
					Content: role.Mention(),
					Embed:   &e,
				}
			}
			m, err := dg.ChannelMessageSendComplex(g.PuzzleChannel, &msg)
			if err != nil {
				log.Println(err)
			}
			queryString = "INSERT INTO x_guilds_puzzles (guild_id, puzzle_id, message_id) VALUES ($1, $2, $3)"
			_, err = Database.Exec(queryString, g.ID, puzzleID, m.ID)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}

}
