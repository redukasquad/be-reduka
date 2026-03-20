package courses

import (
	"github.com/gin-gonic/gin"

	"github.com/redukasquad/be-reduka/modules/courses/answers"
	courses "github.com/redukasquad/be-reduka/modules/courses/index"
	"github.com/redukasquad/be-reduka/modules/courses/questions"
	"github.com/redukasquad/be-reduka/modules/courses/registrations"
)

func CoursesRouter(router *gin.RouterGroup, requireAuth gin.HandlerFunc, requireAdmin gin.HandlerFunc, requireAdminOrTutor gin.HandlerFunc) {
	courses.CourseIndexRouter(router, requireAuth, requireAdmin, requireAdminOrTutor)
	registrations.RegistrationRouter(router, requireAuth, requireAdmin, requireAdminOrTutor)
	questions.QuestionRouter(router, requireAuth, requireAdmin)
	answers.AnswerRouter(router, requireAuth, requireAdmin)
}
