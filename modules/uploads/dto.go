package uploads

import "time"

type CreateImageInput struct {
	URL    string `json:"url" binding:"required"`
	Fileid string `json:"fileId" binding:"required"`
}

type DeleteImageInput struct {
	URL string `json:"url" binding:"required"`
}

type ImageResponse struct {
	ID        uint      `json:"id"`
	URL       string    `json:"url"`
	Fileid    string    `json:"fileId"`
	CreatedAt time.Time `json:"createdAt"`
}
