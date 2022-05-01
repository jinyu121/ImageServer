package csv_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Process(c *gin.Context) {
	name := c.Param("path")
	c.String(http.StatusOK, "CSV File Path %s", name)
}
