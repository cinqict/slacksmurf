package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/cinqict/slacksmurf/cmatrix"
	"github.com/cinqict/slacksmurf/plugins/hello"
	"github.com/nlopes/slack"
	"github.com/spf13/viper"
)

type BotCentral struct {
	Channel *slack.Channel
	Event   *slack.MessageEvent
	UserID  string
}

var (
	//BotID the userid for our bot
	botID string
	api   *slack.Client
)

func main() {
	// set maxprocs to the number of cpus we are running with
	runtime.GOMAXPROCS(runtime.NumCPU())

	// let's get viper to read the config file
	initializeConfig()

	// open the logfiles
	l := openLogfile(viper.GetString("log_file"))
	m := openLogfile(viper.GetString("message_log"))

	// setup the logger
	logger := log.New(l, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	message_logger := log.New(m, "slack-message", log.Lshortfile|log.LstdFlags)

	// instantiate the new slack interface
	api = slack.New(viper.GetString("slack_token"))

	// assing a logger to it
	slack.SetLogger(logger)

	//set debugging according to the specification in the config file
	api.SetDebug(viper.GetBool("slack_debug"))

	//get a new real time messaging connection
	rtm := api.NewRTM()

	botReplyChannel := cmatrix.GetReturnChannel()

	go rtm.ManageConnection()

	go handleBotReply(botReplyChannel)

	hello.Load()
	for msg := range rtm.IncomingEvents {
		fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)
			// Replace #general with your Channel ID
			rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "#general"))
			botID = ev.Info.User.ID

		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev)
			message_logger.Printf("%v\n", ev)

			channelInfo, err := api.GetChannelInfo(ev.Channel)
			if err != nil {
				log.Fatalln(err)
			}

			botCentral := &cmatrix.HandlerChannel{
				Channel: channelInfo,
				Event:   ev,
				UserID:  ev.User,
			}

			if ev.Type == "message" && strings.HasPrefix(ev.Text, "<@"+botID+">") {
				commandArray := strings.Fields(ev.Text)
				if len(commandArray) > 1 {
					hc, err := cmatrix.CGetByP(commandArray[1])
					if err == nil {
						fmt.Printf("%v\n", commandArray[1])
						hc <- *botCentral
					}
				}

			}

		case *slack.PresenceChangeEvent:
			fmt.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			fmt.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:

			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}

func initializeConfig() {
	// get input from config files

	// configfile name is barney
	viper.SetConfigName("slacksmurf")

	// add the filepaths that will be used
	viper.AddConfigPath("/etc/slacksmurf/")
	viper.AddConfigPath("$HOME/.slacksmurf")
	viper.AddConfigPath(".")

	// Handle errors reading the config file
	err := viper.ReadInConfig()
	if err != nil {
		//debugInfo := viper.Debug()
		fmt.Printf("Fatal error config file: %s \n", err)
	}

}

func openLogfile(filename string) io.Writer {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file %s: %v", filename, err)
		os.Exit(-1)
	}

	return f
}

func handleBotReply(c chan cmatrix.AttachmentChannel) {
	for {
		ac := <-c
		params := slack.PostMessageParameters{}
		params.AsUser = true
		params.Attachments = []slack.Attachment{*ac.Attachment}
		_, _, errPostMessage := api.PostMessage(ac.Channel.Name, ac.DisplayTitle, params)
		if errPostMessage != nil {
			log.Fatal(errPostMessage)
		}
	}
}
