package bot

import (
	"context"
	"fmt"
	"log"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (b *Bot) onPhoto(m *tb.Message) {
	rc, err := b.tbot.GetFile(&m.Photo.File)
	if err != nil {
		log.Printf("error reading provided photo: %v", err)
		return
	}

	faceContext := context.Background()
	res, err := b.faceCli.InvokeFace(faceContext, rc)
	if err != nil {
		log.Printf("error invoking face API: %v", err)
		b.tbot.Send(m.Chat, err.Error())
		return
	}

	var message string
	if res.Gender == "male" {
		message = fmt.Sprintf("Hi man, I guess you are %v years old. Why you so %s?", res.Age, res.Sentiment.Adjective())
	} else {
		message = fmt.Sprintf("Hi darling, I guess you are %v years old. Why you so %s?", res.Age, res.Sentiment.Adjective())
	}
	b.tbot.Send(m.Chat, message)
}
