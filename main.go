package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
}

func main() {
	router := gin.Default()

	// Connect to the SQLite database
	db, err := gorm.Open(sqlite.Open("todo.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Auto-migrate the Todo model to create the table
	db.AutoMigrate(&Todo{})

	// Route to create a new Todo
	router.POST("/todos", func(c *gin.Context) {
		var todo Todo
		if err := c.ShouldBindJSON(&todo); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON data"})
			return
		}

		// Save the Todo to the database
		db.Create(&todo)

		c.JSON(200, todo)
	})

	// Route to get all Todos
	router.GET("/todos", func(c *gin.Context) {
		var todos []Todo

		// Retrieve all Todos from the database
		db.Find(&todos)

		c.JSON(200, todos)
	})

	// Route to get a specific Todo by ID
	router.GET("/todos/:id", func(c *gin.Context) {
		var todo Todo
		todoID := c.Param("id")

		// Retrieve the Todo from the database
		result := db.First(&todo, todoID)
		if result.Error != nil {
			c.JSON(404, gin.H{"error": "Todo not found"})
			return
		}

		c.JSON(200, todo)
	})

	// Route to update a Todo by ID
	router.PUT("/todos/:id", func(c *gin.Context) {
		var todo Todo
		todoID := c.Param("id")

		// Retrieve the Todo from the database
		result := db.First(&todo, todoID)
		if result.Error != nil {
			c.JSON(404, gin.H{"error": "Todo not found"})
			return
		}

		var updatedTodo Todo
		if err := c.ShouldBindJSON(&updatedTodo); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON data"})
			return
		}

		// Update the Todo in the database
		todo.Title = updatedTodo.Title
		todo.Description = updatedTodo.Description
		db.Save(&todo)

		c.JSON(200, todo)
	})

	// Route to delete a Todo by ID
	router.DELETE("/todos/:id", func(c *gin.Context) {
		var todo Todo
		todoID := c.Param("id")

		// Retrieve the Todo from the database
		result := db.First(&todo, todoID)
		if result.Error != nil {
			c.JSON(404, gin.H{"error": "Todo not found"})
			return
		}

		// Delete the Todo from the database
		db.Delete(&todo)

		c.JSON(200, gin.H{"message": fmt.Sprintf("Todo with ID %s deleted", todoID)})
	})

	router.Run(":8080")
}

