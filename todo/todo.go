package todo

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"log"
)

var db *sql.DB = nil

func AddDB(db_ *sql.DB) {
	db = db_
}

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
