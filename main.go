package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages + discordgo.IntentsDirectMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	fmt.Println(m.Author.String() + ": " + m.Content)
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}

	if m.Content == "time" {
		s.ChannelMessageSend(m.ChannelID, "<@!"+m.Author.ID+"> "+fmt.Sprint(time.Now()))
	}

	if strings.Contains(strings.ToLower(m.Content), "uwu") {
		s.ChannelMessageSend(m.ChannelID, "<@!"+m.Author.ID+"> "+"????")
	}

	if strings.HasPrefix(m.Content, "!rand") {
		spaced_words := strings.Split(m.Content, " ")
		if len(spaced_words) == 2 {
			max, err := strconv.ParseInt(strings.Split(m.Content, " ")[1], 10, 64)
			if err != nil {
				fmt.Println(err)
				s.ChannelMessageSend(m.ChannelID, "!rand argument must be integer")
			} else {
				rand.Seed(time.Now().UnixNano())
				s.ChannelMessageSend(m.ChannelID, fmt.Sprint(rand.Int()%int(max)))
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, "!rand requires one, and only one argument")
		}
	}
}
