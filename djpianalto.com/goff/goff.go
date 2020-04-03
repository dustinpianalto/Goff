package main

import (
	"djpianalto.com/goff/djpianalto.com/goff/exts"
	"fmt"
	"github.com/MikeModder/anpan"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
)

var (
	Token string
)

//func init() {
//	flag.StringVar(&Token, "t", "", "Bot Token")
//	flag.Parse()
//}

func main() {
	Token = os.Getenv("DISCORDGO_TOKEN")
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("There was an error when creating the Discord Session, ", err)
		return
	}

	prefixes := []string{
		"Go.",
	}
	owners := []string{
		"351794468870946827",
	}

	// Arguments are:
	// prefixes    - []string
	// owner ids   - []string
	// ignore bots - bool
	// check perms - bool
	handler := anpan.NewCommandHandler(prefixes, owners, true, true)

	// Add Command Handlers
	exts.AddCommandHandlers(&handler)

	if _, ok := handler.Commands["help"]; !ok {
		handler.AddDefaultHelpCommand()
	}

	dg.AddHandler(handler.OnMessage)
	dg.AddHandler(handler.StatusHandler.OnReady)

	err = dg.Open()
	if err != nil {
		fmt.Println("There was an error opening the connection, ", err)
		return
	}

	fmt.Println("The Bot is now running. Press Ctrl+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	fmt.Println("Shutting Down...")
	err = dg.Close()
	if err != nil {
		fmt.Println(err)
	}
}
