package classes

import (
	"github.com/gin-gonic/gin"

	"github.com/redukasquad/be-reduka/modules/classes/lessons"
	"github.com/redukasquad/be-reduka/modules/classes/resources"
	"github.com/redukasquad/be-reduka/modules/classes/subjects"
)

func ClassesRouter(router *gin.RouterGroup, requireAuth gin.HandlerFunc, requireAdminOrTutor gin.HandlerFunc) {
	subjects.SubjectRouter(router, requireAuth, requireAdminOrTutor)
	lessons.LessonRouter(router, requireAuth, requireAdminOrTutor)
	resources.ResourceRouter(router, requireAuth, requireAdminOrTutor)
}
