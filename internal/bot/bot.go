package bot

import (
	"github.com/FrancescoIlario/why-so-serious-bot/pkg/tglog"
	"github.com/FrancescoIlario/why-so-serious-bot/pkg/wssface"
	"github.com/FrancescoIlario/why-so-serious-bot/pkg/wssmoderator"
	"github.com/FrancescoIlario/why-so-serious-bot/pkg/wsssentiment"
	"github.com/FrancescoIlario/why-so-serious-bot/pkg/wssvision"
	tb "gopkg.in/tucnak/telebot.v2"
)

//Bot WhySoSerious Bot implementation
type Bot struct {
	tbot             tb.Bot
	faceCli          wssface.FaceServiceClient
	visionCli        wssvision.VisionServiceClient
	textAnalyticsCli wsssentiment.TextAnalyticsServiceClient
	moderatorCli     wssmoderator.ContentModeratorServiceClient
}

//New Bot constructor
func New(tbotSettings tb.Settings,
	faceConf wssface.Configuration,
	visionConf wssvision.Configuration,
	textAnalyticsConf wsssentiment.Configuration,
	moderatorConf wssmoderator.Configuration) (*Bot, error) {
	tbot, err := tb.NewBot(tbotSettings)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		tbot:             *tbot,
		faceCli:          *wssface.NewFaceServiceClient(faceConf),
		visionCli:        *wssvision.NewVisionServiceClient(visionConf),
		textAnalyticsCli: *wsssentiment.NewTextAnalyticsServiceClient(textAnalyticsConf),
		moderatorCli:     *wssmoderator.NewContentModeratorServiceClient(moderatorConf),
	}

	tbot.Handle(tb.OnPhoto, bot.onPhoto)
	tbot.Handle(tb.OnText, tglog.Wrap(bot.onText))

	return bot, nil
}

// Start starts the telegram bot
func (b *Bot) Start() {
	b.tbot.Start()
}
