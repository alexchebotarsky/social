package post

type Post struct {
	ID        string   `json:"id"`
	CreatedAt string   `json:"created_at"`
	Language  string   `json:"language"`
	URL       string   `json:"url"`
	Content   string   `json:"content"`
	Account   *Account `json:"account"`
}

type Account struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Acct        string `json:"acct"`
	URL         string `json:"url"`
}
