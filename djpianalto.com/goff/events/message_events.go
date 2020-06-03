package events

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func OnMessageEdit(session *discordgo.Session, m *discordgo.MessageEdit) {
	fmt.Println(m.Content)
}
