package main

import (
	"github.com/gin-gonic/gin"
	"music-final-project/configuration"
	"music-final-project/controllers/albumcontroller"
	"music-final-project/controllers/artistcontroller"
	"music-final-project/controllers/songcontroller"
)

func main() {
	route := gin.Default()
	configuration.DatabaseConnect()

	// artist route
	route.GET("/api/artists", artistcontroller.Index)
	route.GET("/api/artist/:id", artistcontroller.Show)
	route.POST("/api/artist", artistcontroller.Create)
	route.PUT("/api/artist/:id", artistcontroller.Update)
	route.DELETE("/api/artist/:id", artistcontroller.Delete)

	// album route
	route.GET("/api/albums", albumcontroller.Index)
	route.GET("/api/album/:id", albumcontroller.Show)
	route.POST("/api/album", albumcontroller.Create)
	route.PUT("/api/album/:id", albumcontroller.Update)
	route.DELETE("/api/album/:id", albumcontroller.Delete)

	// song route
	route.GET("/api/songs", songcontroller.Index)
	route.GET("/api/song/:id", songcontroller.Show)
	route.POST("/api/song", songcontroller.Create)
	route.PUT("/api/song/:id", songcontroller.Update)
	route.DELETE("/api/song/:id", songcontroller.Delete)

	route.Run()
}
