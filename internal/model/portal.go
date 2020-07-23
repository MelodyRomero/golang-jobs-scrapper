package model

type Portal struct {
	Web      string
	Name     string
	BaseURL  string
	JobsURL  string
	Keywords []string
	Exclude  []string
}
