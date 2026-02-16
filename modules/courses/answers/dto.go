package answers

type AnswerResponse struct {
	ID             uint   `json:"id"`
	RegistrationID uint   `json:"registrationId"`
	QuestionID     uint   `json:"questionId"`
	QuestionText   string `json:"questionText"`
	AnswerText     string `json:"answerText"`
}

type AnswerListResponse struct {
	Answers []AnswerResponse `json:"answers"`
}

type CreateAnswerRequest struct {
	RegistrationID uint   `json:"registrationId" binding:"required"`
	QuestionID     uint   `json:"questionId" binding:"required"`
	AnswerText     string `json:"answerText" binding:"required"`
}
