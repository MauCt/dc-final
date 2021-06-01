package api

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/MauCt/dc-final/controller"
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
	for i, _ := range users {
		if users[i].User == name {
			isBusy = true
		}
	}

	if isBusy || name == "" {
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

		c.JSON(200, gin.H{
			"message": "Bye " + name + ", your token has been revoked",
		})

		delete(users, tokenKey)
	} else {
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

		c.JSON(200, gin.H{
			"System Name": "Distributed Systems Class",
			"time":        time.Now(),
			"Workloads":   len(controller.Workloads),
		})

	} else {
		c.JSON(200, gin.H{
			"message": "Invalid token",
		})
	}
}

func CreateWorkload(c *gin.Context) {
	loginAuth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
	tokenKey := loginAuth[1]
	_, exist := users[tokenKey]
	if exist {
		workloadName := fmt.Sprintf("%v", len(controller.Workloads))
		filter := c.PostForm("filter")

		workloadFolder := "images/" + fmt.Sprintf("%v", len(controller.Workloads)) + "/"
		_ = os.MkdirAll(workloadFolder, 0755)

		newWL := controller.Workload{
			Id:       fmt.Sprintf("%v", len(controller.Workloads)),
			Filter:   filter,
			Name:     workloadName,
			Status:   "scheduling",
			Jobs:     0,
			Imgs:     []string{},
			Filtered: []string{},
		}

		controller.Workloads[fmt.Sprintf("%v", newWL.Id)] = newWL

		c.JSON(200, gin.H{
			"workload_id":     newWL.Id,
			"filter":          newWL.Filter,
			"workload_name":   newWL.Name,
			"status":          newWL.Status,
			"running_jobs":    newWL.Jobs,
			"filtered_images": []string{},
		})
	} else {
		c.JSON(200, gin.H{
			"message": "Invalid token",
		})
	}
}

func getWorkloads(c *gin.Context) {
	loginAuth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
	tokenKey := loginAuth[1]
	workload_id := c.Param("workload_id")
	_, exist := users[tokenKey]
	if exist {

		tempWorkload := controller.Workloads[workload_id]

		c.JSON(200, gin.H{
			"workload_id":     tempWorkload.Id + "\n",
			"filter":          tempWorkload.Filter,
			"workload_name":   tempWorkload.Name,
			"status":          tempWorkload.Status,
			"running_jobs":    tempWorkload.Jobs,
			"filtered_images": controller.Workloads[workload_id].Filtered,
		})
	} else {
		c.JSON(200, gin.H{
			"message": "Invalid token",
		})
	}
}

//Validates if the user exist using the token key and if the user exists it uploads the test.jpg image to the same folder.
/*func uploadImage(c *gin.Context) {
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
}*/

func Start() {
	r := gin.Default()
	r.POST("/login", login)
	r.DELETE("/logout", Logout)
	r.GET("/status", getStatus)
	r.POST("/workloads", CreateWorkload)
	r.GET("/workloads/:workload_id", getWorkloads)
	r.Run()

}
