package questions

// CreateQuestionInput is the input for creating a registration question
type CreateQuestionInput struct {
	QuestionText  string `json:"questionText" binding:"required"`
	QuestionType  string `json:"questionType" binding:"required,oneof=text textarea select radio checkbox"`
	QuestionOrder int    `json:"questionOrder" binding:"required,min=1"`
	Options       string `json:"options"` // JSON string for select/radio/checkbox options
}

// UpdateQuestionInput is the input for updating a question
type UpdateQuestionInput struct {
	QuestionText  *string `json:"questionText"`
	QuestionType  *string `json:"questionType"`
	QuestionOrder *int    `json:"questionOrder"`
	Options       *string `json:"options"`
}

// QuestionResponse is the response format for a question
type QuestionResponse struct {
	ID            uint   `json:"id"`
	CourseID      uint   `json:"courseId"`
	QuestionText  string `json:"questionText"`
	QuestionType  string `json:"questionType"`
	QuestionOrder int    `json:"questionOrder"`
	Options       string `json:"options,omitempty"`
}
