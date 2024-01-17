package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaoz194/GoFlux/src/internal/example/gin_server/beans"
	"github.com/xiaoz194/GoFlux/src/internal/example/gin_server/serializer"
	"github.com/xiaoz194/GoFlux/src/pkg/e/constant"
	"github.com/xiaoz194/GoFlux/src/pkg/utils/logutil"
	"net/http"
)

func TestGet(c *gin.Context) {
	name := c.Query("name")
	logutil.LogrusObj.Info("name: ", name)
	c.JSON(http.StatusOK, &serializer.Response{
		Status: http.StatusOK,
		Msg:    constant.Ok,
		Data:   name,
	})
}

func TestPost(c *gin.Context) {
	var requestData beans.TestPostRequestData
	err := c.ShouldBind(&requestData)
	if err != nil {
		logutil.LogrusObj.Errorf("bind data err: %v", err)
		c.JSON(http.StatusBadRequest, serializer.Response{
			Status: http.StatusBadRequest,
			Msg:    "bind data err",
		})
		return
	}
	logutil.LogrusObj.Info("---------------------- requestData:  ", requestData)
	c.JSON(http.StatusOK, &serializer.Response{
		Status: http.StatusOK,
		Msg:    constant.Ok,
		Data:   requestData,
	})
}
