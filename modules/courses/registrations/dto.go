package registrations

type RegisterCourseInput struct {
	Answers []AnswerInput `json:"answers"`
}

type AnswerInput struct {
	QuestionID uint   `json:"questionId" binding:"required"`
	AnswerText string `json:"answerText" binding:"required"`
}

type UpdateRegistrationStatusInput struct {
	Status string `json:"status" binding:"required,oneof=approved rejected"`
}

type RegistrationResponse struct {
	ID                uint             `json:"id"`
	UserID            uint             `json:"userId"`
	CourseID          uint             `json:"courseId"`
	Status            string           `json:"status"`
	CourseName        string           `json:"courseName,omitempty"`
	ProgramName       string           `json:"programName,omitempty"`
	WhatsappGroupLink string           `json:"whatsappGroupLink,omitempty"`
	UserName          string           `json:"userName,omitempty"`
	UserEmail         string           `json:"userEmail,omitempty"`
	Answers           []AnswerResponse `json:"answers,omitempty"`
	CreatedAt         string           `json:"createdAt"`
}

type AnswerResponse struct {
	QuestionID   uint   `json:"questionId"`
	QuestionText string `json:"questionText"`
	AnswerText   string `json:"answerText"`
}

type RegistrationListResponse struct {
	Registrations []RegistrationResponse `json:"registrations"`
	TotalCount    int64                  `json:"totalCount"`
}
