package tryouts

import (
	"github.com/gin-gonic/gin"

	"github.com/redukasquad/be-reduka/modules/tryouts/attempts"
	tryouts "github.com/redukasquad/be-reduka/modules/tryouts/index"
	"github.com/redukasquad/be-reduka/modules/tryouts/questions"
	"github.com/redukasquad/be-reduka/modules/tryouts/registrations"
)

func TryOutsRouter(router *gin.RouterGroup, requireAuth gin.HandlerFunc, requireAdminOrTutor gin.HandlerFunc) {
	tryouts.TryOutIndexRouter(router, requireAuth, requireAdminOrTutor)
	questions.QuestionRouter(router, requireAuth, requireAdminOrTutor)
	registrations.RegistrationRouter(router, requireAuth, requireAdminOrTutor)
	attempts.AttemptRouter(router, requireAuth, requireAdminOrTutor)
}
