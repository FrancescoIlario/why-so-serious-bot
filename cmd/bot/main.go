package main

import (
	"fmt"
	"log"
	"time"

	"github.com/FrancescoIlario/why-so-serious-bot/internal/bot"
	"github.com/FrancescoIlario/why-so-serious-bot/internal/conf"
	"github.com/FrancescoIlario/why-so-serious-bot/pkg/wssface"
	"github.com/FrancescoIlario/why-so-serious-bot/pkg/wssmoderator"
	"github.com/FrancescoIlario/why-so-serious-bot/pkg/wsssentiment"
	"github.com/FrancescoIlario/why-so-serious-bot/pkg/wsstranslator"
	"github.com/FrancescoIlario/why-so-serious-bot/pkg/wssvision"
	tb "gopkg.in/tucnak/telebot.v2"
)

type configurations struct {
	pollerInterval    *time.Duration
	token             *string
	faceConf          *wssface.Configuration
	visionConf        *wssvision.Configuration
	textAnalyticsConf *wsssentiment.Configuration
	moderatorConf     *wssmoderator.Configuration
	translatorConf    *wsstranslator.Configuration
}

func main() {
	conf, err := getConfigurations()
	if err != nil {
		log.Fatalln(err)
	}

	// instantiate bot
	settings := tb.Settings{
		Token: *conf.token,
		Poller: &tb.LongPoller{
			Timeout: *conf.pollerInterval,
		},
	}

	fbot, err := bot.New(settings, *conf.faceConf, *conf.visionConf,
		*conf.textAnalyticsConf, *conf.moderatorConf, *conf.translatorConf)
	if err != nil {
		log.Printf("can not instantiate bot: %v", err)
	}

	// start bot
	fbot.Start()

	// wait undefinetly
	shutdown := make(chan struct{})
	<-shutdown
}

func getConfigurations() (*configurations, error) {
	// get configurations
	pollerInterval, err := conf.GetPollerInterval()
	if err != nil {
		return nil, fmt.Errorf("error retrieving poller interval: %v", err)
	}

	token, err := conf.GetToken()
	if err != nil {
		return nil, fmt.Errorf("error retrieving Telegram token: %v", err)
	}

	faceConf, err := wssface.BuildConfigurationFromEnvs()
	if err != nil {
		return nil, fmt.Errorf("error retrieving face service configuration: %v", err)
	}

	visionConf, err := wssvision.BuildConfigurationFromEnvs()
	if err != nil {
		return nil, fmt.Errorf("error retrieving vision service configuration: %v", err)
	}

	textAnalyticsConf, err := wsssentiment.BuildConfigurationFromEnvs()
	if err != nil {
		return nil, fmt.Errorf("error retrieving text analitycs service configuration: %v", err)
	}

	moderatorConf, err := wssmoderator.BuildConfigurationFromEnvs()
	if err != nil {
		return nil, fmt.Errorf("error retrieving moderator service configuration: %v", err)
	}

	translatorConf, err := wsstranslator.BuildConfigurationFromEnvs()
	if err != nil {
		return nil, fmt.Errorf("error retrieving translator service configuration: %v", err)
	}

	return &configurations{
		pollerInterval:    pollerInterval,
		token:             token,
		faceConf:          faceConf,
		visionConf:        visionConf,
		textAnalyticsConf: textAnalyticsConf,
		moderatorConf:     moderatorConf,
		translatorConf:    translatorConf,
	}, nil
}
