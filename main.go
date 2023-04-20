package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Book struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Content string `json:"content"`
}

var db *gorm.DB
var err error

func main() {
	// Open a connection to the database
	db, err = gorm.Open("postgres", "host=127.0.0.1 user=postgres dbname=c2c3 sslmode=disable password=toor")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	db.AutoMigrate(&Book{})

	// Set up Gin router
	r := gin.Default()

	r.GET("/books", GetAllBooks)
	r.GET("/books/:id", GetBookByID)
	r.POST("/books", CreateBook)
	r.PUT("/books/:id", UpdateBook)
	r.DELETE("/books/:id", DeleteBook)

	r.Run(":8080")
}

func GetAllBooks(c *gin.Context) {
	var books []Book
	if err := db.Find(&books).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, books)
	}
}

func GetBookByID(c *gin.Context) {
	id := c.Params.ByName("id")
	var book Book
	if err := db.Where("id = ?", id).First(&book).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, book)
	}
}

func CreateBook(c *gin.Context) {
	var book Book
	c.BindJSON(&book)

	if err := db.Create(&book).Error; err != nil {
		c.AbortWithStatus(400)
		fmt.Println(err)
	} else {
		c.JSON(200, book)
	}
}

func UpdateBook(c *gin.Context) {
	id := c.Params.ByName("id")
	var book Book
	if err := db.Where("id = ?", id).First(&book).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	}

	c.BindJSON(&book)
	db.Save(&book)
	c.JSON(200, book)
}

func DeleteBook(c *gin.Context) {
	id := c.Params.ByName("id")
	var book Book
	if err := db.Where("id = ?", id).First(&book).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	}

	db.Delete(&book)
	c.JSON(200, gin.H{"id #" + id: "deleted"})
}
