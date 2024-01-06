package main

type JsonRequest struct {
	Completed      bool     `json:"completed"`
	CreatedDate    string   `json:"created_date"`
	DueDate        string   `json:"due_date"`
	Priority       string   `json:"priority"`
	Projects       []string `json:"projects"`
	Contexts       []string `json:"contexts"`
	AdditionalTags []string `json:"additional_tags"`
	CompletedDate  string   `json:"completed_date"`
}
