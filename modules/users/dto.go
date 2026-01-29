package users

type UpdateUserInput struct {
	Username     *string `json:"username"`
	NoTelp       *string `json:"no_telp"`
	JenisKelamin *bool   `json:"jenis_kelamin"`
	Kelas        *string `json:"kelas" binding:"omitempty,oneof='Kelas 10' 'Kelas 11' 'Kelas 12' 'Gapyear (Alumni)'"`
	ProfileImage *string `json:"profile_image"`
}

// SetRoleInput is used by admin to set user roles
type SetRoleInput struct {
	Role string `json:"role" binding:"required,oneof=STUDENT TUTOR ADMIN"`
}
