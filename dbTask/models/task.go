package models

type TaskStatus string

const (
	Done       TaskStatus = "done"
	Stale      TaskStatus = "stale"
	InProgress TaskStatus = "in-progress"
)

type Function struct {
	URL string `json:"url"`
}

type Task struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	Function    []Function `json:"function"`
}
