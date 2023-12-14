package main

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

type WebhookData struct {
	Event         string `json:"event"`
	CallerID      string `json:"caller_id"`
	CalledDID     string `json:"called_did"`
	CallStart     string `json:"call_start"`
	Duration      int    `json:"duration,omitempty"`
	Disposition   string `json:"disposition,omitempty"`
	StatusCode    string `json:"status_code,omitempty"`
	IsRecorded    bool   `json:"is_recorded,omitempty"`
	Internal      bool   `json:"internal,omitempty"`
	LastInternal  string `json:"last_internal,omitempty"`
	CallIDWithRec string `json:"call_id_with_rec,omitempty"`
	Destination   string `json:"destination,omitempty"`
}

func sendWebhook(url string, data WebhookData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func main() {
	webhookURL := "http://yourserver.com/webhook" // Замените на ваш URL
	callID := "unique-call-id"                    // Генерируем уникальный ID для звонка
	callerID := "123456789"                       // Пример Caller ID

	// Отправляем NOTIFY_ANSWER
	answerData := WebhookData{
		Event:         "NOTIFY_ANSWER",
		CallerID:      callerID,
		CalledDID:     "987654321",
		CallStart:     time.Now().Format(time.RFC3339),
		CallIDWithRec: callID,
	}
	err := sendWebhook(webhookURL, answerData)
	if err != nil {
		return
	}

	// Ждем 10 секунд перед отправкой NOTIFY_END
	time.Sleep(10 * time.Second)
	endData := WebhookData{
		Event:         "NOTIFY_END",
		CallerID:      callerID,
		CalledDID:     "987654321",
		CallStart:     answerData.CallStart,
		Duration:      600, // Продолжительность в секундах
		CallIDWithRec: callID,
	}
	err = sendWebhook(webhookURL, endData)
	if err != nil {
		return
	}

	// Ждем случайное время от 40 до 60 секунд перед отправкой NOTIFY_RECORD
	time.Sleep(time.Duration(40+rand.Intn(20)) * time.Second)
	recordData := WebhookData{
		Event:         "NOTIFY_RECORD",
		CallIDWithRec: callID,
	}

	err = sendWebhook(webhookURL, recordData)
	if err != nil {
		return
	}
}
