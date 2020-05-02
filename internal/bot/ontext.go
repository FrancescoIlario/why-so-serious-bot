package bot

import (
	"context"
	"log"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (b *Bot) onText(m *tb.Message) {
	log.Printf("received message: %v", m.Text)
	textAnalyticsContext := context.Background()
	res, err := b.textAnalyticsCli.InvokeTextAnalytics(textAnalyticsContext, m.Text)
	if err != nil {
		log.Printf("error invoking analytics service: %v", err)
	}

	if score := res.SentimentScore; score == nil {
		b.tbot.Send(m.Chat, "Sorry, I didn't understand what your sentiment is")
	} else if *score <= 0.3 {
		b.tbot.Send(m.Chat, "I feel you a little bit upset!")
	} else if *score >= 0.7 {
		b.tbot.Send(m.Chat, "What a nice positive message! Thanks!")
	} else {
		b.tbot.Send(m.Chat, "Why so serious?")
	}
}
