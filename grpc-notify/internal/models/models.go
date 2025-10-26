package models

type MessageTelegram struct {
	ChatId  int    `json:"chat_id"`
	Message string `json:"text"`
}
