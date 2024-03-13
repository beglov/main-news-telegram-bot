package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const chatIDSFilePath = "chatIDS.json"
const newsIDSFilePath = "newsIDS.json"
const channel = "lentadnya"
const offset = "0"

type Message struct {
	ID       int    `json:"id"`
	Date     int    `json:"date"`
	Views    int    `json:"views"`
	Forwards int    `json:"forwards"`
	EditDate int    `json:"edit_date"`
	Text     string `json:"text"`
	Html     string `json:"html"`
	Photo    string `json:"photo"`
}

type Response struct {
	Count    int       `json:"count"`
	Messages []Message `json:"messages"`
}

func main() {
	godotenv.Load(".env")

	chatIDS, err := readChatIDS()
	if err != nil {
		log.Fatal(err)
	}
	if chatIDS == nil {
		log.Fatal("chatIDS is empty")
	}

	newsIDS, err := readNewsIDS()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	news, err := getNews()
	if err != nil {
		log.Fatal(err)
	}

	for _, message := range news {
		for _, chatID := range chatIDS {
			if alreadySend(newsIDS, message) {
				continue
			}

			msg := tgbotapi.NewMessage(chatID, message.Html)
			msg.ParseMode = "HTML"
			_, err = bot.Send(msg)
			if err == nil {
				newsIDS, err = saveNewsID(newsIDS, message.ID)
				if err != nil {
					log.Print(err)
				}
			} else {
				log.Print(err)
			}
		}
	}
}

func alreadySend(newsIDS []int, message Message) bool {
	for _, id := range newsIDS {
		if id == message.ID {
			return true
		}
	}

	return false
}

func readChatIDS() (chatIDS []int64, err error) {
	_, err = os.Stat(chatIDSFilePath)

	// File doesn't exists
	if err != nil {
		return chatIDS, nil
	}

	// Read the JSON data from the file
	data, err := ioutil.ReadFile(chatIDSFilePath)
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

func readNewsIDS() (newsIDS []int, err error) {
	_, err = os.Stat(newsIDSFilePath)

	// File doesn't exists
	if err != nil {
		return newsIDS, nil
	}

	// Read the JSON data from the file
	data, err := ioutil.ReadFile(newsIDSFilePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into a []int64 slice
	err = json.Unmarshal(data, &newsIDS)
	if err != nil {
		return nil, err
	}

	return newsIDS, nil
}

func saveNewsID(newsIDS []int, newsID int) ([]int, error) {
	for _, id := range newsIDS {
		if id == newsID {
			return newsIDS, nil // newsID already exists, return without modifying the slice
		}
	}

	// newsID does not exist, add it to the slice
	newsIDS = append(newsIDS, newsID)

	// Marshal the newsIDS slice to JSON
	data, err := json.Marshal(newsIDS)
	if err != nil {
		return nil, err
	}

	// Write the JSON data to a file
	err = ioutil.WriteFile(newsIDSFilePath, data, 0644)
	if err != nil {
		return nil, err
	}

	return newsIDS, nil
}

// getNews возвращает главные новости ленты дня - это главное за ночь и главное за день.
func getNews() ([]Message, error) {
	url := fmt.Sprintf("https://telegram92.p.rapidapi.com/api/history/channel?channel=%s&limit=%s&offset=%s", channel, os.Getenv("POSTS_COUNT"), offset)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-RapidAPI-Key", os.Getenv("RAPID_API_KEY"))
	req.Header.Add("X-RapidAPI-Host", "telegram92.p.rapidapi.com")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("Rapid API response:\n%s", string(body))

	var data Response
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	messages := reverseSlice(data.Messages)

	var news []Message
	for _, message := range messages {
		str := message.Html
		if str == "" {
			continue
		}
		if !isMainNews(str) {
			continue
		}
		// Remove <br /> and <br> tags from the string
		re := regexp.MustCompile(`<br\s+/?>`)
		cleanStr := re.ReplaceAllString(str, "")
		post := Message{
			ID:   message.ID,
			Html: cleanStr,
		}
		news = append(news, post)
	}

	return news, nil
}

func reverseSlice(slice []Message) []Message {
	reversedSlice := make([]Message, len(slice))
	lastIndex := len(slice) - 1

	for i, value := range slice {
		reversedSlice[lastIndex-i] = value
	}

	return reversedSlice
}

func isMainNews(str string) bool {
	return strings.Count(str, "</a>") > 3
}
