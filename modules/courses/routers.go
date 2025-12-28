package courses

import (
	"github.com/gin-gonic/gin"

	"github.com/redukasquad/be-reduka/modules/courses/answers"
	courses "github.com/redukasquad/be-reduka/modules/courses/index"
	"github.com/redukasquad/be-reduka/modules/courses/questions"
	"github.com/redukasquad/be-reduka/modules/courses/registrations"
)

// CoursesRouter registers all course-related routes
// This combines all sub-module routers: index, registrations, questions, answers
func CoursesRouter(router *gin.RouterGroup, requireAuth gin.HandlerFunc, requireAdmin gin.HandlerFunc) {
	// Course CRUD routes
	courses.CourseIndexRouter(router, requireAuth, requireAdmin)

	// Registration routes (user registration, admin approval/rejection)
	registrations.RegistrationRouter(router, requireAuth, requireAdmin)

	// Question routes (registration form questions)
	questions.QuestionRouter(router, requireAuth, requireAdmin)

	// Answer routes (view registration answers)
	answers.AnswerRouter(router, requireAuth, requireAdmin)
}
