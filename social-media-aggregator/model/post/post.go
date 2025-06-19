package post

type Post struct {
	ID        string `json:"id" db:"id"`
	CreatedAt string `json:"created_at" db:"created_at"`
	URL       string `json:"url" db:"url"`
	Language  string `json:"language" db:"language"`
	Content   string `json:"content" db:"content"`
}
