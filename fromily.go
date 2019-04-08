package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Config struct populated in main
var config = Configs{}

// Normal command regex set in main
var prefix_regex *regexp.Regexp

// Admin command regex configured in main
var admin_regex *regexp.Regexp

// Opens a discord session and monitors messages sent
// Processes commands if messages have the appropriate prefix
func main() {
	// Read config file
	config_json, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(config_json, &config)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("%s", config.Token)

	// Create new Discord session
	discord, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("Error creating discord session,", err)
		return
	} else {
		fmt.Println("Discord session is alive")
	}

	// Check if prefixes exist
	// TODO: Eventually get this from database per server
	if config.Prefix == "" {
		panic("FROMILY_ERROR: Invalid prefix!")
	}
	// Compile regexes for use
	prefix_regex = regexp.MustCompile(`^` + config.Prefix + `(\w+)`)

	if config.AdminPrefix == "" {
		panic("FROMILY_ERROR: Invalid admin prefix!")
	}
	admin_regex = regexp.MustCompile(`^` + config.AdminPrefix + `(\w+)`)

	discord.AddHandler(ready)
	discord.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Fromily is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	// event.Guilds retreives a list of connected guild ids
	for _, guild := range event.Guilds {
		guildInfo, _ := s.Guild(guild.ID)
		fmt.Printf("%s:%s\n", guildInfo.Name, guildInfo.ID)
	}
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	go CommandDispatch(s, m)

	go AdminDispatch(s, m)
}