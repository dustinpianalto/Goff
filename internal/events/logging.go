package events

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/dustinpianalto/goff/internal/postgres"
)

var LoggingChannel = make(chan *LogEvent, 10)

type LogEvent struct {
	// Embed with log message
	Embed discordgo.MessageEmbed
	// Guild to log event in
	GuildID string
	// Discordgo Session. Needed for sending messages
	Session *discordgo.Session
}

func LoggingHandler(lc chan *LogEvent) {
	for event := range lc {
		var channelID string
		row := postgres.DB.QueryRow("SELECT logging_channel FROM guilds where id=$1", event.GuildID)
		err := row.Scan(&channelID)
		if err != nil {
			fmt.Println(err)
			return
		}
		if channelID == "" {
			return
		}

		_, _ = event.Session.ChannelMessageSendEmbed(channelID, &event.Embed)
	}
}
