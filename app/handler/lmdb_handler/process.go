package lmdb_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Process(c *gin.Context) {
	name := c.Param("path")
	c.String(http.StatusOK, "LMDB Path: %s", name)
}
