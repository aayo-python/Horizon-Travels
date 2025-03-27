package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/janto-pee/Horizon-Travels.git/util"
)

type Server struct {
	config util.Config
	router *gin.Engine
}

func NewServer(config util.Config) (*Server, error) {
	server := &Server{
		config: config,
	}
	server.setUpRouter()
	return server, nil
}

func (server *Server) setUpRouter() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})
	router.GET("/api/hotels", ListHotels)
	router.GET("/api/hotels/:id", GetHotelByID)
	router.GET("/api/hotels/filter", FilterHotels)
	router.GET("/api/hotels/search", SearchHotels)
	router.GET("/api/hotels/location", SearchHotels)
	/*
	*	MUTATIONS
	 */
	router.POST("/api/hotels", CreateHotels)
	router.PUT("/api/hotels/:id", updateHotel)
	router.DELETE("/api/hotels/:id", deleteHotel)
	router.POST("/movies/aggregations", AggregateHotels)

	server.router = router

	// r.Run()
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
