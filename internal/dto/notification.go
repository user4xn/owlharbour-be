package dto

type (
	Notification struct {
		Title  string `json:"title"`
		Body   string `json:"body"`
		Tokens []Tokens
	}

	Tokens struct {
		Token string `json:"token"`
	}
)
