package model

// Course holds the data for courses
type Course struct {
	ID          int          `json:"id"`
	Hash        string       `json:"hash"`
	Code        string       `json:"code"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Teacher     int          `json:"teacher"`
	Year        string       `json:"year"`
	Semester    string       `json:"semester"`
	Assignments []Assignment `json:"assignments"`
}
