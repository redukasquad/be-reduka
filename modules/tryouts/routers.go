package tryouts

import (
	"github.com/gin-gonic/gin"

	"github.com/redukasquad/be-reduka/modules/tryouts/attempts"
	tryouts "github.com/redukasquad/be-reduka/modules/tryouts/index"
	"github.com/redukasquad/be-reduka/modules/tryouts/questions"
	"github.com/redukasquad/be-reduka/modules/tryouts/registrations"
)

func TryOutsRouter(router *gin.RouterGroup, requireAuth gin.HandlerFunc, requireAdmin gin.HandlerFunc) {
	tryouts.TryOutIndexRouter(router, requireAuth, requireAdmin)
	questions.QuestionRouter(router, requireAuth, requireAdmin)
	registrations.RegistrationRouter(router, requireAuth, requireAdmin)
	attempts.AttemptRouter(router, requireAuth, requireAdmin)
}
