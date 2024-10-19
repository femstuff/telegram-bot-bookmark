package telegram

type Update struct {
	UpdateID int    `json:"update_id"`
	Message  string `json:"message"`
}

type UpdateResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}
