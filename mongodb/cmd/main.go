package main

import (
	"context"
	"log"
	"net/http"

	"github.com/bzhtux/sample_apps/mongodb/pkg/mongodb"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	version = "v0.0.2"
)

func main() {
	log.Printf("\033[32m*** Lauching sample_app-mongo %s...***\033[0m\n", version)
	clt, err := mongodb.NewClient()
	if err != nil {
		log.Printf("Error Getting new MongoDB client: %s\n", err.Error())
	}
	if err := clt.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Printf("\033[31m>>> Can not ping MongoDB innstance <<< \033[0m\n>>> %s <<<\n", err.Error())
	} else {
		log.Printf("\033[32m* PING MongoDB instance is OK *\033[0m\n")
	}
	mh := mongodb.New(clt)

	gin.SetMode(gin.ReleaseMode)
	// gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.MaxMultipartMemory = 16 << 32 // 16 MiB

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "Up",
			"message": "Alive",
		})
	})
	router.GET("/ping", func(c *gin.Context) {
		if err := clt.Ping(context.TODO(), readpref.Primary()); err != nil {
			log.Printf("\033[31m>>> Can not ping MongoDB innstance <<< \033[0m\n>>> %s <<<\n", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "Internal Error",
				"message": "MongoDB does not respond to PING command",
			})
		} else {
			log.Printf("\033[32m* PING MongoDB instance is OK *\033[0m\n")
			c.JSON(http.StatusOK, gin.H{
				"status":  "Ok",
				"message": "MongoDB responded to PING :-)",
			})
		}
	})
	router.POST("/add", mh.AddNewDoc)
	router.GET("/get/byName/:book", mh.GetOneDocByTitle)
	router.GET("/get/byID/:uid", mh.GetOneDocByID)
	router.DELETE("/del/byName/:book", mh.DeleteOnDocByName)
	router.Run(":8080")
}
