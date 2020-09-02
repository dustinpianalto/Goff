package main

import (
	"fmt"
	"log"

	"github.com/dustinpianalto/disgoman"
	"github.com/dustinpianalto/goff/events"
	"github.com/dustinpianalto/goff/exts"
	"github.com/dustinpianalto/goff/utils"

	//"github.com/MikeModder/anpan"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
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
	dg.StateEnabled = true

	dg.Identify = discordgo.Identify{
		Intents: discordgo.MakeIntent(discordgo.IntentsAll),
	}

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
	dg.AddHandler(events.OnMessageUpdate)
	dg.AddHandler(events.OnMessageDelete)
	dg.AddHandler(events.OnGuildMemberAddLogging)
	dg.AddHandler(events.OnGuildMemberRemoveLogging)

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
	go utils.ProcessTasks(dg, 1)

	go utils.RecieveEmail(dg)

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

func getPrefixes(guildID string) []string {
	queryString := "Select prefix from prefixes p, x_guilds_prefixes xgp where xgp.guild_id = $1 and xgp.prefix_id = p.id"
	rows, err := utils.Database.Query(queryString, guildID)
	if err != nil {
		log.Println(err)
		return []string{"Go.", "go."}
	}
	var prefixes []string
	for rows.Next() {
		var prefix string
		err = rows.Scan(&prefix)
		if err != nil {
			log.Println(err)
			return []string{"Go.", "go."}
		}
		prefixes = append(prefixes, prefix)
	}
	return prefixes
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
