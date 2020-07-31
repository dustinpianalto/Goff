package utils

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
)

type Task struct {
	ID           int64
	Type         string
	Content      string
	GuildID      string
	ChannelID    string
	UserID       string
	CreationTime time.Time
	TriggerTime  time.Time
}

func processTask(task *Task, s *discordgo.Session) {
	query := "SELECT completed, processing from tasks where id = $1"
	res, err := Database.Query(query, task.ID)
	if err != nil {
		log.Println(err)
	}
	var completed bool
	var processing bool
	res.Next()
	err = res.Scan(&completed, &processing)
	if err != nil {
		log.Println(err)
		return
	}
	if completed || processing {
		return
	}
	closeQuery := "Update tasks set completed = true where id = $1"
	processQuery := "UPDATE tasks SET processing = true WHERE id = $1"
	defer Database.Exec(closeQuery, task.ID)
	_, err = Database.Exec(processQuery, task.ID)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(fmt.Sprintf("Processing task %v", task.ID))
	guild, err := s.Guild(task.GuildID)
	if err != nil {
		log.Print(fmt.Sprintf("Can't find guild with ID %v. Canceling task %v.", task.GuildID, task.ID))
		return
	}
	channel, err := s.Channel(task.ChannelID)
	if err != nil {
		log.Print(fmt.Sprintf("Can't find channel with ID %v. Canceling task %v.", task.ChannelID, task.ID))
		return
	}
	if channel.GuildID != guild.ID {
		log.Print(fmt.Sprintf("The channel %v is not in guild %v. Canceling task %v.", channel.Name, guild.Name, task.ID))
		return
	}
	member, err := s.GuildMember(guild.ID, task.UserID)
	if err != nil {
		log.Print(fmt.Sprintf("Can't find user with ID %v in guild %v. Canceling task %v.", task.UserID, guild.Name, task.ID))
		return
	}
	if task.Type == "Reminder" {
		color := s.State.UserColor(member.User.ID, channel.ID)
		e := discordgo.MessageEmbed{
			Title:       "REMINDER",
			Description: task.Content,
			Timestamp:   task.CreationTime.Format(time.RFC3339),
			Color:       color,
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Created: ",
			},
		}
		msg := discordgo.MessageSend{
			Content: member.Mention(),
			Embed:   &e,
		}
		_, err = s.ChannelMessageSendComplex(channel.ID, &msg)
		if err != nil {
			log.Println(err)
		}
	}
	processQuery = "UPDATE tasks SET processing = false WHERE id = $1"
	_, err = Database.Exec(processQuery, task.ID)
	if err != nil {
		log.Println(err)
	}

}

func getTasksToRun() []Task {
	query := "SELECT id, type, content, guild_id, channel_id, user_id, creation_time, trigger_time " +
		"from tasks where completed is false and processing is false and trigger_time < $1"
	res, err := Database.Query(query, time.Now())
	if err != nil {
		log.Println(err)
	}
	var tasks []Task
	for res.Next() {
		var t Task
		err = res.Scan(&t.ID, &t.Type, &t.Content, &t.GuildID, &t.ChannelID, &t.UserID, &t.CreationTime, &t.TriggerTime)
		if err != nil {
			log.Println(err)
		}
		for _, task := range tasks {
			if task.ID == t.ID {
				continue
			}
		}
		tasks = append(tasks, t)
	}

	return tasks
}

func ProcessTasks(s *discordgo.Session, interval int) {
	for {
		time.Sleep(time.Duration(interval * 1e9))

		tasks := getTasksToRun()

		if len(tasks) > 0 {
			for _, t := range tasks {
				go processTask(&t, s)
			}
		}
	}
}
