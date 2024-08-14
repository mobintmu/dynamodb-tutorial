package repository

import "time"

type User struct {
	Uid             string `json:"user_id"`
	TgId            string `json:"tg_id"`
	Referrer        string `json:"referrer" `
	Name            string `json:"name" `
	TgFirstName     string `json:"tg_first_name"`
	TgLastName      string `json:"tg_last_name"`
	TgUsername      string `json:"tg_username"`
	ProfilePicture  string `json:"profile_picture"`
	CounterReferrer int    `json:"counter_referrer"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
