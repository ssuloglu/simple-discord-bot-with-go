package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Bot prefix for commands
const prefix string = "!lib"

// urls for specific languages
var lib = map[string][]string{
	"golang": {"https://gobyexample.com/", "https://roadmap.sh/golang", "https://go.dev/doc/effective_go"},
	"python": {"https://www.python.org/about/gettingstarted/", "https://www.learnbyexample.org/python/"},
	"c":      {"https://www.freecodecamp.org/news/the-c-beginners-handbook/", "https://www.w3schools.com/c/"},
}

func getUrls(lang string) string {
	return strings.Join(lib[lang], "\n")
}

func addUrl(lang string, newURL string) {
	lib[lang] = append(lib[lang], newURL)
}

func main() {

	// <token> is replaced with the token generated for the specific bot created and to be
	// integrated which is obtained from https://discord.com/developers/applications
	sess, err := discordgo.New("Bot <token>")
	if err != nil {
		log.Fatal(err)
	}
	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Check if libbot sends a message itself
		if m.Author.ID == s.State.User.ID {
			return
		}

		// Parse message to identify the intent
		args := strings.Split(m.Content, " ")

		// Check if
		if args[0] != prefix || len(args) == 1 {
			return
		}

		if args[1] == "get" {
			s.ChannelMessageSend(m.ChannelID, getUrls(args[2]))
		} else if args[1] == "add" {
			addUrl(args[2], args[3])
			s.ChannelMessageSend(m.ChannelID, "Excellent! New url is added!")
		}
	})

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}

	defer sess.Close()

	fmt.Println("The bot is online!")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
