package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
	"todo-service/todo"
)

var db *sql.DB
var err error

func main() {
	db, err = sql.Open("mysql", "slightlycyborg:foobar@/todo_service")

	//Initialize the Todo Model
	todo.AddDB(db)

	router := gin.Default()

	v1 := router.Group("/api/v1/todos")
	{
		v1.POST("/", handleCreateTodoRequest)
		v1.GET("/", handleFetchAllTodosRequest)
		v1.GET("/:id", handleFetchSingleTodoRequest)
		v1.PUT("/:id", handleUpdateTodoRequest)
		v1.DELETE("/:id", handleDeleteTodoRequest)
	}
	router.Run()
}

func handleCreateTodoRequest(c *gin.Context) {
	completed, _ := strconv.ParseBool(c.PostForm("completed"))
	title := c.PostForm("title")

	t := todo.New(title, completed)

	response := gin.H{
		"status":     http.StatusCreated,
		"message":    "Todo item created successfully!",
		"resourceId": t.ID}

	c.JSON(http.StatusCreated, response)
}

func handleFetchAllTodosRequest(c *gin.Context) {
	todos := todo.All()
	response := gin.H{
		"data": todos}
	c.JSON(http.StatusFound, response)
}

func handleFetchSingleTodoRequest(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	todos := todo.ByID(id)
	if len(todos) > 0 {
		response := gin.H{
			"data": todos[0]}
		c.JSON(http.StatusFound, response)
	}

}

func handleUpdateTodoRequest(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	title := c.PostForm("title")
	completed := c.PostForm("completed")

	//This may be slow to do a fetch and then an update.
	//The ByID does not have to actually retrieve the data from the database
	//This is the beauty of abstraction.
	//We can an build accesors to get the attributes as needed
	t := todo.ByID(id)[0]

	if title != "" {
		t.Title = title
	}
	if completed != "" {
		b_completed, err := strconv.ParseBool(completed)
		if err == nil {
			t.Completed = b_completed
		}
	}
	t.Persist()

	response := gin.H{
		"message": "Todo updated successfully"}
	c.JSON(http.StatusOK, response)
}

func handleDeleteTodoRequest(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	t := todo.ByID(id)[0]
	b_deleted := t.Delete()

	var response gin.H
	var status int

	if b_deleted {
		response = gin.H{
			"message": "Todo deleted successfully"}

		status = http.StatusOK

	} else {
		response = gin.H{
			"message": "Todo was not deleted"}

		status = http.StatusInternalServerError
	}

	c.JSON(status, response)

}
