package main

// length 2 error for any problem with keys
// including having any extra key or one missing key
import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"

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
	router.Use(cors.Default())
	router.LoadHTMLGlob("FrontEnd/*.html")
	router.StaticFile("/register.css", "./FrontEnd/register.css")
	router.StaticFile("/timeline.css", "./FrontEnd/timeline.css")
	router.StaticFile("/register.js", "./FrontEnd/register.js")
	router.StaticFile("/timeline.js", "./FrontEnd/timeline.js")
	router.StaticFile("/images/nani.png", "./FrontEnd/images/nani.png")
	router.StaticFile("/images/no-image.jpg", "./FrontEnd/images/no-image.jpg")
	router.StaticFile("/images/profile-photo.jpg", "./FrontEnd/images/profile-photo.jpg")

	router.GET("", serveHTML)

	router.POST("/api/signup", signUp)
	router.POST("/api/signin", signIn)
	router.POST("/api/edit", editProfile)
	router.GET("/api/profile", getMyProfile)
	router.GET("/api/otherProfile/:username", getOtherProfile)
	router.POST("/api/createPost", createPost)
	router.POST("/api/like/:postId", likePost)
	router.POST("/api/mark", markPost)
	router.GET("/api/getPost/:postId", getPost)
	router.GET("/api/getComments/:postId", getComent)
	router.POST("/api/follow/:username", followUnfollow)
	router.GET("/api/timeline", getTimeline)
	router.GET("/api/followers", getFollowers)
	router.GET("/api/followings", getFollowings)
	router.GET("/api/bookMarks", getBookmarls)

	router.Run(":8080")
}

func serveHTML(c *gin.Context) {
	_, ok := validateToken(c)
	if !ok {
		c.HTML(http.StatusOK, "register.html", nil)
		return
	}
	c.HTML(http.StatusOK, "timeline.html", nil)
}

func validateToken(c *gin.Context) (username string, ok bool) {
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
	return claims.Username, true
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

func getUserById(username string) User {
	collection := client.Database("web_project").Collection("User")
	filter := bson.M{"username": username}
	var user User
	err = collection.FindOne(context.TODO(), filter).Decode(&user)
	return user
}

func validateEditRequest(c *gin.Context) bool {
	c.Request.ParseMultipartForm(1000)
	myform := c.Request.PostForm
	if !isFormValid(myform, []string{"bio", "fullName", "email", "password"}) {
		c.JSON(400, gin.H{
			"message": "Request Length should be 4",
		})
		return false
	}
	return true
}

func makeToken(username string) (string, bool) {
	expirationTime := time.Now().Add(6000 * time.Hour)
	// build token
	claims := &Claims{
		Username: username,
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
	Username string `json:"username"`
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

	collection := client.Database("web_project").Collection("User")
	filter := bson.M{"username": username}
	var result User
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err == nil {
		c.JSON(409, gin.H{
			"message": "username already exist.",
		})
		return
	}
	collection.InsertOne(context.TODO(), User{Email: email, Username: username, Password: password, Created_at: time.Now().Format("01-02-2006")})
	c.JSON(201, gin.H{
		"message": "user has been created",
	})
}

func getMyProfile(c *gin.Context) {
	username, ok := validateToken(c)
	if !ok {
		c.JSON(401, gin.H{
			"message": "permission denied.",
		})
		return
	}
	getProfile(c, username)
}
func getOtherProfile(c *gin.Context) {
	usernameM, ok := validateToken(c)
	if !ok {
		c.JSON(401, gin.H{
			"message": "permission denied.",
		})
		return
	}
	fmt.Println(usernameM)
	username := c.Param("username")
	getProfile(c, username)
}
func getProfile(c *gin.Context, username string) {

	collection := client.Database("web_project").Collection("User")
	filter := bson.M{"username": username}
	var profile User
	err = collection.FindOne(context.TODO(), filter).Decode(&profile)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "username is not valid",
		})
		return
	}
	collectionP := client.Database("web_project").Collection("Post")
	var posts []Post
	filterP := bson.M{
		"creator": bson.M{
			"$eq": username,
		},
	}
	curP, errP := collectionP.Find(context.TODO(), filterP)
	if errP != nil {
		c.JSON(409, gin.H{
			"message": "error in find posts",
		})
		return
	}
	for curP.Next(context.TODO()) {
		var post Post
		if err = curP.Decode(&post); err != nil {
			return
		}
		posts = append(posts, post)
	}
	c.JSON(200, bson.M{"bio": profile.Bio,
		"email":     profile.Email,
		"posts":     convertPosts(posts, profile),
		"following": len(profile.Followings),
		"followers": len(profile.Followers)})
}

