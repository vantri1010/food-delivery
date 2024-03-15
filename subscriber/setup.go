package subscriber

import (
	"context"
	"food-delivery/component/appctx"
)

func Setup(appCtx appctx.AppContext, ctx context.Context) {
	InCreaseLikeCountAfterUserLikeRestaurant(appCtx, ctx)
	DeCreaseLikeCountAfterUserDisLikeRestaurant(appCtx, ctx)
	PushNotiUserLikeRestaurant(appCtx, ctx)
}
