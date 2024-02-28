package ginrstlike

import (
	"food-delivery/common"
	"food-delivery/component/appctx"
	rstlikebiz "food-delivery/module/restaurantlike/biz"
	restaurantlikemodel "food-delivery/module/restaurantlike/model"
	restaurantlikestore "food-delivery/module/restaurantlike/store"
	"github.com/gin-gonic/gin"
	"net/http"
)

// POST/v1/restaurants/:id/like

func UserLikeRestaurant(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, err := common.FromBase58(c.Param("id"))

		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		requester := c.MustGet(common.CurrentUser).(common.Requester)

		data := restaurantlikemodel.Like{
			RestaurantId: int(uid.GetLocalID()),
			UserId:       requester.GetUserId(),
		}

		store := restaurantlikestore.NewSQLstore(appCtx.GetMainDBConnection())
		biz := rstlikebiz.NewUserLikeRestaurantBiz(store)

		if err := biz.LikeRestaurant(c.Request.Context(), &data); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data))
	}
}
