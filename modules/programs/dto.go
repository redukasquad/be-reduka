package programs

// CreateProgramInput is the input for creating a new program
type CreateProgramInput struct {
	ProgramName  string `json:"programName" binding:"required"`
	Description  string `json:"description"`
	ImageProgram string `json:"imageProgram"`
}

// UpdateProgramInput is the input for updating a program
type UpdateProgramInput struct {
	ProgramName  *string `json:"programName"`
	Description  *string `json:"description"`
	ImageProgram *string `json:"imageProgram"`
}

// ProgramResponse is the response format for a program
type ProgramResponse struct {
	ID           uint   `json:"id"`
	ProgramName  string `json:"programName"`
	Description  string `json:"description"`
	ImageProgram string `json:"imageProgram"`
}
