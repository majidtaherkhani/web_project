package main

// length 2 error for any problem with keys
// including having any extra key or one missing key
import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
)

var jwtKey = []byte("majidT&matinF")
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// var hexRegex = regexp.MustCompile("[0-9a-fA-F]+")
var clientOptions = options.Client().ApplyURI("mongodb://localhost:27017")
var client, err = mongo.Connect(context.TODO(), clientOptions)

func main() {

	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.POST("/api/signup", signUp)
	router.POST("/api/signin", signIn)
	router.POST("/api/edit", editProfile)
	router.GET("/api/:username/profile", getProfile)
	router.POST("/api/createPost", createPost)
	router.POST("api/like/:postId/:username", likePost)
	router.GET("/api/getPost/:postId", getPost)

	router.Run(":8080")
}

func validateToken(c *gin.Context) (userEmail string, ok bool) {
	tknString, err := c.Cookie("token")
	if err != nil {
		return "not_set", false
	}
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !tkn.Valid {
		return "not_set", false
	}
	return claims.Email, true
}

// checks if form has keys matching the keys array.
// also checks if form length is equal to keys length.
func isFormValid(form url.Values, keys []string) bool {
	if len(form) != len(keys) {
		return false
	}
	for _, val := range keys {
		_, ok1 := form[val]
		if !ok1 {
			return false
		}
	}
	return true
}

// valiadates email
func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	if !emailRegex.MatchString(e) {
		return false
	}
	parts := strings.Split(e, "@")
	mx, err := net.LookupMX(parts[1])
	if err != nil || len(mx) == 0 {
		return false
	}
	return true
}

// # length error
func validateSignUpRequest(c *gin.Context) bool {
	c.Request.ParseMultipartForm(1000)
	myform := c.Request.PostForm
	if !isFormValid(myform, []string{"email", "username", "password"}) {
		c.JSON(400, gin.H{
			"message": "Request Length should be 3",
		})
		return false
	}
	email := myform["email"][0]
	password := myform["password"][0]
	if !isEmailValid(email) {
		c.JSON(400, gin.H{
			"message": "filed `email` is not valid",
		})
		return false
	} else if len(password) <= 5 {
		c.JSON(400, gin.H{
			"message": "field 'password'.length should be gt 5.",
		})
		return false
	}
	return true
}
func validateSigninRequest(c *gin.Context) bool {
	c.Request.ParseMultipartForm(1000)
	myform := c.Request.PostForm
	if !isFormValid(myform, []string{"username", "password"}) {
		c.JSON(400, gin.H{
			"message": "Request Length should be 2",
		})
		return false
	}
	password := myform["password"][0]
	if len(password) <= 5 {
		c.JSON(400, gin.H{
			"message": "field 'password'.length should be gt 5.",
		})
		return false
	}
	return true
}

func validateEditRequest(c *gin.Context) bool {
	c.Request.ParseMultipartForm(1000)
	myform := c.Request.PostForm
	if !isFormValid(myform, []string{"username", "bio", "email", "password"}) {
		c.JSON(400, gin.H{
			"message": "Request Length should be 4",
		})
		return false
	}
	return true
}

