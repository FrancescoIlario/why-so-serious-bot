package bot

import (
	"time"

	"github.com/FrancescoIlario/why-so-serious-bot/pkg/wssface"
	"github.com/FrancescoIlario/why-so-serious-bot/pkg/wssformrecognizer"
	"github.com/FrancescoIlario/why-so-serious-bot/pkg/wssmoderator"
	"github.com/FrancescoIlario/why-so-serious-bot/pkg/wsssentiment"
	"github.com/FrancescoIlario/why-so-serious-bot/pkg/wsstranslator"
	"github.com/FrancescoIlario/why-so-serious-bot/pkg/wssvision"
)

//Configuration Bot's Configuration
type Configuration struct {
	PollerInterval     time.Duration
	Token              string
	FaceConf           *wssface.Configuration
	VisionConf         *wssvision.Configuration
	TextAnalyticsConf  *wsssentiment.Configuration
	ModeratorConf      *wssmoderator.Configuration
	TranslatorConf     *wsstranslator.Configuration
	FormRecognizerConf *wssformrecognizer.Configuration
}
