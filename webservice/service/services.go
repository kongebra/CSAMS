package service

import "database/sql"

// Services struct
type Services struct {
	Assignment       *AssignmentService
	Course           *CourseService
	Field            *FieldService
	Form             *FormService
	Review           *ReviewService
	ReviewAnswer     *ReviewAnswerService
	Submission       *SubmissionService
	SubmissionAnswer *SubmissionAnswerService
	User             *UserService
}

// NewServices func
func NewServices(db *sql.DB) *Services {
	return &Services{
		Assignment:       NewAssignmentService(db),
		Course:           NewCourseService(db),
		Field:            NewFieldService(db),
		Form:             NewFormService(db),
		Review:           NewReviewService(db),
		ReviewAnswer:     NewReviewAnswerService(db),
		Submission:       NewSubmissionService(db),
		SubmissionAnswer: NewSubmissionAnswerService(db),
		User:             NewUserService(db),
	}
}