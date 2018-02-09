package main

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
)

var db *sql.DB
var err error

func main() {
	db, err = sql.Open("mysql", "slightlycyborg:foobar@/todo_service")

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

	t := New(title, completed)

	response := gin.H{
		"status":     http.StatusCreated,
		"message":    "Todo item created successfully!",
		"resourceId": t.ID}

	c.JSON(http.StatusCreated, response)
}

func handleFetchAllTodosRequest(c *gin.Context) {
	todos := All()
	response := gin.H{
		"data": todos}
	c.JSON(http.StatusFound, response)
}

func handleFetchSingleTodoRequest(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	todos := ByID(id)
	if len(todos) > 0 {
		response := gin.H{
			"data": todos[0]}
		c.JSON(http.StatusFound, response)
	} else {
		response := gin.H{
			"data": nil}
		c.JSON(http.StatusNotFound, response)
	}
}

func handleUpdateTodoRequest(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	title := c.PostForm("title")
	completed := c.PostForm("completed")

	todos := ByID(id)

	if len(todos) < 1 {
		response := gin.H{
			"message": "Todo not found"}
		c.JSON(http.StatusFound, response)
		return
	}

	t := todos[0]

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
	t := ByID(id)[0]
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

// Todo Struct and fns

type Todo struct {
	ID        int64
	Title     string
	Completed bool
	in_db     bool
}

//PUBLIC
func New(title string, completed bool) Todo {
	if db == nil {
		log.Fatal("Todo model not connected to DB. Implement another for of persistence perhaps?")
	}
	todo := Todo{ID: -1, Title: title, Completed: completed, in_db: false}
	id := todo.Persist()
	todo.ID = id
	todo.in_db = true

	return todo
}

func All() []Todo {
	sql_str, args, _ := sq.Select("*").From("todos").ToSql()
	rv := fromSQL(sql_str, args)
	return rv
}

func ByID(id int64) []Todo {
	sql_str, args, _ := sq.Select("*").From("todos").Where(sq.Eq{"ID": id}).ToSql()
	rv := fromSQL(sql_str, args)
	return rv
}

func (t Todo) Delete() bool {
	if t.ID == -1 {
		return false
	}

	query, args, _ := sq.Delete("todos").Where(sq.Eq{"ID": t.ID}).ToSql()

	_, err := db.Exec(query, args...)

	if err == nil {
		return true
	} else {
		return false
	}
}

func (t Todo) Persist() int64 {

	if !t.in_db {
		query, args, _ := sq.Insert("todos").
			Columns("title", "completed").
			Values(t.Title, t.Completed).ToSql()

		result, _ := db.Exec(query, args...)
		id, _ := result.LastInsertId()
		return id
	} else {
		query, args, _ := sq.Update("todos").
			Set("title", t.Title).
			Set("completed", t.Completed).
			Where(sq.Eq{"ID": t.ID}).ToSql()

		db.Exec(query, args...)
		return t.ID
	}
}

//private
func fromSQL(sql_str string, args []interface{}) []Todo {
	fmt.Println(args)
	var rows *sql.Rows
	var err error
	rows, err = db.Query(sql_str, args...)

	var rv []Todo
	print(rows)
	if err == nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var title string
		var completed bool

		err = rows.Scan(&id, &title, &completed)
		fmt.Print("\n\n<err>\n")
		fmt.Print(err)
		fmt.Print("</err>\n")

		rv = append(rv, Todo{ID: id, Title: title, Completed: completed, in_db: true})
	}
	return rv
}
