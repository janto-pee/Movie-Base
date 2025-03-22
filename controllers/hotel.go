package controllers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/janto-pee/Horizon-Travels.git/model"
	"github.com/janto-pee/Horizon-Travels.git/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Implemention of the /movies route that returns all of the movies from our movies collection.
func GetMovies(c *gin.Context) {
	// Find movies
	cursor, err := util.Db.Find(context.TODO(), bson.D{{}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Map results
	var movies []bson.M
	if err = cursor.All(context.TODO(), &movies); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return movies
	c.JSON(http.StatusOK, movies)
}

// The implementation of our /movies/{id} endpoint that returns a single movie based on the provided ID
func GetMovieByID(c *gin.Context) {

	// Get movie ID from URL
	idStr := c.Param("id")

	// Convert id string to ObjectId
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find movie by ObjectId
	var movie bson.M
	err = util.Db.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return movie
	c.JSON(http.StatusOK, movie)
}

// The implementation of our /movies/aggregations endpoint that allows a user to pass in an aggregation to run our the movies collection.
func AggregateMovies(c *gin.Context) {
	// Get aggregation pipeline from request body
	var pipeline interface{}
	if err := c.ShouldBindJSON(&pipeline); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Run aggregations
	cursor, err := util.MongoClient.Database("sample_mflix").Collection("movies").Aggregate(context.TODO(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Map results
	var result []bson.M
	if err = cursor.All(context.TODO(), &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return result
	c.JSON(http.StatusOK, result)
}

type createHotelRequest struct {
	Title         string
	Content       string
	PrimaryInfo   string
	SecondaryInfo string
	AccentedLabel string
	Provider      string
	PriceDetails  string
	PriceSummary  string
}

func CreateMovies(c *gin.Context) {
	// Get aggregation pipeline from request body
	var req createHotelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := model.Hotel{
		Title:         req.Title,
		Content:       req.Content,
		PrimaryInfo:   req.PrimaryInfo,
		SecondaryInfo: req.SecondaryInfo,
		AccentedLabel: req.AccentedLabel,
		Provider:      req.Provider,
		PriceDetails:  req.PriceDetails,
		PriceSummary:  req.PriceSummary,
	}

	hotel, err := util.Db.InsertOne(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, hotel)

}

type updateHotelParam struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func updateHotel(c *gin.Context) {
	var id updateHotelParam
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req createHotelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := model.Hotel{
		Title:         req.Title,
		Content:       req.Content,
		PrimaryInfo:   req.PrimaryInfo,
		SecondaryInfo: req.SecondaryInfo,
		AccentedLabel: req.AccentedLabel,
		Provider:      req.Provider,
		PriceDetails:  req.PriceDetails,
		PriceSummary:  req.PriceSummary,
	}
	if err := util.Db.FindOne(c, bson.M{"id": id}).Decode(&arg); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if err := c.BindJSON(&arg); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	hotel, err := util.Db.UpdateByID(c, bson.M{"id": id}, bson.M{"$set": arg})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, hotel)
}

func deleteHotel(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	_, err := util.Db.DeleteOne(c, bson.M{"id": id})
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
}
