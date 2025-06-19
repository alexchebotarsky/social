package post

type Post struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	URL       string `json:"url"`
	Language  string `json:"language"`
	Content   string `json:"content"`
}
