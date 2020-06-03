package events

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func OnMessageUpdate(session *discordgo.Session, m *discordgo.MessageUpdate) {
	fmt.Println(m.Content)
	fmt.Println(m.BeforeUpdate.Content)
}
