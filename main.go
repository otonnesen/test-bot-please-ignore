package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
)

const threadTemplate = `Daily game thread for %s.
Wordle: https://www.nytimes.com/games/wordle/index.html
Travle: https://travle.earth/
Globle: https://globle-game.com/
Timeguessr: https://timeguessr.com/
Costcodle: https://costcodle.com/`

var (
	GuildId       = flag.String("guild", "", "Guild ID")
	TextChannelId = flag.String("channel", "", "Text Channel ID")
	BotToken      = flag.String("token", "", "Bot Token")
)

var stop = make(chan os.Signal, 1)

func main() {
	flag.Parse()
	if *GuildId == "" || *TextChannelId == "" || *BotToken == "" {
		flag.Usage()
		os.Exit(1)
	}

	discord, err := discordgo.New("Bot " + *BotToken)
	if err != nil {
		panic(err)
	}

	discord.AddHandler(createDailyGameThread)

	err = discord.Open()
	if err != nil {
		panic(err)
	}
	defer discord.Close()

	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Printf("Shutting down...")
}

func createDailyGameThread(s *discordgo.Session, r *discordgo.Ready) {
	today := time.Now()
	threadName := fmt.Sprintf("Daily game thread for %s", today.Format("Mon, 02 Jan 2006"))
	threadContent := fmt.Sprintf(threadTemplate, today.Format("Mon, 02 Jan 2006"))

	ch, err := s.ThreadStart(*TextChannelId, threadName, discordgo.ChannelTypeGuildPublicThread, 1440)
	if err != nil {
		panic(err)
	}

	_, err = s.ChannelMessageSend(ch.ID, threadContent)
	if err != nil {
		panic(err)
	}
	log.Printf("Created thread %s, stopping...", ch.ID)
	stop <- os.Interrupt
}
