package utils

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	imap "github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
)

var (
	emailUsername = os.Getenv("GOFF_EMAIL_USERNAME")
	emailPassword = os.Getenv("GOFF_EMAIL_PASSWORD")
)

var EmailClient client.Client

func RecieveEmail() {
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
		defer EmailClient.Logout()

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
				mr, err := mail.CreateReader(r)
				if err != nil {
					log.Println(err)
					continue
				}
				header := mr.Header
				if date, err := header.Date(); err == nil {
					log.Println("Date:", date)
				}
				if from, err := header.AddressList("From"); err == nil {
					log.Println("From:", from)
				}
				if to, err := header.AddressList("To"); err == nil {
					log.Println("To:", to)
				}
				if subject, err := header.Subject(); err == nil {
					log.Println("Subject:", subject)
				}
				for {
					p, err := mr.NextPart()
					if err == io.EOF {
						break
					} else if err != nil {
						log.Println(err)
						break
					}

					switch h := p.Header.(type) {
					case *mail.InlineHeader:
						// This is the message's text (can be plain-text or HTML)
						b, _ := ioutil.ReadAll(p.Body)
						log.Printf("Got text: %v\n", string(b))
					case *mail.AttachmentHeader:
						// This is an attachment
						filename, _ := h.Filename()
						log.Printf("Got attachment: %v\n", filename)
					}
				}
			}
		}
		time.Sleep(300 * time.Second)
	}
}
