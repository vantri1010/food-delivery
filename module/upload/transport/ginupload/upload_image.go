package ginupload

import (
	"fmt"
	"food-delivery/common"
	"food-delivery/component/appctx"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpLoadImage(appCtx appctx.AppContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		fileHeader, err := c.FormFile("file")

		if err != nil {
			panic(err)
		}

		if err := c.SaveUploadedFile(fileHeader, fmt.Sprintf("static/%s", fileHeader.Filename)); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(common.Image{
			Id:        0,
			Url:       "http://localhost:8080/static/" + fileHeader.Filename,
			Width:     80,
			Height:    80,
			CloudName: "local",
			Extension: "png",
		}))
	}
}
