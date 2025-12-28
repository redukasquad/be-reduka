package registrations

// RegisterCourseInput is the input for registering to a course
type RegisterCourseInput struct {
	Answers []AnswerInput `json:"answers"`
}

// AnswerInput represents an answer to a registration question
type AnswerInput struct {
	QuestionID uint   `json:"questionId" binding:"required"`
	AnswerText string `json:"answerText" binding:"required"`
}

// UpdateRegistrationStatusInput is the input for updating registration status
type UpdateRegistrationStatusInput struct {
	Status string `json:"status" binding:"required,oneof=approved rejected"`
}

// RegistrationResponse is the response format for a registration
type RegistrationResponse struct {
	ID                uint             `json:"id"`
	UserID            uint             `json:"userId"`
	CourseID          uint             `json:"courseId"`
	Status            string           `json:"status"`
	CourseName        string           `json:"courseName,omitempty"`
	ProgramName       string           `json:"programName,omitempty"`
	WhatsappGroupLink string           `json:"whatsappGroupLink,omitempty"` // Only shown when approved
	UserName          string           `json:"userName,omitempty"`
	UserEmail         string           `json:"userEmail,omitempty"`
	Answers           []AnswerResponse `json:"answers,omitempty"`
	CreatedAt         string           `json:"createdAt"`
}

// AnswerResponse is the response format for a registration answer
type AnswerResponse struct {
	QuestionID   uint   `json:"questionId"`
	QuestionText string `json:"questionText"`
	AnswerText   string `json:"answerText"`
}

// RegistrationListResponse is the response for listing registrations
type RegistrationListResponse struct {
	Registrations []RegistrationResponse `json:"registrations"`
	TotalCount    int64                  `json:"totalCount"`
}
