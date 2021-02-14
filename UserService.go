package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func signUp(c *gin.Context) {
	if !validateSigninSignUpRequest(c) {
		return
	}
	myform := c.Request.PostForm
	email := myform["email"][0]
	username := myform["username"][0]
	password := myform["password"][0]

	collection := client.Database("Web_HW3").Collection("User")
	filter := bson.M{"username": username}
	var result User
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err == nil {
		c.JSON(409, gin.H{
			"message": "username already exist.",
		})
		return
	} else {
		collection.InsertOne(context.TODO(), User{email, password, time.Now().Format("01-02-2006")})
		c.JSON(201, gin.H{
			"massage": "user has been created",
		})
	}
}

func editProfile(c *gin.Context) {
	if !validateSigninSignUpRequest(c) {
		return
	}

	myform := c.Request.PostForm
	email := myform["email"][0]
	username := myform["username"][0]
	password := myform["password"][0]
	bio := myform["bio"][0]

	collection := client.Database("Web_HW3").Collection("User")

	filter := bson.M{"username": bson.M{"$eq": username}}
	update := bson.M{"$set": bson.M{"bio": bio}}
	result, err := collection.UpdateOne(
		context.Background(),
		filter,
		update,
	)
	if err == nil {
		c.JSON(409, gin.H{
			"message": "username dosn't exist.",
		})
		return
	} else {
		c.JSON(201, gin.H{
			"username": "shod",
		})
	}
}

func signIn(c *gin.Context) {
	if !validateSigninSignUpRequest(c) {
		return
	}
	myform := c.Request.PostForm
	username := myform["username"][0]
	password := myform["password"][0]

	collection := client.Database("Web_HW3").Collection("User")
	filter := bson.M{"username": username}
	var result User
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil || result.Password != password {
		c.JSON(400, gin.H{
			"message": "wrong username or password.",
		})
		return
	}

	tokenString, err := makeToken(email)
	if err {
		fmt.Println("token string error")
		return
	}

	c.SetCookie("token", tokenString, 3600, "/", "localhost", false, true)
	c.JSON(200, gin.H{
		"message": tokenString,
	})
}

type User struct {
	Email      string `json:"email,omitempty"`
	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	Created_at string `json:"created-at,omitempty"`
	Bio        string `json:"bio,omitempty"`
}
