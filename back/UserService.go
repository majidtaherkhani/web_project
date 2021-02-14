package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func getProfile(c *gin.Context) {
	userEmail, ok := validateToken(c)
	if !ok {
		c.JSON(401, gin.H{
			"message": "permission denied.",
		})
		return
	}
	collection := client.Database("Web_HW3").Collection("User")
	username := c.Param("username")
	fmt.Println(username)
	fmt.Println(userEmail)
	filter := bson.M{"username": username}
	var profile User
	err = collection.FindOne(context.TODO(), filter).Decode(&profile)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "username is not valid",
		})
		return
	}
	c.JSON(200, bson.M{"bio": profile.Bio,
		"email": profile.Email})
}
