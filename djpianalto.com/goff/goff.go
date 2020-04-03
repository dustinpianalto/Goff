package main

import (
	"fmt"
	"github.com/MikeModder/anpan"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
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

	// Arguments:
	// name - command name - string
	// desc - command description - string
	// owneronly - only allow owners to run - bool
	// hidden - hide command from non-owners - bool
	// perms - permissisions required - anpan.Permission (int)
	// type - command type, sets where the command is available
	// run - function to run - func(anpan.Context, []string) / CommandRunFunc
	handler.AddCommand("ping", "Check the bot's ping", false, false, 0, anpan.CommandTypeEverywhere, pingCommand)
	handler.AddCommand("say", "Repeat a message", false, false, 0, anpan.CommandTypeEverywhere, sayCommand)
	handler.AddCommand("user", "Show info about a user", false, false, 0, anpan.CommandTypeEverywhere, userCommand)

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
	dg.Close()
}

func pingCommand(ctx anpan.Context, _ []string) error {
	timeBefore := time.Now()
	msg, _ := ctx.Reply("Pong!")
	took := time.Now().Sub(timeBefore)
	ctx.Session.ChannelMessageEdit(ctx.Message.ChannelID, msg.ID, fmt.Sprintf("Pong!\nPing Took **%s**", took.String()))
	return nil
}

func sayCommand(ctx anpan.Context, args []string) error {
	resp := strings.Join(args, " ")
	resp = strings.ReplaceAll(resp, "@everyone", "@\ufff0everyone")
	resp = strings.ReplaceAll(resp, "@here", "@\ufff0here")
	ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, resp)
	return nil
}

func userCommand(ctx anpan.Context, args []string) error {
	var member *discordgo.Member
	if len(args) == 0 {
		member, _ = ctx.Session.GuildMember(ctx.Guild.ID, ctx.Message.Author.ID)
	} else {
		var err error
		if len(ctx.Message.Mentions) > 0 {
			member, err = ctx.Session.GuildMember(ctx.Guild.ID, ctx.Message.Mentions[0].ID)
		} else {
			member, err = ctx.Session.GuildMember(ctx.Guild.ID, args[0])
		}
		if err != nil {
			return err
		}
	}
	thumb := &discordgo.MessageEmbedThumbnail{
		URL: member.User.AvatarURL(""),
	}

	var botString string
	if member.User.Bot {
		botString = "BOT"
	} else {
		botString = ""
	}

	var roles []*discordgo.Role
	for _, roleID := range member.Roles {
		role, _ := ctx.Session.State.Role(ctx.Guild.ID, roleID)
		roles = append(roles, role)
	}
	sort.Slice(roles, func(i, j int) bool { return roles[i].Position > roles[j].Position })
	var roleMentions []string
	for _, role := range roles {
		roleMentions = append(roleMentions, role.Mention())
	}
	rolesString := strings.Join(roleMentions, " ")

	rolesField := &discordgo.MessageEmbedField{
		Name:   "Roles:",
		Value:  rolesString,
		Inline: false,
	}

	guildJoinTime, _ := member.JoinedAt.Parse()
	guildJoinedField := &discordgo.MessageEmbedField{
		Name:   "Joined Guild:",
		Value:  parseDateString(guildJoinTime),
		Inline: false,
	}

	int64ID, _ := strconv.ParseInt(member.User.ID, 10, 64)
	s := parseSnowflake(int64ID)
	discordJoinedField := &discordgo.MessageEmbedField{
		Name:   "Joined Discord:",
		Value:  parseDateString(s.CreationTime),
		Inline: false,
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("%v#%v  %v", member.User.Username, member.User.Discriminator, botString),
		Description: fmt.Sprintf("**%v** (%v)", member.Nick, member.User.ID),
		Color:       0,
		Thumbnail:   thumb,
		Fields: []*discordgo.MessageEmbedField{
			guildJoinedField,
			discordJoinedField,
			rolesField,
		},
	}
	ctx.Session.ChannelMessageSendEmbed(ctx.Channel.ID, embed)
	return nil
}

func parseDateString(inTime time.Time) string {
	d := time.Now().Sub(inTime)
	s := int64(d.Seconds())
	days := s / 86400
	s = s - (days * 86400)
	hours := s / 3600
	s = s - (hours * 3600)
	minutes := s / 60
	seconds := s - (minutes * 60)
	dateString := ""
	if days != 0 {
		dateString += fmt.Sprintf("%v days ", days)
	}
	if hours != 0 {
		dateString += fmt.Sprintf("%v hours ", hours)
	}
	if minutes != 0 {
		dateString += fmt.Sprintf("%v minutes ", minutes)
	}
	if seconds != 0 {
		dateString += fmt.Sprintf("%v seconds ", seconds)
	}
	if dateString != "" {
		dateString += " ago."
	} else {
		dateString = "Now"
	}
	stamp := inTime.Format("2006-01-02 15:04:05")
	return fmt.Sprintf("%v\n%v", dateString, stamp)
}

type Snowflake struct {
	CreationTime time.Time
	WorkerID     int8
	ProcessID    int8
	Increment    int16
}

func parseSnowflake(s int64) Snowflake {
	const (
		DISCORD_EPOCH   = 1420070400000
		TIME_BITS_LOC   = 22
		WORKER_ID_LOC   = 17
		WORKER_ID_MASK  = 0x3E0000
		PROCESS_ID_LOC  = 12
		PROCESS_ID_MASK = 0x1F000
		INCREMENT_MASK  = 0xFFF
	)
	creationTime := time.Unix(((s>>TIME_BITS_LOC)+DISCORD_EPOCH)/1000.0, 0)
	workerID := (s & WORKER_ID_MASK) >> WORKER_ID_LOC
	processID := (s & PROCESS_ID_MASK) >> PROCESS_ID_LOC
	increment := s & INCREMENT_MASK
	return Snowflake{
		CreationTime: creationTime,
		WorkerID:     int8(workerID),
		ProcessID:    int8(processID),
		Increment:    int16(increment),
	}
}
