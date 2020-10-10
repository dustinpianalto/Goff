package main

import (
	"fmt"
	"log"

	"github.com/dustinpianalto/disgoman"
	"github.com/dustinpianalto/goff/internal/exts"
	"github.com/dustinpianalto/goff/internal/exts/guild_management"
	"github.com/dustinpianalto/goff/internal/exts/logging"
	"github.com/dustinpianalto/goff/internal/exts/tasks"
	"github.com/dustinpianalto/goff/internal/exts/user_management"
	"github.com/dustinpianalto/goff/internal/postgres"
	"github.com/dustinpianalto/goff/internal/services"
	"github.com/dustinpianalto/goff/pkg/email"

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

	postgres.ConnectDatabase(os.Getenv("DATABASE_URL"))
	postgres.InitializeDatabase()
	//utils.LoadTestData()

	us := &postgres.UserService{DB: postgres.DB}
	gs := &postgres.GuildService{DB: postgres.DB}

	//prefixes := []string{
	//	"Go.",
	//}
	owners := []string{
		"351794468870946827",
	}

	manager := disgoman.CommandManager{
		Prefixes:         getPrefixes,
		Owners:           owners,
		StatusManager:    disgoman.GetDefaultStatusManager(),
		ErrorChannel:     make(chan disgoman.CommandError, 10),
		Commands:         make(map[string]*disgoman.Command),
		IgnoreBots:       true,
		CheckPermissions: false,
	}

	// Add Command Handlers
	exts.AddCommandHandlers(&manager)
	services.InitalizeServices(us, gs)

	//if _, ok := handler.Commands["help"]; !ok {
	//	handler.AddDefaultHelpCommand()
	//}

	dg.AddHandler(manager.OnMessage)
	dg.AddHandler(manager.StatusManager.OnReady)
	dg.AddHandler(guild_management.OnMessageUpdate)
	dg.AddHandler(guild_management.OnMessageDelete)
	dg.AddHandler(user_management.OnGuildMemberAddLogging)
	dg.AddHandler(user_management.OnGuildMemberRemoveLogging)

	err = dg.Open()
	if err != nil {
		fmt.Println("There was an error opening the connection, ", err)
		return
	}

	// Start the Error handler in a goroutine
	go ErrorHandler(manager.ErrorChannel)

	// Start the Logging handler in a goroutine
	go logging.LoggingHandler(logging.LoggingChannel)

	// Start the task handler in a goroutine
	go tasks.ProcessTasks(dg, 1)

	go email.RecieveEmail(dg)

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
	rows, err := postgres.DB.Query(queryString, guildID)
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
	if len(prefixes) == 0 {
		prefixes = append(prefixes, "Godev.", "godev.")
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
