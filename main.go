package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	gorm.Model
	Name     string
	Email    string
	Password string
}

type TransformedUser struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	//Migrate the schema
	db := Database()
	defer db.Close()
	db.AutoMigrate(&User{})

	router := gin.Default()
	v1 := router.Group("/api/v1/users")
	{
		v1.POST("/", CreateUser)
		v1.GET("/", GetUsers)
		// v1.GET("/:id", FetchSingleTodo)
		// v1.PUT("/:id", UpdateTodo)
		// v1.DELETE("/:id", DeleteTodo)
	}
	router.Run(":9293")

}

func Database() *gorm.DB {
	//open a db connection
	db, err := gorm.Open("mysql", "root:123456@/homestead?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func CreateUser(c *gin.Context) {
	user := User{Name: c.PostForm("name"), Email: c.PostForm("email"), Password: "$2y$10$SJPqybtTRV9gbXQBEcte5.6pYKh.0SDhYKytTgxunbSuq2iCzOQpC"}
	db := Database()
	defer db.Close()
	db.Save(&user)
	db.Close()

	status := http.StatusCreated
	message := "Create User Successfully"
	header := make(map[string]interface{})

	header["status"] = status
	header["message"] = message

	body := make(map[string]interface{})
	body["user_ID"] = user.ID
	body["user_name"] = user.Name
	body["user_email"] = user.Email
	c.JSON(status, gin.H{"header": header, "body": body})
}

func GetUsers(c *gin.Context) {
	var users []User
	var _users []TransformedUser

	db := Database()
	db.Find(&users)
	db.Close()

	if len(users) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "body": "No todo found!"})
		return
	}

	//transforms the users for building a good response
	for _, item := range users {
		_users = append(_users, TransformedUser{ID: item.ID, Name: item.Name, Email: item.Email})
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _users})

}
