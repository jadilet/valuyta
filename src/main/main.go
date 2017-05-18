package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"gopkg.in/telegram-bot-api.v4"
)

type Joke struct {
	ID   uint32 `json:"id"`
	Joke string `json:"joke"`
}

type JokeResponse struct {
	Value Joke   `json:"value"`
	Type  string `json:"type"`
}

var buttons = []tgbotapi.KeyboardButton{
	tgbotapi.KeyboardButton{Text: "Валюта"},
}

// При старте приложения, оно скажет телеграму ходить с обновлениями по этому URL
const WebhookURL = "https://msu-go-2017.herokuapp.com/"

func getJoke() string {
	client := http.Client{}
	response, error := client.Get("http://api.icndb.com/jokes/random?limitTo=[nerdy]")

	if error != nil {
		return "Joke api not responding!"
	}

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	joke := JokeResponse{}

	err := json.Unmarshal(body, &joke)

	if err != nil {
		return "Joke error while parsing!"
	}

	return joke.Value.Joke
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
			message = tgbotapi.NewMessage(update.Message.Chat.ID, "https://valuta.kg/")
		default:
			message = tgbotapi.NewMessage(update.Message.Chat.ID, `Нажмите "Валюта"`)
		}

		// В ответном сообщении просим показать клавиатуру
		message.ReplyMarkup = tgbotapi.NewReplyKeyboard(buttons)

		bot.Send(message)
	}

	fmt.Println(getJoke())
}
