// SPDX-FileCopyrightText: 2024 Bogdan Alekseevich Zazhigin <zaboal@tuta.io>
// SPDX-License-Identifier: 0BSD

package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"regexp"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	ansi "github.com/zaboal/tilde/internal/ansi"
	tilde "github.com/zaboal/tilde/internal/tilde"
)

var logger = log.New(os.Stdout, ansi.Bold("telegram bot "), log.Ldate+log.Ltime+log.Lmsgprefix)

var token = flag.String("token", "", "token for telegram")
var provider = flag.String("provider", "", "provider token for telegram")

func init() { flag.Parse() }

func main() {
	var bot, _ = api.NewBotAPI(*token)
	logger.Print("authorized as " + tme(bot.Self.UserName))

	u := api.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if message := update.Message; message != nil {
			switch message.Command() {
			// платная регистрация пользователя на сервере
			case "subcribe":
				username := chooseUsername(*message)
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
			logger.Printf("got a precheckout with payload \"%s\"", username)
			password, error := tilde.Subscribe(username)

			if error == nil {
				bot.Send(api.PreCheckoutConfig{
					PreCheckoutQueryID: precheckout.ID,
					OK:                 true,
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
				logger.Printf("did "+ansi.Italic("not")+" subscribed %s: %s", tme(precheckout.From.UserName), error)

				var errorMessage string
				var userNameExistsError *tilde.UserNameExistsError
				if errors.As(error, &userNameExistsError) {
					errorMessage = "имя пользователя занято, напишите другое после /subcribe"
				}

				bot.Send(api.PreCheckoutConfig{
					PreCheckoutQueryID: precheckout.ID,
					OK:                 false,
					ErrorMessage:       errorMessage,
				})
			}
		}
	}
}

// hyperlink telegram usernames for console
func tme(username string) string {
	return ansi.Link("@"+username, "https://t.me/"+username)
}

// choose a username from the message
func chooseUsername(message api.Message) (username string) {
	argument := message.CommandArguments()
	isValidLogin, _ := regexp.MatchString("[a-z_][a-z0-9_-]*[$]?", argument)
	if argument != "" && isValidLogin {
		return argument
	}

	return message.From.UserName
}
