package main

import (
	"log"

	"github.com/FrancescoIlario/wss-bot/internal/bot"
	"github.com/FrancescoIlario/wss-bot/internal/conf"
	"github.com/FrancescoIlario/wss-bot/pkg/wssface"
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	// get configurations
	pollerInterval, err := conf.GetPollerInterval()
	if err != nil {
		log.Fatalf("error retrieving poller interval: %v", err)
	}

	token, err := conf.GetToken()
	if err != nil {
		log.Fatalf("error retrieving Telegram token: %v", err)
	}

	faceConf, err := wssface.BuildConfigurationFromEnvs()
	if err != nil {
		log.Fatalf("error retrieving face service configuration: %v", err)
	}

	// instantiate bot
	settings := tb.Settings{
		Token: *token,
		Poller: &tb.LongPoller{
			Timeout: *pollerInterval,
		},
	}

	fbot, err := bot.New(settings, *faceConf)
	if err != nil {
		log.Printf("can not instantiate bot: %v", err)
	}

	// start bot
	fbot.Start()

	// wait undefinetly
	shutdown := make(chan struct{})
	<-shutdown
}
