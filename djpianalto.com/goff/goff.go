package main

import (
	"djpianalto.com/goff/djpianalto.com/goff/exts"
	"djpianalto.com/goff/djpianalto.com/goff/utils"
	"fmt"
	"github.com/dustinpianalto/disgoman"
	//"github.com/MikeModder/anpan"
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
	dg.State.MaxMessageCount = 100

	utils.ConnectDatabase(os.Getenv("DATABASE_URL"))
	utils.InitializeDatabase()
	//utils.LoadTestData()

	//prefixes := []string{
	//	"Go.",
	//}
	owners := []string{
		"351794468870946827",
	}

	// Arguments are:
	// prefixes    - []string
	// owner ids   - []string
	// ignore bots - bool
	// check perms - bool
	handler := disgoman.CommandManager{
		Prefixes:         getPrefixes,
		Owners:           owners,
		StatusManager:    disgoman.GetDefaultStatusManager(),
		ErrorChannel:     make(chan disgoman.CommandError, 10),
		Commands:         make(map[string]*disgoman.Command),
		IgnoreBots:       true,
		CheckPermissions: false,
	}

	// Add Command Handlers
	exts.AddCommandHandlers(&handler)

	//if _, ok := handler.Commands["help"]; !ok {
	//	handler.AddDefaultHelpCommand()
	//}

	dg.AddHandler(handler.OnMessage)
	dg.AddHandler(handler.StatusManager.OnReady)

	err = dg.Open()
	if err != nil {
		fmt.Println("There was an error opening the connection, ", err)
		return
	}

	// Start the Error handler in a goroutine
	go ErrorHandler(handler.ErrorChannel)

	// Start the Logging handler in a goroutine
	go utils.LoggingHandler(utils.LoggingChannel)

	// Start the task handler in a goroutine
	go utils.ProcessTasks(dg, 10)

	fmt.Println("The Bot is now running.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	fmt.Println("Shutting Down...")
	err = dg.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func getPrefixes(guild_id string) []string {
	return []string{"Go.", "go."}
}

func ErrorHandler(ErrorChan chan disgoman.CommandError) {
	for ce := range ErrorChan {
		msg := ce.Message
		if msg == "" {
			msg = ce.Error.Error()
		}
		_, _ = ce.Context.Send(msg)
		fmt.Println(ce.Error)
	}
}
