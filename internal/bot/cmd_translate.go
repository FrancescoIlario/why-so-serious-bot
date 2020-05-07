package bot

import (
	"context"
	"log"
	"regexp"

	"github.com/FrancescoIlario/why-so-serious-bot/pkg/wsstranslator"
	tb "gopkg.in/tucnak/telebot.v2"
)

const langPattern = "-to:([A-z]{2}) "

func (b *Bot) translate(m *tb.Message) {
	if !m.Private() {
		return
	}

	language := extractLang(m.Payload)
	message := m.Payload
	if language != nil {
		message = message[7:]
	}

	ctx := context.Background()
	res, err := b.translatorCli.InvokeTranslator(ctx, message, language)
	if err != nil {
		log.Printf("error invoking translator service: %v", err)
		return
	}

	response := generateTranslationResponse(res)
	b.tbot.Send(m.Chat, response)
}

func generateTranslationResponse(res *wsstranslator.TranslatorResult) string {
	var message string
	if res.IdentifiedLang != nil {
		message += "Identified language " + *res.IdentifiedLang + ". "
	}

	if res.Translation != nil {
		message += "Translation:\n\n" + *res.Translation
	} else {
		message += "Unfortunately, I was not able to translate what you said"
	}

	return message
}

func extractLang(payload string) *string {
	exp, err := regexp.Compile(langPattern)
	if err != nil {
		log.Printf("can not compile regexp: %s", langPattern)
		return nil
	}

	if submatches := exp.FindStringSubmatch(payload); len(submatches) >= 1 {
		return &submatches[0]
	}
	return nil
}