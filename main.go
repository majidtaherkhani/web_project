package main

// length 2 error for any problem with keys
// including having any extra key or one missing key
import (
	"context"
	"log"
	"math/rand"
	"net"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
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
	// router.POST("/api/signup", signUp)
	// router.POST("/api/signin", signIn)
	// router.GET("/api/post/", getPosts)
	// router.POST("/api/admin/post/crud", createPost)
	// router.PUT("/api/admin/post/crud/:id", updatePost)
	// router.DELETE("/api/admin/post/delet/:id", deletePost)
	// router.GET("/api/admin/post/crud/*id", readPost)
	// router.GET("/api/admin/user/crud/:id", getUser)
	// router.Run(":8081")

	router.POST("/api/signup", signUp)
	router.POST("/api/signin", signIn)
	// router.GET("/api/user/posts", getPosts)
	// router.GET("/api/user/profile", getProfile)
	// router.GET("/api/post", getPost)
	// router.POST("/api/send/comment", sendComment)
	router.POST("/api/edit/profile", editProfile)
	// router.POST("/api/like", saveLike)
	// router.GET("/api/timelie", getTimeLine)
	router.Run(":8080")
}

// func signUp(c *gin.Context) {
// 	if !validateSigninSignUpRequest(c) {
// 		return
// 	}
// 	myform := c.Request.PostForm
// 	email := myform["email"][0]
// 	password := myform["password"][0]

	collection := client.Database("Web_HW3").Collection("User")
// 	filter := bson.M{"email": email}
// 	var result User
// 	err = collection.FindOne(context.TODO(), filter).Decode(&result)
// 	if err == nil {
// 		c.JSON(409, gin.H{
// 			"message": "email already exist.",
// 		})
// 		return
// 	} else {
// 		collection.InsertOne(context.TODO(), User{email, password, time.Now().Format("01-02-2006")})
// 		c.JSON(201, gin.H{
// 			"massage": "user has been created",
// 		})
// 	}
// }

// func signIn(c *gin.Context) {
// 	if !validateSigninSignUpRequest(c) {
// 		return
// 	}
// 	myform := c.Request.PostForm
// 	email := myform["email"][0]
// 	password := myform["password"][0]

// 	collection := client.Database("Web_HW3").Collection("User")
// 	filter := bson.M{"email": email}
// 	var result User
// 	err = collection.FindOne(context.TODO(), filter).Decode(&result)
// 	if err != nil || result.Password != password {
// 		c.JSON(400, gin.H{
// 			"message": "wrong email or password.",
// 		})
// 		return
// 	}

// 	tokenString, err := makeToken(email)
// 	if err {
// 		fmt.Println("token string error")
// 		return
// 	}

// 	c.SetCookie("token", tokenString, 3600, "/", "localhost", false, true)
// 	c.JSON(200, gin.H{
// 		"message": tokenString,
// 	})
// }

// // returns 401 if UNAUTHORIZED ( token has been expired or no token )
// func createPost(c *gin.Context) {
// 	fmt.Println(StringWithCharset())
// 	userEmail, ok := validateToken(c)
// 	if !ok {
// 		c.JSON(401, gin.H{
// 			"message": "permission denied..",
// 		})
// 		return
// 	} else if !validatePost(c) {
// 		return
// 	}
// 	myform := c.Request.PostForm
// 	title := myform["title"][0]
// 	content := myform["content"][0]
// 	collection := client.Database("Web_HW3").Collection("Post")
// 	insertResult, _ := collection.InsertOne(context.TODO(), Post{Title: title, Content: content, Created_by: userEmail,
// 		Created_at: time.Now().Format("01-02-2006")})
// 	c.JSON(200, gin.H{
// 		"id": insertResult.InsertedID,
// 	})

// }

