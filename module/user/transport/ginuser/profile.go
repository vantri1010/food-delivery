package ginuser

import (
	"food-delivery/common"
	"food-delivery/component/appctx"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Profile(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		u := c.MustGet(common.CurrentUser) // middleware get info from context
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(u))
	}
}
