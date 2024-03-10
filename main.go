package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
)

type HelloHandler struct{}
type WorldHandler struct{}
type Todo struct {
	ID int `json:"id"`
	Title string `json:"title"`
}

func (h *HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true",
		"root",
		"",
		"localhost",
		"3306",
		"test",
	)

	connection, err := sql.Open("mysql", dsn)

	if err != nil {
		fmt.Println(err)
		panic("Failed to connect to database")
	}

	rows, err := connection.Query("select * from todos;")

	todos := []Todo{}

	if err != nil {
		panic(err)
	}

	for rows.Next() {
		todo := Todo{}
		rows.Scan(&todo.ID, &todo.Title)
		todos = append(todos, todo)
	}

	rows.Close()

	response, _ := json.Marshal(todos)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (h *WorldHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "World")
}

func main() {
	hello := HelloHandler{}
	world := WorldHandler{}

	server := http.Server{
		Addr: ":8080",
		Handler: nil,
	}

	http.Handle("/hello", &hello)
	http.Handle("/world", &world)

	server.ListenAndServe()
}