// // # last line: c.String(204, "") ??
// func updatePost(c *gin.Context) {
// 	userEmail, ok := validateToken(c)
// 	if !ok {
// 		c.JSON(401, gin.H{
// 			"message": "permission denied..",
// 		})
// 		return
// 	} else if !validatePost(c) {
// 		return
// 	}
// 	myform := c.Request.PostForm
// 	title := myform["title"][0]
// 	content := myform["content"][0]
// 	postId := c.Param("id")
// 	postObjectId, err := primitive.ObjectIDFromHex(postId)
// 	if err != nil {
// 		c.JSON(400, gin.H{
// 			"message": "url id is not valid",
// 		})
// 		return
// 	}
// 	collection := client.Database("Web_HW3").Collection("Post")
// 	filter := bson.M{"_id": postObjectId}
// 	var targetPost Post
// 	err = collection.FindOne(context.TODO(), filter).Decode(&targetPost)
// 	if err != nil {
// 		c.JSON(400, gin.H{
// 			"message": "url id is not valid",
// 		})
// 		return
// 	}
// 	if targetPost.Created_by != userEmail {
// 		c.JSON(401, gin.H{
// 			"message": "permission denied.",
// 		})
// 		return
// 	}
// 	collection.UpdateOne(context.TODO(), filter,
// 		bson.D{
// 			{"$set", bson.D{{"title", title}, {"content", content}}},
// 		})
// 	c.String(204, "")
// }

// func deletePost(c *gin.Context) {
// 	userEmail, ok := validateToken(c)
// 	if !ok {
// 		c.JSON(401, gin.H{
// 			"message": "you are UNAUTHORIZED.",
// 		})
// 		return
// 	}
// 	postId := c.Param("id")
// 	postObjectId, err := primitive.ObjectIDFromHex(postId)
// 	if err != nil {
// 		c.JSON(400, gin.H{
// 			"message": "url id is not valid",
// 		})
// 		return
// 	}
// 	collection := client.Database("Web_HW3").Collection("Post")
// 	filter := bson.M{"_id": postObjectId}
// 	var targetPost Post
// 	err = collection.FindOne(context.TODO(), filter).Decode(&targetPost)
// 	if err != nil {
// 		c.JSON(400, gin.H{
// 			"message": "url id is not valid",
// 		})
// 		return
// 	}
// 	if targetPost.Created_by != userEmail {
// 		c.JSON(401, gin.H{
// 			"message": "permission denied.",
// 		})
// 		return
// 	}
// 	collection.DeleteMany(context.TODO(), filter)
// 	c.String(204, "")
// }

// func readPost(c *gin.Context) {
// 	userEmail, ok := validateToken(c)
// 	if !ok {
// 		c.JSON(401, gin.H{
// 			"message": "permission denied.",
// 		})
// 		return
// 	}
// 	postId := c.Param("id")[1:]
// 	postsCollection := client.Database("Web_HW3").Collection("Post")
// 	if postId == "" {
// 		sendUserPosts(c, userEmail)
// 		return
// 	}

// 	postObjectId, err := primitive.ObjectIDFromHex(postId)
// 	if err != nil {
// 		c.JSON(400, gin.H{
// 			"message": "url id is not valid",
// 		})
// 		return
// 	}

// 	filter := bson.M{"_id": postObjectId}
// 	var targetPost PostWithId
// 	err = postsCollection.FindOne(context.TODO(), filter).Decode(&targetPost)
// 	if err != nil {
// 		c.JSON(400, gin.H{
// 			"message": "url id is not valid",
// 		})
// 		return
// 	}
// 	c.JSON(200, bson.M{"posts": targetPost})

// }

// func getPosts(c *gin.Context) {
// 	postsCollection := client.Database("Web_HW3").Collection("Post")
// 	var result []PostWithId
// 	cur, _ := postsCollection.Find(context.TODO(), bson.M{})
// 	for cur.Next(context.TODO()) {
// 		var elem PostWithId
// 		err := cur.Decode(&elem)
// 		// fmt.Println(elem)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		result = append(result, elem)
// 	}

