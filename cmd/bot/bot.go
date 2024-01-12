package main

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"os"
)

const filePath = "chatIDS.json"

func main() {
	godotenv.Load(".env")

	chatIDS, err := readChatIDS()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Create a new UpdateConfig struct with an offset of 0. Offsets are used
	// to make sure Telegram knows we've handled previous values and we don't
	// need them repeated.
	updateConfig := tgbotapi.NewUpdate(0)

	// Tell Telegram we should wait up to 30 seconds on each request for an
	// update. This way we can get information just as quickly as making many
	// frequent requests without having to send nearly as many.
	updateConfig.Timeout = 30

	// Start polling Telegram for updates.
	updates := bot.GetUpdatesChan(updateConfig)

	// Let's go through each update that we're getting from Telegram.
	for update := range updates {
		// Telegram can send many types of updates depending on what your Bot
		// is up to. We only want to look at messages for now, so we can
		// discard any other updates.
		if update.Message == nil {
			continue
		}

		chatIDS, err = saveChatID(chatIDS, update.Message.Chat.ID)
		if err != nil {
			log.Fatal(err)
		}

		// Now that we know we've gotten a new message, we can construct a
		// reply! We'll take the Chat ID and Text from the incoming message
		// and use it to create a new message.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вы подписалиль на получение главных новостей дня! Первое сообщение придет в течение суток.")

		bot.Send(msg)
	}
}

func saveChatID(chatIDS []int64, chatID int64) ([]int64, error) {
	for _, id := range chatIDS {
		if id == chatID {
			return chatIDS, nil // chatID already exists, return without modifying the slice
		}
	}

	// chatID does not exist, add it to the slice
	chatIDS = append(chatIDS, chatID)

	// Marshal the chatIDS slice to JSON
	data, err := json.Marshal(chatIDS)
	if err != nil {
		return nil, err
	}

	// Write the JSON data to a file
	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		return nil, err
	}

	return chatIDS, nil
}

func readChatIDS() (chatIDS []int64, err error) {
	_, err = os.Stat(filePath)

	if err != nil {
		// File does't exists
		return chatIDS, nil
	}

	// Read the JSON data from the file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into a []int64 slice
	err = json.Unmarshal(data, &chatIDS)
	if err != nil {
		return nil, err
	}

	return chatIDS, nil
}
