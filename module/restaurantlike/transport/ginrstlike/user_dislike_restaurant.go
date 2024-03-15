package ginrstlike

import (
	"food-delivery/common"
	"food-delivery/component/appctx"
	rstlikebiz "food-delivery/module/restaurantlike/biz"
	restaurantlikestore "food-delivery/module/restaurantlike/store"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserDislikeRestaurant(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, err := common.FromBase58(c.Param("id"))

		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		requester := c.MustGet(common.CurrentUser).(common.Requester)

		store := restaurantlikestore.NewSQLstore(appCtx.GetMainDBConnection())
		//decStore := restaurantstorage.NewSQLStore(appCtx.GetMainDBConnection())
		biz := rstlikebiz.NewUserDislikeRestaurantBiz(store, appCtx.GetPubSub())

		if err := biz.DislikeRestaurant(c.Request.Context(), requester.GetUserId(), int(uid.GetLocalID())); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
