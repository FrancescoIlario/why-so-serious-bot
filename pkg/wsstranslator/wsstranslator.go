package wsstranslator

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v3.0/translatortext"
	"github.com/Azure/go-autorest/autorest"
)

//TranslatorServiceClient client for the Azure Translator Service
type TranslatorServiceClient struct {
	conf          Configuration
	translatorCli *translatortext.TranslatorClient
}

//NewTranslatorServiceClient TranslatorServiceClient constructor
func NewTranslatorServiceClient(conf Configuration) *TranslatorServiceClient {
	client := translatortext.NewTranslatorClient(conf.ServiceEnpoint)
	client.Authorizer = autorest.NewCognitiveServicesAuthorizer(conf.TranslatorSubscription)

	return &TranslatorServiceClient{
		conf:          conf,
		translatorCli: &client,
	}
}

//InvokeTranslator invokes the Translator APIs with the provided photo
func (s *TranslatorServiceClient) InvokeTranslator(translatorContext context.Context, message string, lang *string) (*TranslatorResult, error) {
	toLang := []string{"en"}
	if lang != nil {
		toLang[0] = *lang
	}

	input := []translatortext.TranslateTextInput{{Text: &message}}
	res, err := s.translatorCli.Translate(translatorContext, toLang, input, "", "", "", "", "", nil, nil, "en", "", nil, "")
	if err != nil {
		return nil, err
	}

	v := res.Value
	if v == nil || len(*v) == 0 {
		return &TranslatorResult{}, nil
	}
	fv := (*v)[0]

	var translation *translatortext.TranslateResultAllItemTranslationsItem
	if translations := fv.Translations; len(*translations) > 0 {
		translation = &(*translations)[0]
	}

	return &TranslatorResult{
		Translation:    translation.Text,
		IdentifiedLang: fv.DetectedLanguage.Language,
	}, nil
}

// TranslatorResult result of the TranslatorAPI
type TranslatorResult struct {
	Translation    *string
	IdentifiedLang *string
}
