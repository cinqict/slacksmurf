package hello

import (
	"fmt"
	"strings"

	"github.com/cinqict/slacksmurf/cmatrix"
	"github.com/davecgh/go-spew/spew"
	"github.com/nlopes/slack"
)

var (
	returnChannel chan cmatrix.AttachmentChannel
)

const (
	prefix = "hello"
)

var help = map[string]string{
	"hello": "an average greeting",
	"help":  "display help info",
}

//Load initializes the hello plugin
func Load() {
	simpleChannel := make(chan cmatrix.HandlerChannel)
	returnChannel = cmatrix.GetReturnChannel()
	go handler(returnChannel, simpleChannel)
	cmatrix.Add(prefix, simpleChannel)
}

func handler(c chan cmatrix.AttachmentChannel, bc chan cmatrix.HandlerChannel) {
	commands := map[string]string{
		"":     "an average greeting",
		"help": "display help info",
	}

	//retrieve own handler

	var attachmentChannel cmatrix.AttachmentChannel

	for {

		botChannel := <-bc

		attachmentChannel.Channel = botChannel.Channel
		commandArray := strings.Fields(botChannel.Event.Text)
		if len(commandArray) > 2 {
			switch commandArray[2] {

			case "help":
				fmt.Println("test")
				fields := make([]slack.AttachmentField, 0)
				for k, v := range commands {
					fields = append(fields, slack.AttachmentField{
						Title: prefix + " " + k,
						Value: v,
					})

				}

				attachment := &slack.Attachment{
					Pretext: "Guru Command List",
					Color:   "#B733FF",
					Fields:  fields,
				}
				fmt.Printf("%v\n", attachment)
				attachmentChannel.Attachment = attachment
				c <- attachmentChannel

			case "hello":

				fmt.Println("hello")

				spew.Dump(botChannel)

				attachmentChannel.Attachment = &slack.Attachment{
					Text: fmt.Sprintf("<@%s> Greetings Human", botChannel.Event.User),
				}

				c <- attachmentChannel
			default:
				fmt.Println("default")
			}
		}
		if len(commandArray) == 2 {
			fmt.Println("test")
			attachmentChannel.Attachment = &slack.Attachment{
				Text: fmt.Sprintf("<@%s> Greetings Human", botChannel.Event.User),
			}

			c <- attachmentChannel

		}
	}
}

// func getHelp() slack.Attachment {
//
// }
