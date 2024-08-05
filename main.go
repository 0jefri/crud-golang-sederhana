package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// membuat model
type User struct {
	gorm.Model
	Name  string `json:"name"`
	Email string `json:"email"`
}

// membuat koneksi ke database
var DB *gorm.DB
var err error

func main() {
	//koneksi ke database  SQLite
	DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		log.Fatal("gagal menghubungkan ke database")
	}

	//migrasi model ke databse
	DB.AutoMigrate(&User{})

	router := gin.Default()

	//routes
	router.POST("/users", CreateUser)
	router.GET("/users", GetUsers)
	router.GET("/users/:id", GetUserById)
	router.PUT("/users/:id", UpdateUserById)
	router.DELETE("/users/:id", DeleteUserById)

	//jalankan server
	router.Run()
}

// function handler
// function untuk membuat user baru
func CreateUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	DB.Create(&user)
	c.JSON(http.StatusOK, user)
}

// function handler all users
func GetUsers(c *gin.Context) {
	var users []User
	DB.Find(&users)
	c.JSON(http.StatusOK, users)
}

// function handler users by id
func GetUserById(c *gin.Context) {
	id := c.Param("id")
	var user User
	if result := DB.First(&user, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// function update data user berdasarkan id
func UpdateUserById(c *gin.Context) {
	id := c.Param("id")
	var user User
	if result := DB.First(&user, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Bind JSON data yang dikirimkan ke dalam struct user yang ditemukan
	var input User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Perbarui field yang sesuai
	user.Name = input.Name
	user.Email = input.Email
	//simpan perubahan e database
	DB.Save(&user)
	//kirim respon ke data yang sudah diperbarui
	c.JSON(http.StatusOK, user)
}

// function delete user by id
func DeleteUserById(c *gin.Context) {
	id := c.Param("id")
	var user User
	if result := DB.First(&user, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User Not Found"})
		return
	}
	DB.Delete(&user)
	c.JSON(http.StatusOK, gin.H{"message": "Success Delete User"})
}
