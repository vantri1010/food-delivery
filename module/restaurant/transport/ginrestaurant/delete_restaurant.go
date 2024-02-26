package ginrestaurant

import (
	"food-delivery/common"
	"food-delivery/component/appctx"
	"food-delivery/module/restaurant/biz"
	"food-delivery/module/restaurant/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DeleteRestaurant(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := appCtx.GetMainDBConnection()
		requester := c.MustGet(common.CurrentUser).(common.Requester)

		//id, err := strconv.Atoi(c.Param("id"))

		uid, err := common.FromBase58(c.Param("id"))

		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		store := restaurantstorage.NewSQLStore(db)
		biz := restaurantbiz.NewDeleteRestaurantBiz(store, requester)

		if err := biz.DeleteRestaurant(c.Request.Context(), int(uid.GetLocalID())); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
