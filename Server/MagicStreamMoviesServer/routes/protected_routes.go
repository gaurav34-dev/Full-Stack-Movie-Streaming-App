package routes

import (
	controller "github.com/gaurav34-dev/MagicStreamMovies/Server/MagicStreamMoviesServer/controllers"
	"github.com/gaurav34-dev/MagicStreamMovies/Server/MagicStreamMoviesServer/middleware"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)


func SetupProtectedRoutes(router *gin.Engine, client *mongo.Client) {
	router.Use(middleware.AuthMiddleWare())

	router.GET("/movie/:imdb_id", controller.GetMovieByID(client))
	router.POST("/addmovie", controller.AddMovie(client))
	router.GET("/recommendedmovies", controller.GetRecommendedMovies(client))
	router.PATCH("/updatereview/:imdb_id", controller.AdminReviewUpdate(client))
}

