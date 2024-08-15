package repository

type User struct {
	Uid             string `json:"user_id"`
	TgId            string `json:"tg_id"`
	Referrer        string `json:"referrer" `
	Name            string `json:"name" `
	TgFirstName     string `json:"tg_first_name"`
	TgLastName      string `json:"tg_last_name"`
	TgUsername      string `json:"tg_user_name"`
	ProfilePicture  string `json:"profile_picture"`
	CounterReferrer int    `json:"counter_referrer"`
	CreatedAt       int64  `json:"created_at"`
	UpdatedAt       int64  `json:"updated_at"`
}
