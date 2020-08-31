package utils

import (
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	imap "github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
)

const ()

var (
	emailUsername = os.Getenv("GOFF_EMAIL_USERNAME")
	emailPassword = os.Getenv("GOFF_EMAIL_PASSWORD")
	puzzleAddress = mail.Address{
		Name:    "Daily Coding Problem",
		Address: "founders@dailycodingproblem.com",
	}
)

var EmailClient client.Client

func RecieveEmail(dg *discordgo.Session) {
	for {
		log.Println("Connecting to Email server.")

		EmailClient, err := client.DialTLS("mail.djpianalto.com:993", nil)
		if err != nil {
			log.Println(err)
			return
		}
		if err = EmailClient.Login(emailUsername, emailPassword); err != nil {
			log.Println(err)
			return
		}
		log.Println("Connected to Email server.")

		mbox, err := EmailClient.Select("INBOX", false)
		if err != nil {
			log.Println(err)
			return
		}

		if mbox.Messages == 0 {
			log.Println("No Messages in Mailbox")
		}

		criteria := imap.NewSearchCriteria()
		criteria.WithoutFlags = []string{"\\Seen"}
		uids, err := EmailClient.Search(criteria)
		if err != nil {
			log.Println(err)
		}
		if len(uids) > 0 {
			seqset := new(imap.SeqSet)
			seqset.AddNum(uids...)
			section := &imap.BodySectionName{}
			items := []imap.FetchItem{section.FetchItem()}
			messages := make(chan *imap.Message, 10)
			go func() {
				if err = EmailClient.Fetch(seqset, items, messages); err != nil {
					log.Println(err)
					return
				}
			}()

			var wg sync.WaitGroup

			for msg := range messages {
				if msg == nil {
					log.Println("No New Messages")
					continue
				}
				r := msg.GetBody(section)
				if r == nil {
					log.Println("Server didn't send a message body")
					continue
				}
				wg.Add(1)
				go processEmail(r, dg, &wg)
			}
			wg.Wait()
		}

		EmailClient.Logout()
		time.Sleep(300 * time.Second)
	}
}

func processEmail(r io.Reader, dg *discordgo.Session, wg *sync.WaitGroup) {
	defer wg.Done()
	mr, err := mail.CreateReader(r)
	if err != nil {
		log.Println(err)
		return
	}
	header := mr.Header
	from, err := header.AddressList("From")
	if err != nil {
		log.Println(err)
		return
	}
	subject, err := header.Subject()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(from)
	log.Println(subject)
	if addressIn(from, puzzleAddress) &&
		strings.Contains(subject, "Daily Coding Problem:") {
		log.Println("Processing Puzzle")
		ProcessPuzzleEmail(mr, dg)
	}

}

func addressIn(s []*mail.Address, a mail.Address) bool {
	for _, item := range s {
		if item.String() == a.String() {
			return true
		}
	}
	return false
}
