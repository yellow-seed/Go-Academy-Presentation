package user

type User struct {
	Id           int    `json:"id"`
	LineUserId   string `json:"line_user_id"`
	LanguageCode string `json:"language_code"`
	SearchMode   string `json:"search_mode"`
}