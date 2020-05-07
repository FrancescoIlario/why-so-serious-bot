package bot

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/FrancescoIlario/why-so-serious-bot/pkg/wssface"
	"github.com/FrancescoIlario/why-so-serious-bot/pkg/wssmoderator"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (b *Bot) onPhoto(m *tb.Message) {
	rc, err := b.tbot.GetFile(&m.Photo.File)
	if err != nil {
		log.Printf("error reading provided photo: %v", err)
		return
	}
	defer rc.Close()

	image, err := ioutil.ReadAll(rc)
	if err != nil {
		log.Printf("error reading provided photo: %v", err)
		return
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mChan := b.invokeModeratorAPI(ctx, m, image)
	frChan := b.invokeFaceAPI(ctx, m, image)

	if moderatorResult := <-mChan; moderatorResult != nil {
		if moderatorResult.Adult || moderatorResult.Racy {
			b.tbot.Send(m.Chat, "This content is inappropriate: it seems to have adult or racy content")
			return
		}
	}

	faceResult := <-frChan
	faces := faceResult.Faces
	switch len(faces) {
	case 0:
		b.processNoFacePhoto(m.Chat, image)
		break
	case 1:
		b.processSingleFacePhoto(m.Chat, faces[0])
		break
	default:
		b.processGroupPhoto(m.Chat, faces)
		break
	}
}

func (b *Bot) invokeModeratorAPI(ctx context.Context, m *tb.Message, image []byte) chan *wssmoderator.ContentModeratorPhotoResult {
	mChan := make(chan *wssmoderator.ContentModeratorPhotoResult)
	go func() {
		defer close(mChan)
		rct := ioutil.NopCloser(bytes.NewReader(image))
		defer rct.Close()
		moderatorResult, err := b.moderatorCli.InvokeContentModeratorPhoto(ctx, rct)
		if err != nil {
			log.Printf("error invoking face API: %v", err)
		}
		mChan <- moderatorResult
	}()

	return mChan
}

func (b *Bot) invokeFaceAPI(ctx context.Context, m *tb.Message, image []byte) chan *wssface.FaceResult {
	frChan := make(chan *wssface.FaceResult)
	go func() {
		defer close(frChan)

		rct := ioutil.NopCloser(bytes.NewReader(image))
		defer rct.Close()
		faceResult, err := b.faceCli.InvokeFace(ctx, rct)

		if err != nil {
			log.Printf("error invoking face API: %v", err)
		}
		frChan <- faceResult
	}()

	return frChan
}

func (b *Bot) processGroupPhoto(chat *tb.Chat, faces []wssface.FaceDetails) {
	message := fmt.Sprintf("What a nice group picture of %v of you", len(faces))
	b.tbot.Send(chat, message)
}

func (b *Bot) processSingleFacePhoto(chat *tb.Chat, face wssface.FaceDetails) {
	message := fmt.Sprintf("Hello %s, I guess you are %v years old. Why you so %s?",
		b.genderGreet(face.Gender), face.Age, face.Sentiment.Adjective())

	b.tbot.Send(chat, message)
}

func (b *Bot) processNoFacePhoto(chat *tb.Chat, image []byte) {
	visionContext := context.Background()
	rct := ioutil.NopCloser(bytes.NewReader(image))
	defer rct.Close()
	res, err := b.visionCli.InvokeVision(visionContext, rct)
	if err != nil {
		log.Printf(`error invoking computer vision service: %v`, err)
		b.tbot.Send(chat, `This picture makes me feel a sick! I'm sorry, I can't handle this request!`)
		return
	}

	if res.Description == nil {
		b.tbot.Send(chat, `I'm sorry but I can't figure out what this picture represents`)
		return
	}

	message := fmt.Sprintf("It seems %s", *res.Description)
	b.tbot.Send(chat, message)
}

func (b *Bot) genderGreet(gender string) string {
	if gender == "male" {
		return "man"
	}
	return "darling"
}
