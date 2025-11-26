package busController

import (
	"wejh-go/app/apiException"
	"wejh-go/app/services/yxyServices"
	"wejh-go/app/utils"

	"github.com/gin-gonic/gin"
)

type getBusInfoRequest struct {
	Page     string `form:"page"`
	PageSize string `form:"page_size"`
	Search   string `form:"search"`
}

func GetBusInfo(c *gin.Context) {
	var req getBusInfoRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	busInfo, err := yxyServices.GetBusInfo(req.Page, req.PageSize, req.Search)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	utils.JsonSuccessResponse(c, busInfo)
}
