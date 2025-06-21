package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/Knetic/govaluate"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	exprRegex := regexp.MustCompile(`^[0-9+\-*/().\s]+$`)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		input := update.Message.Text

		if !exprRegex.MatchString(input) {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка: допустимы только числа и операторы (+ - * /)")
			bot.Send(msg)
			continue
		}

		expression, err := govaluate.NewEvaluableExpression(input)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка парсинга выражения"))
			continue
		}

		result, err := expression.Evaluate(nil)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка вычисления"))
			continue
		}

		reply := fmt.Sprintf("Результат: %v", result)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)
	}
}
