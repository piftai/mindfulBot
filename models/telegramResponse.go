package models

type TelegramResponse struct {
	Ok     bool `json:"ok"`
	Result struct {
		ID        int    `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string `json:"username"`
		Type      string `json:"type"`
	} `json:"result"`
}
