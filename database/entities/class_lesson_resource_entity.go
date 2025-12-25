package entities

import "gorm.io/gorm"

type ClassLessonResource struct {
	gorm.Model
	ClassLessonID uint

	Type string `json:"type" binding:"required,oneof=video document link zoom recording" gorm:"type:enum('video','document','link','zoom','recording')"`

	Title string `json:"title"`
	URL   string `json:"url" binding:"required"`

	ClassLesson ClassLesson `gorm:"foreignKey:ClassLessonID"`
}
