package dto

type RegisterInput struct {
	Username     string `json:"username" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=6"`
	Role         string `json:"role" binding:"required,oneof=Students Tutor Admin"`
	NoTelp       string `json:"no_telp" binding:"required"`
	JenisKelamin bool   `json:"jenis_kelamin"`
	Kelas        string `json:"kelas" binding:"required,oneof='Kelas 10' 'Kelas 11' 'Kelas 12' 'Gapyear (Alumni)'"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
