package main

import (
	"log"

	"github.com/FrancescoIlario/why-so-serious-bot/internal/bot"
	"github.com/FrancescoIlario/why-so-serious-bot/internal/conf"
	"github.com/FrancescoIlario/why-so-serious-bot/pkg/wssface"
	"github.com/FrancescoIlario/why-so-serious-bot/pkg/wssvision"
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

	visionConf, err := wssvision.BuildConfigurationFromEnvs()
	if err != nil {
		log.Fatalf("error retrieving vision service configuration: %v", err)
	}

	// instantiate bot
	settings := tb.Settings{
		Token: *token,
		Poller: &tb.LongPoller{
			Timeout: *pollerInterval,
		},
	}

	fbot, err := bot.New(settings, *faceConf, *visionConf)
	if err != nil {
		log.Printf("can not instantiate bot: %v", err)
	}

	// start bot
	fbot.Start()

	// wait undefinetly
	shutdown := make(chan struct{})
	<-shutdown
}