// 	c.JSON(200, bson.M{"posts": result})
// }

// func getUser(c *gin.Context) {
// 	userEmail, ok := validateToken(c)
// 	if !ok {
// 		c.JSON(401, gin.H{
// 			"message": "permission denied.",
// 		})
// 		return
// 	}
// 	userId := c.Param("id")
// 	userObjectId, err := primitive.ObjectIDFromHex(userId)
// 	if err != nil {
// 		c.JSON(400, gin.H{
// 			"message": "url id is not valid",
// 		})
// 		return
// 	}
// 	collection := client.Database("Web_HW3").Collection("User")
// 	filter := bson.M{"_id": userObjectId}
// 	var targetUser UserWithId
// 	err = collection.FindOne(context.TODO(), filter).Decode(&targetUser)
// 	if err != nil {
// 		c.JSON(400, gin.H{
// 			"message": "url id is not valid",
// 		})
// 		return
// 	}
// 	fmt.Println(targetUser.Email)
// 	fmt.Println(userEmail)

// 	if targetUser.Email != userEmail {
// 		c.JSON(401, gin.H{
// 			"message": "permission denied.",
// 		})
// 		return
// 	}

// 	c.JSON(200, bson.M{"user": targetUser})
// }

// func sendUserPosts(c *gin.Context, userEmail string) {
// 	postsCollection := client.Database("Web_HW3").Collection("Post")
// 	var result []PostWithId
// 	filter := bson.M{"created_by": userEmail}
// 	cur, err := postsCollection.Find(context.TODO(), filter)

// 	if err != nil {
// 		c.JSON(400, gin.H{
// 			"message": "url id is not valid",
// 		})
// 		return
// 	}
// 	for cur.Next(context.TODO()) {
// 		var elem PostWithId
// 		err := cur.Decode(&elem)
// 		// fmt.Println(elem)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		result = append(result, elem)
// 	}

// 	c.JSON(200, bson.M{"posts": result})
// }

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

// validates post to have just title and content keys
// generates the error itself
func validatePost(c *gin.Context) bool {
	c.Request.ParseMultipartForm(1000)
	form := c.Request.PostForm
	if !isFormValid(form, []string{"title", "content"}) {
		c.JSON(400, gin.H{
			"message": "Request Length should be 2",
		})
		return false
	}
	title := form["title"][0]
	content := form["content"][0]
	if title == "" {
		c.JSON(400, gin.H{
			"message": "filed `title` is not valid",
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
func validateSigninSignUpRequest(c *gin.Context) bool {
	c.Request.ParseMultipartForm(1000)
	myform := c.Request.PostForm
	if !isFormValid(myform, []string{"email", "password"}) {
		c.JSON(400, gin.H{
			"message": "Request Length should be 2",
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

func makeToken(email string) (string, bool) {
	expirationTime := time.Now().Add(60 * time.Minute)
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

type User struct {
	Email      string `json:"email,omitempty"`
	Password   string `json:"password,omitempty"`
	Created_at string `json:"created-at,omitempty"`
}

type UserWithId struct {
	Id         string `bson:"_id" json:"id,omitempty"`
	Email      string `json:"email,omitempty"`
	Password   string `json:"password,omitempty"`
	Created_at string `json:"created-at,omitempty"`
}

type Post struct {
	Title      string `json:"title,omitempty"`
	Content    string `json:"content,omitempty"`
	Created_by string `json:"created_by,omitempty"`
	Created_at string `json:"created-at,omitempty"`
}

type PostWithId struct {
	Id         string `bson:"_id" json:"id,omitempty"`
	Title      string `json:"title,omitempty"`
	Content    string `json:"content,omitempty"`
	Created_by string `json:"created_by,omitempty"`
	Created_at string `json:"created-at,omitempty"`
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset() string {
	length := 10
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
