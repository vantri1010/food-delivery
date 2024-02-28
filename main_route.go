package main

import (
	"food-delivery/component/appctx"
	"food-delivery/middleware"
	"food-delivery/module/restaurant/transport/ginrestaurant"
	"food-delivery/module/restaurantlike/transport/ginrstlike"
	"food-delivery/module/upload/uploadtransport/ginupload"
	"food-delivery/module/user/transport/ginuser"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func setupRoute(appContext appctx.AppContext, v1 *gin.RouterGroup) {
	v1.POST("/upload", ginupload.UpLoad(appContext))

	v1.POST("/register", ginuser.Register(appContext))
	v1.POST("/authenticate", ginuser.Login(appContext))
	v1.GET("/profile", middleware.RequiredAuth(appContext), ginuser.Profile(appContext))

	restaurants := v1.Group("/restaurants", middleware.RequiredAuth(appContext))

	restaurants.POST("", ginrestaurant.CreateRestaurant(appContext))

	restaurants.GET("/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		var data Restaurant

		appContext.GetMainDBConnection().Where("id = ?", id).First(&data)

		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	})

	restaurants.GET("", ginrestaurant.ListRestaurant(appContext))

	restaurants.PATCH("/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		var data RestaurantUpdate

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		appContext.GetMainDBConnection().Where("id = ?", id).Updates(&data)

		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	})

	restaurants.DELETE("/:id", ginrestaurant.DeleteRestaurant(appContext))

	restaurants.POST("/:id/like", ginrstlike.UserLikeRestaurant(appContext))
	restaurants.DELETE("/:id/dislike", ginrstlike.UserDislikeRestaurant(appContext))
	restaurants.GET("/:id/liked-users", ginrstlike.ListUsers(appContext))
}
