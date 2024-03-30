package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type MisskeyPostRequest struct {
	I          string `json:"i"`
	Text       string `json:"text"`
	Visibility string `json:"visibility"`
}

func main() {
	logFile, err := os.OpenFile("jihoulog.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("ログファイルを開くのに失敗しました: ", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal(".env の読み取りに失敗しました。: ", err)
	}

	host := os.Getenv("MISSKEY_HOST")
	token := os.Getenv("MISSKEY_TOKEN")
	content := os.Getenv("MISSKEY_CONTENT")

	if host == "" || token == "" || content == "" {
		log.Fatal("MISSKEY_HOST、MISSKEY_TOKEN、またはMISSKEY_CONTENTが設定されていません")
	}

	postRequest := MisskeyPostRequest{
		I:          token,
		Text:       content,
		Visibility: "home",
	}
	postData, err := json.Marshal(postRequest)
	if err != nil {
		log.Fatal("JSONのマーシャリングに失敗しました: ", err)
	}

	url := fmt.Sprintf("https://%s/api/notes/create", host)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postData))
	if err != nil {
		log.Fatal("リクエストの作成に失敗しました: ", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("リクエストの送信に失敗しました: ", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("リクエスト失敗。ステータスコード: %d", resp.StatusCode)
	}

	log.Println("つぶやきに成功しました")
}
