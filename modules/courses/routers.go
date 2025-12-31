package courses

import (
	"github.com/gin-gonic/gin"

	"github.com/redukasquad/be-reduka/modules/courses/answers"
	courses "github.com/redukasquad/be-reduka/modules/courses/index"
	"github.com/redukasquad/be-reduka/modules/courses/questions"
	"github.com/redukasquad/be-reduka/modules/courses/registrations"
)

func CoursesRouter(router *gin.RouterGroup, requireAuth gin.HandlerFunc, requireAdmin gin.HandlerFunc) {
	courses.CourseIndexRouter(router, requireAuth, requireAdmin)
	registrations.RegistrationRouter(router, requireAuth, requireAdmin)
	questions.QuestionRouter(router, requireAuth, requireAdmin)
	answers.AnswerRouter(router, requireAuth, requireAdmin)
}
