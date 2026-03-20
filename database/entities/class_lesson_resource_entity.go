package entities

import "gorm.io/gorm"

type LessonResource struct {
	gorm.Model

	LessonID uint `json:"lessonId"`

	Type  string `json:"type" binding:"required,oneof=video document link zoom recording" gorm:"type:enum('video','document','link','zoom','recording')"`
	Title string `json:"title"`
	URL   string `json:"url" binding:"required"`

	Lesson Lesson `json:"lesson,omitempty" gorm:"foreignKey:LessonID"`
}
