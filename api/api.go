package api

import (
	"encoding/base64"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

//Struct to save the data of the users
type userData struct {
	User     string
	Password string
	Token    string
}

//Users map
var users = make(map[string]userData)

//Login function that takes the parameters and decode them to have the username and password.
//Validates if the user is already created.
func login(c *gin.Context) {

	loginAuth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
	l, _ := base64.StdEncoding.DecodeString(loginAuth[1])
	l2 := strings.SplitN(string(l), ":", 2)

	name := l2[0]
	password := l2[1]
	isBusy := false
	//Checks if the user exist
	for i, _ := range users {
		if users[i].User == name {
			isBusy = true
		}
	}

	if isBusy || name == "" {
		//print error
		c.JSON(200, gin.H{
			"message": "Username already taken",
		})
	} else {
		tokenNumber := loginAuth[1]
		users[tokenNumber] = userData{
			User:     name,
			Password: password,
			Token:    tokenNumber,
		}
		//print correct
		c.JSON(200, gin.H{
			"message": "Hi " + name + ", welcome to the DPIP System",
			"token":   tokenNumber,
		})
	}

}

//Logout function that uses the token key to see if the user exist or not.
func Logout(c *gin.Context) {
	loginAuth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
	tokenKey := loginAuth[1]

	_, exist := users[tokenKey]

	if exist {
		name := users[tokenKey].User
		//print

		c.JSON(200, gin.H{
			"message": "Bye " + name + ", your token has been revoked",
		})

		delete(users, tokenKey)
	} else {
		//print error
		c.JSON(200, gin.H{
			"message": "Invalid token",
		})
	}
}

//Status function use the token key to know if the user exist and gives the time of the day.
func getStatus(c *gin.Context) {
	loginAuth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
	tokenKey := loginAuth[1]
	_, exist := users[tokenKey]
	if exist {
		name := users[tokenKey].User

		c.JSON(200, gin.H{
			"message": "Hi " + name + ", the DPIP System is Up and Running",
			"time":    time.Now(),
		})

	} else {
		c.JSON(200, gin.H{
			"message": "Invalid token",
		})
	}
}

//Validates if the user exist using the token key and if the user exists it uploads the test.jpg image to the same folder.
func uploadImage(c *gin.Context) {
	loginAuth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
	tokenKey := loginAuth[1]
	_, exist := users[tokenKey]

	file, header, err := c.Request.FormFile("data")
	if err != nil {
		log.Fatal(err)
	}

	if exist {
		filename := header.Filename
		fileSize := header.Size

		imageOut, err := os.Create("copy" + filename)
		if err != nil {
			log.Fatal(err)
		}
		defer imageOut.Close()
		_, err = io.Copy(imageOut, file)
		if err != nil {
			log.Fatal(err)
		}
		fileSize = fileSize / 1000
		str := strconv.FormatInt(fileSize, 10)
		c.JSON(200, gin.H{
			"message":  "An image has been successfully uploaded",
			"filename": filename,
			"size":     str + "kb",
		})

	} else {
		c.JSON(200, gin.H{
			"message": "Invalid token",
		})
	}
}

func main() {
	r := gin.Default()
	r.GET("/login", login)
	r.GET("/logout", Logout)
	r.GET("/status", getStatus)
	r.POST("/upload", uploadImage)

	r.Run()

}

/*Links of the codes and information we use to make this project:
https://github.com/gin-gonic/gin
https://www.youtube.com/watch?v=RkmvVFZJJvs
https://gist.github.com/schollz/f25d77afc9130b72390748bdbce0d9a3
https://github.com/vaksi/go_auth_api
*/
