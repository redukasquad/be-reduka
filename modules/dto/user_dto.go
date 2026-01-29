package dto

import "github.com/redukasquad/be-reduka/database/entities"

type CreatorResponse struct {
	ID           uint   `json:"id"`
	Username     string `json:"username"`
	ProfileImage string `json:"profileImage,omitempty"`
}

func ToCreatorResponse(user entities.User) CreatorResponse {
	return CreatorResponse{
		ID: user.ID,
		Username: user.Username,
		ProfileImage: user.ProfileImage,
	}
}