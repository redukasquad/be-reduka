package answers

// AnswerResponse is the response format for a registration answer
type AnswerResponse struct {
	ID             uint   `json:"id"`
	RegistrationID uint   `json:"registrationId"`
	QuestionID     uint   `json:"questionId"`
	QuestionText   string `json:"questionText"`
	AnswerText     string `json:"answerText"`
}

// AnswerListResponse is the response for listing answers
type AnswerListResponse struct {
	Answers []AnswerResponse `json:"answers"`
}
