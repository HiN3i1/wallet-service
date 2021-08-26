package controller

import (
	"github.com/HiN3i1/wallet-service/types/api"
	"github.com/gin-gonic/gin"
)

func HandleNotFound(c *gin.Context) {
	err := api.NotFound
	c.JSON(err.StatusCode, err)
	return
}
