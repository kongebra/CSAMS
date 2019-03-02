package model

import "github.com/JohanAanesen/NTNU-Bachelor-Management-System-For-CS-Assignments/webservice/shared/db"

// Review struct
type Review struct {
	ID     int   `json:"id" db:"id"`
	FormID int   `json:"-" db:"form_id"`
	Form   Form `json:"form"`
}

type ReviewRepository struct {}

func (repo *ReviewRepository) Insert(form Form) error {
	// Declare FormRepository
	var formRepo = FormRepository{}

	// Insert form to database, and get last inserted Id from it
	formID, err := formRepo.Insert(form)
	if err != nil {
		return err
	}

	// Insertions query
	query := "INSERT INTO reviews (form_id) VALUES(?)"
	// Insert form_id into submissions
	rows, err := db.GetDB().Query(query, formID)
	// Check for error
	if err != nil {
		return err
	}
	// Close the connections
	defer rows.Close()

	// Loop trough fields in the forms
	for _, field := range form.Fields {
		// Insertion query
		query := "INSERT INTO fields (form_id, type, name, label, description, priority, weight, choices) VALUES (?, ?, ?, ?, ?, ?, ?, ?);"
		// Execute the query
		rows, err := db.GetDB().Query(query, int(formID), field.Type, field.Name, field.Label, field.Description, field.Order, field.Weight, field.Choices)
		// Check for error
		if err != nil {
			return err
		}
		// Close the connection
		rows.Close()
	}

	// Return no error
	return nil
}

// GetAll returns all submission in the database
func (repo *ReviewRepository) GetAll() ([]Review, error) {
	// Declare return slice
	var result []Review
	// Create query-string
	query := "SELECT id, form_id FROM review"
	// Perform query
	rows, err := db.GetDB().Query(query)
	// Check for error
	if err != nil {
		return result, err
	}
	// Close connection
	defer rows.Close()

	// Loop through rows
	for rows.Next() {
		// Declare a single Submission
		var review = Review{}
		// Scan the data from the rows
		err = rows.Scan(&review.ID, &review.FormID)
		// Check for error
		if err != nil {
			return result, err
		}

		// Append scan-result to result
		result = append(result, review)
	}

	// Declare a FormRepository
	var formRepo = FormRepository{}
	// Loop through all submissions
	for index, review := range result {
		// Get form from database
		form, err := formRepo.Get(review.FormID)
		// Check for error
		if err != nil {
			return result, nil
		}
		// Get the form
		review.Form = form
		// Set the new values
		result[index] = review
	}

	return result, nil
}

func (repo *ReviewRepository) Update(form Form) error {
	query := "UPDATE forms SET prefix=?, name=?, description=? WHERE id=?"
	tx, err := db.GetDB().Begin()
	if err != nil {
		return err
	}

	_, err = db.GetDB().Exec(query, form.Prefix, form.Name, form.Description, form.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	query = "DELETE FROM fields WHERE form_id=?"
	rows, err := db.GetDB().Query(query, form.ID)
	if err != nil {
		return err
	}
	rows.Close()

	// Loop trough fields in the forms
	for _, field := range form.Fields {
		// Insertion query
		query := "INSERT INTO fields (form_id, type, name, label, description, priority, weight, choices) VALUES (?, ?, ?, ?, ?, ?, ?, ?);"
		// Execute the query
		rows, err := db.GetDB().Query(query, form.ID, field.Type, field.Name, field.Label, field.Description, field.Order, field.Weight, field.Choices)
		// Check for error
		if err != nil {
			return err
		}
		// Close the connection
		rows.Close()
	}

	// Return no error
	return nil
}