func editProfile(c *gin.Context) {
	username, ok := validateToken(c)
	if !ok {
		c.JSON(401, gin.H{
			"message": "permission denied.",
		})
		return
	} else if !validateEditRequest(c) {
		return
	}
	fmt.Println(username)
	myform := c.Request.PostForm
	bio := myform["bio"][0]
	email := myform["email"][0]
	password := myform["password"][0]
	fullName := myform["fullName"][0]
	if !isEmailValid(email) {
		c.JSON(400, gin.H{
			"message": "filed `email` is not valid",
		})
		return
	}
	collection := client.Database("web_project").Collection("User")

	filter := bson.M{
		"username": bson.M{
			"$eq": username,
		},
	}
	update := bson.M{"$set": bson.M{
		"bio":      bio,
		"fullName": fullName,
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
			"message": "username is not valid",
		})
		return
	}
	c.JSON(400, gin.H{
		"profile": result,
	})
	return

}

func signIn(c *gin.Context) {
	if !validateSigninRequest(c) {
		return
	}
	myform := c.Request.PostForm
	username := myform["username"][0]
	password := myform["password"][0]

	collection := client.Database("web_project").Collection("User")
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

	c.SetCookie("token", tokenString, 3600, "/", "localhost", false, false)
	c.JSON(200, gin.H{
		"message": tokenString,
	})
}

type User struct {
	Email      string   `json:"email,omitempty"`
	Username   string   `json:"username,omitempty"`
	FullName   string   `json:"fullName,omitempty"`
	Password   string   `json:"password,omitempty"`
	Bio        string   `json:"bio,omitempty"`
	Followers  []string `json:"followers,omitempty"`
	Followings []string `json:"following,omitempty"`
	BookMark   []string `json:"bookMark,omitempty"`
	Created_at string   `json:"created-at,omitempty"`
}

//////////////////post
type Post struct {
	Id         string   `bson:"_id" json:"id,omitempty"`
	Creator    string   `json:"creator,omitempty"`
	Content    string   `json:"content,omitempty"`
	Parent     string   `json:"parent,omitempty"`
	Likes      []string `json:"likes,omitempty"`
	Created_at int64    `json:"created-at,omitempty"`
}

type PostForUser struct {
	Id           string   `bson:"_id" json:"id,omitempty"`
	Creator      string   `json:"creator"`
	FullName     string   `json:"fullName"`
	Content      string   `json:"content"`
	Parent       string   `json:"parent"`
	Likes        []string `json:"likes"`
	Like         bool     `json:"like"`
	Mark         bool     `json:"mark"`
	LikeNumber   int      `json:"likeNumber"`
	ComentNumber int      `json:"cumentNumber"`
	Created_at   int64    `json:"created-at"`
}

func markPost(c *gin.Context) {
	username, ok := validateToken(c)
	if !ok {
		c.JSON(401, gin.H{
			"message": "permission denied.",
		})
		return
	}
	collection := client.Database("web_project").Collection("User")
	postId := c.Param("postId")
	filterU := bson.M{
		"username": bson.M{
			"$eq": username,
		},
	}
	var user User
	err = collection.FindOne(context.TODO(), filterU).Decode(&user)
	if err != nil {
		c.JSON(409, gin.H{
			"message": "error",
		})
		return
	}
	for i := range user.BookMark {
		if user.BookMark[i] == postId {
			removeMark(c, i, user)
			return
		}
	}
	saveMark(c, user, postId)
	return
}
func removeMark(c *gin.Context, index int, user User) {
	collection := client.Database("web_project").Collection("User")
	marks := user.BookMark
	copy(marks[index:], marks[index+1:])
	marks[len(marks)-1] = ""
	marks = marks[:len(marks)-1]
	filter := bson.M{
		"username": bson.M{
			"$eq": user.Username,
		},
	}
	update := bson.M{"$set": bson.M{"bookMark": marks}}
	result, err := collection.UpdateOne(
		context.Background(),
		filter,
		update,
	)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "update failed",
		})
		return
	}
	fmt.Println(result.UpsertedID)
}

