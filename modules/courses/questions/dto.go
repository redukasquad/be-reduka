package questions

type CreateQuestionInput struct {
	QuestionText  string `json:"questionText" binding:"required"`
	QuestionType  string `json:"questionType" binding:"required,oneof=text select radio checkbox"`
	QuestionOrder int    `json:"questionOrder" binding:"required,min=1"`
}

type UpdateQuestionInput struct {
	QuestionText  *string `json:"questionText"`
	QuestionType  *string `json:"questionType"`
	QuestionOrder *int    `json:"questionOrder"`
}

type QuestionResponse struct {
	ID            uint   `json:"id"`
	CourseID      uint   `json:"courseId"`
	QuestionText  string `json:"questionText"`
	QuestionType  string `json:"questionType"`
	QuestionOrder int    `json:"questionOrder"`
}
