package busController

import (
	"wejh-go/app/apiException"
	"wejh-go/app/services/sessionServices"
	"wejh-go/app/services/yxyServices"
	"wejh-go/app/utils"

	"github.com/gin-gonic/gin"
)

type infoForm struct {
	Page     string `form:"page" json:"page"`
	PageSize string `form:"page_size" json:"page_size"`
	Search   string `form:"search,optional" json:"search"`
}

type recordsForm struct {
	Page     string `form:"page" json:"page"`
	PageSize string `form:"page_size" json:"page_size"`
	Status   string `form:"status" json:"status"`
}

type messageForm struct {
	Page     string `form:"page" json:"page"`
	PageSize string `form:"page_size" json:"page_size"`
}

func GetBusInfo(c *gin.Context) {
	var postForm infoForm
	err := c.ShouldBindQuery(&postForm)
	if err != nil {
		_ = c.AbortWithError(200, apiException.ParamError)
		return
	}

	busInfo, err := yxyServices.BusInfo(postForm.Page, postForm.PageSize, postForm.Search)
	if err != nil {
		_ = c.AbortWithError(200, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, busInfo)
}

func GetBusRecords(c *gin.Context) {
	var postForm recordsForm
	err := c.ShouldBindQuery(&postForm)
	if err != nil {
		_ = c.AbortWithError(200, apiException.ParamError)
		return
	}

	user, err := sessionServices.GetUserSession(c)
	if err != nil {
		_ = c.AbortWithError(200, apiException.NotLogin)
		return
	}

	token, err := yxyServices.GetBusAuthToken(user.YxyUid)
	if err != nil {
		_ = c.AbortWithError(200, apiException.ServerError)
		return
	}

	busRecords, err := yxyServices.BusRecords(*token, postForm.Page, postForm.PageSize, postForm.Status)
	if err != nil {
		_ = c.AbortWithError(200, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, busRecords)
}

func GetBusQrcode(c *gin.Context) {
	user, err := sessionServices.GetUserSession(c)
	if err != nil {
		_ = c.AbortWithError(200, apiException.NotLogin)
		return
	}

	token, err := yxyServices.GetBusAuthToken(user.YxyUid)
	if err != nil {
		_ = c.AbortWithError(200, apiException.ServerError)
		return
	}

	qrcode, err := yxyServices.BusQrcode(*token)
	if err != nil {
		_ = c.AbortWithError(200, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, qrcode)
}

func GetBusMessage(c *gin.Context) {
	var postForm messageForm
	err := c.ShouldBindQuery(&postForm)
	if err != nil {
		_ = c.AbortWithError(200, apiException.ParamError)
		return
	}

	user, err := sessionServices.GetUserSession(c)
	if err != nil {
		_ = c.AbortWithError(200, apiException.NotLogin)
		return
	}

	token, err := yxyServices.GetBusAuthToken(user.YxyUid)
	if err != nil {
		_ = c.AbortWithError(200, apiException.ServerError)
		return
	}

	message, err := yxyServices.BusMessage(*token, postForm.Page, postForm.PageSize)
	if err != nil {
		_ = c.AbortWithError(200, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, message)
}
