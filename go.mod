module github.com/FrancescoIlario/wss-bot

go 1.14

replace github.com/FrancescoIlario/wss-bot => ./

require (
	github.com/Azure/azure-sdk-for-go v42.0.0+incompatible
	github.com/Azure/go-autorest/autorest v0.10.0
	github.com/Azure/go-autorest/autorest/validation v0.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	gopkg.in/tucnak/telebot.v2 v2.0.0-20200426184946-59629fe0483e
)