package wsstranslator

import (
	"fmt"
	"os"
)

//Configuration structure for the Translator service configuration
type Configuration struct {
	TranslatorSubscription string
	ServiceEnpoint         string
}

const (
	//TranslatorSubscriptionKey Azure Translator Subscription env key
	TranslatorSubscriptionKey = "WSS_TRANSLATOR_SUBSCRIPTION_KEY"
	//TranslatorEndpointKey Azure Translator Endpoint env key
	TranslatorEndpointKey = "WSS_TRANSLATOR_ENDPOINT"
)

//BuildConfigurationFromEnvs builds the configuration from env variables
func BuildConfigurationFromEnvs() (*Configuration, error) {
	sub, err := getTranslatorSubscription()
	if err != nil {
		return nil, err
	}

	endStr, err := getTranslatorEndpoint()
	if err != nil {
		return nil, err
	}

	return &Configuration{
		TranslatorSubscription: *sub,
		ServiceEnpoint:         *endStr,
	}, nil
}

func getTranslatorSubscription() (*string, error) {
	if v := os.Getenv(TranslatorSubscriptionKey); v != "" {
		return &v, nil
	}

	return nil, fmt.Errorf("translator subscription env key (%s) not set", TranslatorSubscriptionKey)
}

func getTranslatorEndpoint() (*string, error) {
	if v := os.Getenv(TranslatorEndpointKey); v != "" {
		return &v, nil
	}

	return nil, fmt.Errorf("translator subscription env key (%s) not set", TranslatorEndpointKey)
}
