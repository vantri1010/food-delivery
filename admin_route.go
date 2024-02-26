package main

import (
	"food-delivery/component/appctx"
	"food-delivery/middleware"
	"food-delivery/module/user/transport/ginuser"
	"github.com/gin-gonic/gin"
)

func setupAdminRoute(appContext appctx.AppContext, v1 *gin.RouterGroup) {
	admin := v1.Group("/admin",
		middleware.RequiredAuth(appContext),
		middleware.RoleRequired(appContext, "admin", "mod"),
	)

	{
		// local scope can get outside var (admin) but in the reverse is not true
		admin.GET("/profile", ginuser.Profile(appContext))
	}
}
