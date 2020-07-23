package model

type JobOffer struct {
	ID           string
	Portal       string
	URL          string
	LinkText     string
	Title        string
	Date         string
	Location     string
	Keywords     []string
	Seniority    string
	Salary       string
	Description  string
	Mode         string
	Organization string
}
