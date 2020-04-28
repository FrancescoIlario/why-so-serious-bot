package bot

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

func (b *Bot) onText(m *tb.Message) {
	message := "Hi dude"
	b.tbot.Send(m.Chat, message)
}