func saveMark(c *gin.Context, user User, postId string) {
	marks := append(user.BookMark, postId)
	collection := client.Database("web_project").Collection("Post")
	filter := bson.M{
		"username": bson.M{
			"$eq": user.Username,
		},
	}
	update := bson.M{"$set": bson.M{"bookMark": marks}}
	result, err := collection.UpdateOne(
		context.Background(),
		filter,
		update,
	)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "update falled",
		})
		return
	}
	fmt.Println(result.UpsertedID)

}

func likePost(c *gin.Context) {
	username, ok := validateToken(c)
	if !ok {
		c.JSON(401, gin.H{
			"message": "permission denied.",
		})
		return
	}
	collection := client.Database("web_project").Collection("Post")
	postId := c.Param("postId")
	filterF := bson.M{"_id": postId}
	var post Post
	err = collection.FindOne(context.TODO(), filterF).Decode(&post)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "post id is not valid",
		})
		return
	}
	for i := range post.Likes {
		if post.Likes[i] == username {
			removeLike(c, i, post.Likes, postId)
			c.JSON(200, gin.H{
				"message": "post disliked",
			})
			return
		}
	}
	saveLike(c, post.Likes, postId, username)
	return

}
func saveLike(c *gin.Context, likes []string, postId string, username string) {
	collection := client.Database("web_project").Collection("Post")
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
			"message": "update falled",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "post liked",
	})
	fmt.Println(result.UpsertedID)
}

func removeLike(c *gin.Context, index int, likes []string, postId string) {
	collection := client.Database("web_project").Collection("Post")
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
			"message": "update failed",
		})
		return
	}
	fmt.Println(result.UpsertedID)
}

func getPostById(postId string) Post {
	collection := client.Database("web_project").Collection("Post")
	filter := bson.M{"_id": postId}
	var post Post
	err = collection.FindOne(context.TODO(), filter).Decode(&post)
	return post
}

func createPost(c *gin.Context) {
	username, ok := validateToken(c)
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
	collection := client.Database("web_project").Collection("Post")
	id := generateRandom()
	insertResult, _ := collection.InsertOne(context.TODO(), Post{Id: id, Creator: username, Content: content, Parent: parent,
		Created_at: time.Now().Unix()})
	c.JSON(201, gin.H{
		"id": insertResult.InsertedID,
	})
}

func validatePost(c *gin.Context) bool {
	c.Request.ParseMultipartForm(1000)
	form := c.Request.PostForm
	if !isFormValid(form, []string{"content", "parent"}) {
		c.JSON(400, gin.H{
			"message": "Request Length should be 3",
		})
		return false
	}
	content := form["content"][0]
	if content == "" {
		c.JSON(400, gin.H{
			"message": "filed `content` is not valid",
		})
		return false
	}
	return true
}

