package main

import (
	"flag"
	"log"
	"os/exec"
	"regexp"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var token = flag.String("token", "", "token for telegram")
var provider = flag.String("provider", "", "provider token for telegram")

func init() { flag.Parse() }

func linkUserName(username string) string {
	return link("https://t.me/"+username, "@"+username)
}

func main() {
	command := exec.Command("echo", "world")
	command.Output()

	var bot, _ = api.NewBotAPI(*token)
	log.Printf("authorized as %s", linkUserName(bot.Self.UserName))

	u := api.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if message := update.Message; message != nil {
			switch message.Command() {
			// платная регистрация пользователя на сервере
			case "subcribe":
				username := Username(*message)
				bot.Send(api.InvoiceConfig{
					BaseChat: api.BaseChat{
						ChatConfig: api.ChatConfig{ChatID: message.Chat.ID},
					},
					Title:          "трехмесячный пропуск",
					Description:    "доступ к коллективному серверу по адресу tit.zba.su для юзернейма " + username + ". вы платите на три месяца вперёд по 42,66 ₽.",
					Payload:        username,
					ProviderToken:  *provider,
					StartParameter: "",
					Currency:       "RUB",
					Prices: []api.LabeledPrice{
						{Label: "2 гб для домашней директории", Amount: 7882},
						{Label: "заработная плата администрации", Amount: 1085},
						{Label: "взимания юкассы и фнс", Amount: 1164},
					},
					MaxTipAmount:        45000000,
					SuggestedTipAmounts: []int{2667, 5334, 10668},
				})
			default:
				continue
			}
		}
		if precheckout := update.PreCheckoutQuery; precheckout != nil {
			username := precheckout.InvoicePayload
			password, error := register(username)

			if error == nil {
				bot.Send(api.PreCheckoutConfig{
					PreCheckoutQueryID: precheckout.ID,
					OK:                 true,
					ErrorMessage:       "привет",
				})
				bot.Send(api.MessageConfig{
					BaseChat: api.BaseChat{
						ChatConfig: api.ChatConfig{ChatID: precheckout.From.ID},
					},
					Text: "*создан пользователь* `" + username + "`" +
						", пароль: `" + password + "`" +
						"\n```sh\nssh " + username + "@tit.zba.su\n```",
					ParseMode: "MarkdownV2",
				})
			} else {
				log.Printf("didn't registered %s, \"%s\"", linkUserName(precheckout.From.UserName), error)
				bot.Send(api.PreCheckoutConfig{
					PreCheckoutQueryID: precheckout.ID,
					OK:                 false,
					ErrorMessage:       "имя пользователя занято",
				})
			}
		}
	}
}

// определение юзернейма
func Username(message api.Message) (username string) {
	argument := message.CommandArguments()
	isValidLogin, _ := regexp.MatchString("[a-z_][a-z0-9_-]*[$]?", argument)
	if argument != "" && isValidLogin {
		return argument
	}

	return message.From.UserName
}
