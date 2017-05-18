package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/html/charset"
	"gopkg.in/telegram-bot-api.v4"
)

type CurrencyRates struct {
	Date     string `xml:"Date,attr"`
	Currency []struct {
		ISOCode string `xml:"ISOCode,attr"`
		Nominal int    `xml:"Nominal"`
		Value   string `xml:"Value"`
	}
}

var buttons = []tgbotapi.KeyboardButton{
	tgbotapi.KeyboardButton{Text: "Валюта"},
}

// При старте приложения, оно скажет телеграму ходить с обновлениями по этому URL
const WebhookURL = "https://valuyta.herokuapp.com/"

func getCurrency() CurrencyRates {
	client := http.Client{}
	response, error := client.Get("http://www.nbkr.kg/XML/daily.xml")

	if error != nil {
		log.Println("NBKR not responding")
		return CurrencyRates{}
	}

	defer response.Body.Close()

	var currencyRates CurrencyRates

	decoder := xml.NewDecoder((response.Body))
	decoder.CharsetReader = charset.NewReaderLabel
	decoder.Decode(&currencyRates)

	return currencyRates
}

func main() {
	// Heroku прокидывает порт для приложения в переменную окружения PORT
	port := os.Getenv("PORT")
	bot, err := tgbotapi.NewBotAPI("351660528:AAGp-U1bseZwIzWvxTwroQNHm1JCPA7TSXY")

	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	log.Printf("First Name %s", bot.Self.FirstName)
	log.Printf("Last Name %s", bot.Self.LastName)
	log.Printf("Telephone %s", bot.Self.UserName)

	// Устанавливаем вебхук
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL))

	if err != nil {
		log.Fatal(err)
	}

	updates := bot.ListenForWebhook("/")
	go http.ListenAndServe(":"+port, nil)

	// получаем все обновления из канала updates
	for update := range updates {
		var message tgbotapi.MessageConfig
		log.Println("received text: ", update.Message.Text)

		switch update.Message.Text {
		case "Валюта":
			// Если пользователь нажал на кнопку, то придёт сообщение "Валюта"
			var currencyRates CurrencyRates
			currencyRates = getCurrency()
			currencyInfo := fmt.Sprintf("КУРС НБКР\nпродажа\n")

			for i, _ := range currencyRates.Currency {
				data := fmt.Sprintf("%s %s \n", currencyRates.Currency[i].ISOCode, currencyRates.Currency[i].Value)
				currencyInfo = currencyInfo + data
			}

			message = tgbotapi.NewMessage(update.Message.Chat.ID, currencyInfo)
		default:
			message = tgbotapi.NewMessage(update.Message.Chat.ID, `Нажмите "Валюта"`)
		}

		// В ответном сообщении просим показать клавиатуру
		message.ReplyMarkup = tgbotapi.NewReplyKeyboard(buttons)

		bot.Send(message)
	}
}