func getPost(c *gin.Context) {
	username, ok := validateToken(c)
	if !ok {
		c.JSON(401, gin.H{
			"message": "permission denied.",
		})
		return
	}
	fmt.Println(username)
	postId := c.Param("postId")
	collection := client.Database("web_project").Collection("Post")
	filter := bson.M{"_id": postId}
	var post Post
	err = collection.FindOne(context.TODO(), filter).Decode(&post)
	if err != nil {
		c.JSON(409, gin.H{
			"message": "failed",
		})
		return
	}
	c.JSON(200, gin.H{
		"post": post,
	})
	return
}
func getComent(c *gin.Context) {
	username, ok := validateToken(c)
	if !ok {
		c.JSON(401, gin.H{
			"message": "permission denied.",
		})
		return
	}
	fmt.Println(username)
	collection := client.Database("web_project").Collection("Post")
	postId := c.Param("postId")
	commentFilter := bson.M{"parent": postId}
	var coments []Post
	cur, err := collection.Find(context.TODO(), commentFilter)
	if err != nil {
		c.JSON(409, gin.H{
			"message": "failed",
		})
		return
	}
	for cur.Next(context.TODO()) {
		var coment Post
		if err = cur.Decode(&coment); err != nil {
			return
		}
		coments = append(coments, coment)
	}
	user := getUserById(username)
	c.JSON(200, gin.H{
		"comment": convertPosts(coments, user),
	})
	return
}

///follower following
func followUnfollow(c *gin.Context) {
	username, ok := validateToken(c)
	if !ok {
		c.JSON(401, gin.H{
			"message": "permission denied.",
		})
		return
	}
	targetUsername := c.Param("username")
	collection := client.Database("web_project").Collection("User")
	filterMain := bson.M{"username": username}
	filterTarget := bson.M{"username": targetUsername}
	var mainUser User
	var targetUser User
	err = collection.FindOne(context.TODO(), filterMain).Decode(&mainUser)
	if err != nil {
		c.JSON(409, gin.H{
			"message": "failed",
		})
		return
	}
	err = collection.FindOne(context.TODO(), filterTarget).Decode(&targetUser)
	if err != nil {
		c.JSON(409, gin.H{
			"message": "failed",
		})
		return
	}
	for i := range mainUser.Followings {
		if mainUser.Followings[i] == targetUsername {
			for j := range targetUser.Followers {
				if targetUser.Followers[j] == username {
					unfollow(c, mainUser, targetUser, j, i)
					return
				}
			}
		}
	}
	follow(c, mainUser, targetUser)
	return
}

func unfollow(c *gin.Context, user User, target User, followerIndex int, followingIndex int) {
	collection := client.Database("web_project").Collection("User")
	followings := user.Followings
	copy(followings[followingIndex:], followings[followingIndex+1:])
	followings[len(followings)-1] = ""
	followings = followings[:len(followings)-1]
	filterM := bson.M{
		"username": bson.M{
			"$eq": user.Username,
		},
	}
	updateM := bson.M{"$set": bson.M{"followings": followings}}
	resultM, errM := collection.UpdateOne(
		context.Background(),
		filterM,
		updateM,
	)
	if errM != nil {
		c.JSON(400, gin.H{
			"message": "error",
		})
		return
	}
	fmt.Println(resultM.UpsertedID)

	followers := target.Followers
	copy(followers[followerIndex:], followers[followerIndex+1:])
	followers[len(followers)-1] = ""
	followers = followers[:len(followers)-1]
	filterT := bson.M{
		"username": bson.M{
			"$eq": target.Username,
		},
	}
	updateT := bson.M{"$set": bson.M{"followers": followers}}
	resultT, errT := collection.UpdateOne(
		context.Background(),
		filterT,
		updateT,
	)
	if errT != nil {
		c.JSON(400, gin.H{
			"message": "error",
		})
		return
	}
	fmt.Println(resultT.UpsertedID)

}

func follow(c *gin.Context, user User, target User) {
	collection := client.Database("web_project").Collection("User")
	followings := append(user.Followings, target.Username)
	filterM := bson.M{
		"username": bson.M{
			"$eq": user.Username,
		},
	}
	updateM := bson.M{"$set": bson.M{"followings": followings}}
	resultM, err := collection.UpdateOne(
		context.Background(),
		filterM,
		updateM,
	)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "error",
		})
		return
	}
	fmt.Println(resultM.UpsertedID)
	followers := append(target.Followers, user.Username)
	filterT := bson.M{
		"username": bson.M{
			"$eq": target.Username,
		},
	}
	updateT := bson.M{"$set": bson.M{"followers": followers}}
	resultT, errt := collection.UpdateOne(
		context.Background(),
		filterT,
		updateT,
	)
	if errt != nil {
		c.JSON(400, gin.H{
			"message": "error",
		})
		return
	}
	fmt.Println(resultT.UpsertedID)
}

