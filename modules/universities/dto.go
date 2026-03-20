package universities

import "github.com/redukasquad/be-reduka/database/entities"

type CreateUniversityInput struct {
	Name string `json:"name" binding:"required"`
	Type string `json:"type" binding:"required,oneof=PTN PTS PTK"`
}

type UpdateUniversityInput struct {
	Name *string `json:"name"`
	Type *string `json:"type" binding:"omitempty,oneof=PTN PTS PTK"`
}

type CreateMajorInput struct {
	UniversityID uint    `json:"universityId"`
	Name         string  `json:"name" binding:"required"`
	PassingGrade float64 `json:"passingGrade" binding:"min=0,max=100"`
}

type UpdateMajorInput struct {
	Name         *string  `json:"name"`
	PassingGrade *float64 `json:"passingGrade" binding:"omitempty,min=0,max=100"`
}

type SetUserTargetInput struct {
	UniversityMajorID uint `json:"universityMajorId" binding:"required"`
	Priority          int  `json:"priority" binding:"required,min=1"`
}

// ── Response DTOs (lowercase id for frontend consistency) ────────────────────

type MajorResponse struct {
	ID           uint                `json:"id"`
	UniversityID uint                `json:"universityId"`
	Name         string              `json:"name"`
	PassingGrade float64             `json:"passingGrade"`
	University   *UniversityResponse `json:"university,omitempty"`
}

type UniversityResponse struct {
	ID       uint            `json:"id"`
	Name     string          `json:"name"`
	Type     string          `json:"type"`
	Programs []MajorResponse `json:"programs,omitempty"`
}

type UserResponse struct {
	ID           uint   `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Kelas        string `json:"kelas"`
	ProfileImage string `json:"profileImage"`
}

type UserTargetResponse struct {
	ID                uint           `json:"id"`
	UserID            uint           `json:"userId"`
	UniversityMajorID uint           `json:"universityMajorId"`
	Priority          int            `json:"priority"`
	UniversityProgram *MajorResponse `json:"universityProgram,omitempty"`
}

// ── Mapping helpers ───────────────────────────────────────────────────────────

func ToMajorResponse(m entities.UniversityMajor) MajorResponse {
	resp := MajorResponse{
		ID:           m.ID,
		UniversityID: m.UniversityID,
		Name:         m.Name,
		PassingGrade: m.PassingGrade,
	}
	if m.University.ID != 0 {
		uniResp := UniversityResponse{
			ID:   m.University.ID,
			Name: m.University.Name,
			Type: m.University.Type,
		}
		resp.University = &uniResp
	}
	return resp
}

func ToUniversityResponse(u entities.University) UniversityResponse {
	resp := UniversityResponse{
		ID:   u.ID,
		Name: u.Name,
		Type: u.Type,
	}
	for _, m := range u.Major {
		resp.Programs = append(resp.Programs, ToMajorResponse(m))
	}
	return resp
}

func ToUserTargetResponse(t entities.UserTarget) UserTargetResponse {
	resp := UserTargetResponse{
		ID:                t.ID,
		UserID:            t.UserID,
		UniversityMajorID: t.UniversityMajorID,
		Priority:          t.Priority,
	}
	if t.Major.ID != 0 {
		mr := ToMajorResponse(t.Major)
		resp.UniversityProgram = &mr
	}
	return resp
}

func ToUserResponse(u entities.User) UserResponse {
	kelas := ""
	if u.Kelas != nil {
		kelas = *u.Kelas
	}
	return UserResponse{
		ID:           u.ID,
		Username:     u.Username,
		Email:        u.Email,
		Kelas:        kelas,
		ProfileImage: u.ProfileImage,
	}
}
