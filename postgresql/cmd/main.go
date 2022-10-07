package main

import (
	"fmt"
	"net/http"

	"github.com/bzhtux/sample_apps/postgresql/models"
	pg "github.com/bzhtux/sample_apps/postgresql/pkg/postgresql"
	"github.com/gin-gonic/gin"
)

var (
	version = "v0.0.1"
)

func main() {
	fmt.Printf("\033[32mLauching sample_app-postgres %s...\033[0m\n", version)
	pgConn := pg.OpenDB()
	pgh := pg.New(pgConn)
	pgConn.AutoMigrate(&models.Books{})

	// gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.MaxMultipartMemory = 16 << 32 // 16 MiB

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "Up",
			"message": "Alive",
		})
	})
	router.POST("/add", pgh.AddNewBook)
	router.GET("/get/:uid", pgh.GetBookByID)
	router.DELETE("/del/:uid", pgh.DeleteBook)
	router.Run(":8080")
}
