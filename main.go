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

}

func handleDeleteTodoRequest(c *gin.Context) {

}
