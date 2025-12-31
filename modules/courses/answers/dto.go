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