func makeToken(email string) (string, bool) {
	expirationTime := time.Now().Add(60 * time.Hour)
	// build token
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", true
	}
	return tokenString, false
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func generateRandom() string {
	length := 10
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func signUp(c *gin.Context) {
	if !validateSignUpRequest(c) {
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
		collection.InsertOne(context.TODO(), User{Email: email, Username: username, Password: password, Created_at: time.Now().Format("01-02-2006")})
		c.JSON(201, gin.H{
			"massage": "user has been created",
		})
	}
}

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

func editProfile(c *gin.Context) {
	userEmail, ok := validateToken(c)
	if !ok {
		c.JSON(401, gin.H{
			"message": "permission denied.",
		})
		return
	} else if !validateEditRequest(c) {
		return
	}
	fmt.Println(userEmail)
	myform := c.Request.PostForm
	username := myform["username"][0]
	bio := myform["bio"][0]
	email := myform["email"][0]
	password := myform["password"][0]
	if !isEmailValid(email) {
		c.JSON(400, gin.H{
			"message": "filed `email` is not valid",
		})
		return
	}
	collection := client.Database("Web_HW3").Collection("User")

	filter := bson.M{
		"username": bson.M{
			"$eq": username, // check if bool field has value of 'false'
		},
	}
	update := bson.M{"$set": bson.M{
		"bio":      bio,
		"email":    email,
		"password": password,
	}}
	result, err := collection.UpdateOne(
		context.Background(),
		filter,
		update,
	)
	if err != nil {
		c.JSON(400, gin.H{
			"ine": result.UpsertedID,
		})
		return
	}
}

func signIn(c *gin.Context) {
	if !validateSigninRequest(c) {
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

	tokenString, err := makeToken(username)
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
	Bio        string `json:"bio,omitempty"`
	Created_at string `json:"created-at,omitempty"`
}

//////////////////post
type Post struct {
	Id         string   `bson:"_id" json:"id,omitempty"`
	Creator    string   `json:"creator,omitempty"`
	Content    string   `json:"content,omitempty"`
	Parent     string   `json:"parent,omitempty"`
	Likes      []string `json:"likes,omitempty"`
	Created_at string   `json:"created-at,omitempty"`
}

func likePost(c *gin.Context) {
	userEmail, ok := validateToken(c)
	if !ok {
		c.JSON(401, gin.H{
			"message": "permission denied.",
		})
		return
	}
	fmt.Println(userEmail)
	collection := client.Database("Web_HW3").Collection("Post")
	postId := c.Param("postId")
	username := c.Param("username")
	filterF := bson.M{"_id": postId}
	var post Post
	err = collection.FindOne(context.TODO(), filterF).Decode(&post)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "url id is not valid",
		})
		return
	}
	for i := range post.Likes {
		if post.Likes[i] == username {
			removeLike(c, i, post.Likes, postId)
			return
		}
	}
	saveLike(c, post.Likes, postId, username)
	return

}
func saveLike(c *gin.Context, likes []string, postId string, username string) {
	collection := client.Database("Web_HW3").Collection("Post")
	nLikes := append(likes, username)
	filterU := bson.M{
		"_id": bson.M{
			"$eq": postId,
		},
	}
	update := bson.M{"$set": bson.M{"likes": nLikes}}
	result, err := collection.UpdateOne(
		context.Background(),
		filterU,
		update,
	)
	if err != nil {
		c.JSON(400, gin.H{
			"ine": result.UpsertedID,
		})
		return
	}
}

func removeLike(c *gin.Context, index int, likes []string, postId string) {
	collection := client.Database("Web_HW3").Collection("Post")
	copy(likes[index:], likes[index+1:])
	likes[len(likes)-1] = ""
	likes = likes[:len(likes)-1]
	filterU := bson.M{
		"_id": bson.M{
			"$eq": postId,
		},
	}
	update := bson.M{"$set": bson.M{"likes": likes}}
	result, err := collection.UpdateOne(
		context.Background(),
		filterU,
		update,
	)
	if err != nil {
		c.JSON(400, gin.H{
			"ine": result.UpsertedID,
		})
		return
	}
}

func createPost(c *gin.Context) {
	userEmail, ok := validateToken(c)
	if !ok {
		c.JSON(401, gin.H{
			"message": "permission denied.",
		})
		return
	} else if !validatePost(c) {
		return
	}
	myform := c.Request.PostForm
	// creator := myform["creator"][0]
	content := myform["content"][0]
	parent := myform["parent"][0]
	collection := client.Database("Web_HW3").Collection("Post")
	id := generateRandom()
	insertResult, _ := collection.InsertOne(context.TODO(), Post{Id: id, Creator: userEmail, Content: content, Parent: parent,
		Created_at: time.Now().Format("01-02-2006")})
	c.JSON(201, gin.H{
		"id": insertResult.InsertedID,
	})

}
func validatePost(c *gin.Context) bool {
	c.Request.ParseMultipartForm(1000)
	form := c.Request.PostForm
	if !isFormValid(form, []string{"content", "creator", "parent"}) {
		c.JSON(400, gin.H{
			"message": "Request Length should be 3",
		})
		return false
	}
	creator := form["creator"][0]
	content := form["content"][0]
	if creator == "" {
		c.JSON(400, gin.H{
			"message": "filed `creator` is not valid",
		})
		return false
	} else if content == "" {
		c.JSON(400, gin.H{
			"message": "filed `content` is not valid",
		})
		return false
	}
	return true
}

func getPost(c *gin.Context) {
	userEmail, ok := validateToken(c)
	if !ok {
		c.JSON(401, gin.H{
			"message": "permission denied.",
		})
		return
	}
	fmt.Println(userEmail)
	postId := c.Param("postId")
	return findPostById(c, postId)

}

func findPostById(c *gin.Context, postId string) {
	collection := client.Database("Web_HW3").Collection("Post")
	filter := bson.M{"_id": postId}
	var post Post
	err = collection.FindOne(context.TODO(), filter).Decode(&post)
	if err == nil {
		c.JSON(409, gin.H{
			"message": "username already exist.",
		})
		return
	} else {

	}
}
