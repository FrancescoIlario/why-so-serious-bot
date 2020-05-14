package bot

import (
	"github.com/FrancescoIlario/why-so-serious-bot/pkg/tglog"
	"github.com/FrancescoIlario/why-so-serious-bot/pkg/wssface"
	"github.com/FrancescoIlario/why-so-serious-bot/pkg/wssformrecognizer"
	"github.com/FrancescoIlario/why-so-serious-bot/pkg/wssmoderator"
	"github.com/FrancescoIlario/why-so-serious-bot/pkg/wsssentiment"
	"github.com/FrancescoIlario/why-so-serious-bot/pkg/wsstranslator"
	"github.com/FrancescoIlario/why-so-serious-bot/pkg/wssvision"
	tb "gopkg.in/tucnak/telebot.v2"
)

//Bot WhySoSerious Bot implementation
type Bot struct {
	tbot              tb.Bot
	faceCli           *wssface.FaceServiceClient
	visionCli         *wssvision.VisionServiceClient
	textAnalyticsCli  *wsssentiment.TextAnalyticsServiceClient
	moderatorCli      *wssmoderator.ContentModeratorServiceClient
	translatorCli     *wsstranslator.TranslatorServiceClient
	formRecognizerCli *wssformrecognizer.FormRecognizerServiceClient
}

//New Bot constructor
func New(c Configuration) (*Bot, error) {
	bot, err := new(c)
	if err != nil {
		return nil, err
	}

	{ // Handler: Photo
		bot.tbot.Handle(tb.OnPhoto, bot.onPhoto)
	}
	{ // Handler: Text
		bot.tbot.Handle(tb.OnText, tglog.Wrap(bot.onText))
	}
	{ // Command: Translate
		bot.tbot.Handle("/t", tglog.Wrap(bot.translate))
		bot.tbot.Handle("/translate", tglog.Wrap(bot.translate))
	}
	{ // Command: Languages
		bot.tbot.Handle("/l", tglog.Wrap(bot.languages))
		bot.tbot.Handle("/languages", tglog.Wrap(bot.languages))
	}

	return bot, nil
}

func new(c Configuration) (*Bot, error) {
	// telebot settings
	tbSettings := tb.Settings{
		Token: c.Token,
		Poller: &tb.LongPoller{
			Timeout: c.PollerInterval,
		},
	}

	// instantiate telebot bot
	tbot, err := tb.NewBot(tbSettings)
	if err != nil {
		return nil, err
	}

	// instantiate our bot
	bot := &Bot{
		tbot:              *tbot,
		faceCli:           wssface.NewFaceServiceClient(c.FaceConf),
		visionCli:         wssvision.NewVisionServiceClient(c.VisionConf),
		textAnalyticsCli:  wsssentiment.NewTextAnalyticsServiceClient(c.TextAnalyticsConf),
		moderatorCli:      wssmoderator.NewContentModeratorServiceClient(c.ModeratorConf),
		translatorCli:     wsstranslator.NewTranslatorServiceClient(c.TranslatorConf),
		formRecognizerCli: wssformrecognizer.NewFormRecognizerServiceClient(c.FormRecognizerConf),
	}
	return bot, nil
}

// Start starts the telegram bot
func (b *Bot) Start() {
	b.tbot.Start()
}
