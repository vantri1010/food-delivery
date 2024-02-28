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

func ListUsers(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, err := common.FromBase58(c.Param("id"))

		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		filter := restaurantlikemodel.Filter{
			RestaurantId: int(uid.GetLocalID()),
		}

		var paging common.Paging

		if err := c.ShouldBind(&paging); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		paging.Fulfill()

		store := restaurantlikestore.NewSQLstore(appCtx.GetMainDBConnection())
		biz := rstlikebiz.NewListUserLikeRestaurantBiz(store)

		result, err := biz.ListUsers(c.Request.Context(), &filter, &paging)

		if err != nil {
			panic(err)
		}

		for i := range result {
			result[i].Mask(false)
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(result, paging, filter))
	}
}