//////timeline

func getTimeline(c *gin.Context) {
	username, ok := validateToken(c)
	if !ok {
		c.JSON(401, gin.H{
			"message": "permission denied.",
		})
		return
	}
	collectionP := client.Database("web_project").Collection("Post")
	collectionU := client.Database("web_project").Collection("User")
	var posts []Post
	filterP := bson.M{}
	curP, errP := collectionP.Find(context.TODO(), filterP)
	if errP != nil {
		c.JSON(409, gin.H{
			"message": "failed",
		})
		return
	}
	for curP.Next(context.TODO()) {
		var post Post
		if err = curP.Decode(&post); err != nil {
			return
		}
		posts = append(posts, post)
	}
	// fmt.Println(username)
	filterU := bson.M{
		"username": bson.M{
			"$eq": username,
		},
	}
	var user User
	err = collectionU.FindOne(context.TODO(), filterU).Decode(&user)
	if err != nil {
		c.JSON(409, gin.H{
			"message": "error",
		})
		return
	}
	var finalPosts []Post
	// fmt.Println(len(posts))
	// fmt.Println(user.Username)
	for i := range posts {
		for j := range user.Followings {
			if posts[i].Creator == user.Followings[j] {
				if posts[i].Parent == "" {
					finalPosts = append(finalPosts, posts[i])
				}
			}
		}
	}

	c.JSON(200, gin.H{
		"timeLine": convertPosts(finalPosts, user),
	})
}

func getFollowers(c *gin.Context) {
	username, ok := validateToken(c)
	if !ok {
		c.JSON(401, gin.H{
			"message": "permission denied.",
		})
		return
	}
	user := getUserById(username)
	c.JSON(200, gin.H{
		"followers": user.Followers,
	})
	return
}

func getFollowings(c *gin.Context) {
	username, ok := validateToken(c)
	if !ok {
		c.JSON(401, gin.H{
			"message": "permission denied.",
		})
		return
	}
	user := getUserById(username)
	c.JSON(200, gin.H{
		"followers": user.Followings,
	})
	return
}

func getBookmarls(c *gin.Context) {
	username, ok := validateToken(c)
	if !ok {
		c.JSON(401, gin.H{
			"message": "permission denied.",
		})
		return
	}
	user := getUserById(username)
	var posts []Post
	for i := range user.BookMark {
		posts = append(posts, getPostById(user.BookMark[i]))
	}
	c.JSON(200, gin.H{
		"posts": posts,
	})
	return
}

func convertPosts(posts []Post, user User) []PostForUser {
	var finalPosts []PostForUser
	for i := range posts {
		finalPosts = append(finalPosts, PostForUser{Id: posts[i].Id, Creator: posts[i].Creator, FullName: user.FullName, Content: posts[i].Content,
			Parent: posts[i].Parent, Likes: posts[i].Likes, Like: checkLike(posts[i], user.Username), Mark: checkMark(posts[i], user),
			LikeNumber: len(posts[i].Likes), ComentNumber: getComentNumber(posts[i].Id), Created_at: posts[i].Created_at})
	}
	return finalPosts
}

func getComentNumber(postId string) int {
	collectionP := client.Database("web_project").Collection("Post")
	filter := bson.M{
		"parent": bson.M{
			"$eq": postId,
		},
	}
	cur, err := collectionP.Find(context.TODO(), filter)
	fmt.Println(err)
	counter := 0
	for cur.Next(context.TODO()) {
		counter++
	}
	return counter
}

func checkLike(post Post, username string) bool {
	for i := range post.Likes {
		if post.Likes[i] == username {
			return true
		}
	}
	return false
}

func checkMark(post Post, user User) bool {
	for i := range user.BookMark {
		if user.BookMark[i] == post.Id {
			return true
		}
	}
	return false
}
