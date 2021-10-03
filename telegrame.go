package telegrame

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/vianamjr/telegrame/config"
)

const (
	ENV_URI     = "TELEGRAME_URI"
	ENV_TOKEN   = "TELEGRAME_TOKEN"
	ENV_CHAT_ID = "TELEGRAME_CHAT_ID"
)

type BOT struct {
	token     string
	baseURI   string
	botMethod string
	chatID    string
}

type messageModel struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func NewBot() (*BOT, error) {
	err := config.LoadConfig()
	if err != nil {
		log.Println("error on config:", err.Error())
		return nil, err
	}
	fmt.Println("env uri", os.Getenv(ENV_URI))
	return &BOT{
		baseURI: os.Getenv(ENV_URI),
		token:   os.Getenv(ENV_TOKEN),
		chatID:  os.Getenv(ENV_CHAT_ID),
	}, nil
}

func (b *BOT) SendMessage(message string) error {
	b.botMethod = "sendMessage"
	uri := fmt.Sprintf("%v%v/%v", b.baseURI, b.token, b.botMethod)
	b.doRequest(uri, message)
	return nil
}

func (b *BOT) doRequest(uri string, message string) {
	m := messageModel{
		ParseMode: "HTML",
		ChatID:    b.chatID,
		Text:      message,
	}

	payload, err := json.Marshal(m)
	if err != nil {
		fmt.Println("error on marshal message:", err.Error())
		return
	}

	request, err := http.NewRequest(http.MethodPost, uri, bytes.NewReader(payload))
	if err != nil {
		fmt.Println("request :", err)
		return
	}

	request.Header.Add("Content-Type", "application/json")
	fmt.Println(uri)
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("do :", err)
		return
	}
	defer response.Body.Close()

	log.Println("sataus code :", response.StatusCode)

	if response.StatusCode >= 300 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println("body :", err)
			return
		}
		fmt.Println("--------------- erro ini", response.StatusCode, "---------------")
		fmt.Println(string(body))
		fmt.Println()
		fmt.Println(string(payload))
		fmt.Println()
		fmt.Println(response.Header)
		fmt.Println("--------------- erro fim", response.StatusCode, "---------------")
	}

}
