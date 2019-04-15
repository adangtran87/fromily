package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/adangtran87/fromily/fromilyclient"
	"github.com/bwmarrin/discordgo"
)

type Configs struct {
	Token        string   `json:"TOKEN"`
	Prefix       string   `json:"PREFIX"`
	AdminPrefix  string   `json:"ADMIN_PREFIX"`
	Admins       []string `json:"ADMINS"`
	FromilyToken string   `json:"FROMILY_TOKEN"`
}

// Config struct populated in main
var config = Configs{}

// Normal command regex set in main
var prefix_regex *regexp.Regexp

// Admin command regex configured in main
var admin_regex *regexp.Regexp

var Backend = ServerBackend{}

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

	// Create fromily-server session
	Backend.Init()
	Backend.Client = fromilyclient.New(config.FromilyToken)
	if err != nil {
		fmt.Println("Error creating fromily session,", err)
	}

	if Backend.RefreshInfo() == false {
		fmt.Println("Error refreshing data from server")
		return
	}

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

	discord.AddHandler(ready)
	discord.AddHandler(messageCreate)
	discord.AddHandler(guildMemberAdd)
	discord.AddHandler(guildCreate)

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
		DUTIL_UpdateGuildInfo(s, guild)
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

	go Commands.Dispatch(s, m, config.Prefix, m.Content)
}

// GuildMemberAdd event
func guildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	DUTIL_UpdateMember(m.Member)
}

func guildCreate(s *discordgo.Session, g *discordgo.GuildCreate) {
	DUTIL_UpdateGuildInfo(s, g.Guild)
}
