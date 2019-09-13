package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" //ให้เรียก function init ของ SQL เฉยๆ
)

//export DATABASE_URL=postgres://losikmkt:OLxHUThWjfIpCVskUJlNfbaCeD1VgEjR@arjuna.db.elephantsql.com:5432/losikmkt

//DB Pointer
var db *sql.DB
var stmt *sql.Stmt
var rows *sql.Rows

var err error
var globalTestInsertCount int

type Todo struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

func connectDB() {
	url := os.Getenv("DATABASE_URL")
	fmt.Println("url :", url)
	db, err = sql.Open("postgres", url)
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
}

//Query string parameters
func getAllTodosHandler(c *gin.Context) {

	fmt.Println("Console => Start getAllTodosHandler")
	//http://localhost:1234/todos?title=buysnack
	//t1 := c.DefaultQuery("title", "") //get Query para from URL
	var todosSelected []Todo
	var todoItem Todo
	var sqlStatement string

	t1 := c.Query("title") //get Query para from URL

	if t1 == "" {
		sqlStatement = "SELECT id, title, status From todos"
		stmt, err = db.Prepare(sqlStatement)
		if err != nil {
			//log.Fatal("can't prepare query all statement : ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "can't prepare query all statement : " + err.Error()})
			c.Abort()
			return
		}
		rows, err = stmt.Query()
		if err != nil {
			//log.Fatal("can't query all statement : ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "can't query all statement : " + err.Error()})
			c.Abort()
			return
		}

	} else {
		t1 = "%" + t1 + "%"
		sqlStatement = "SELECT id, title, status From todos where title LIKE $1"
		stmt, err = db.Prepare(sqlStatement)
		if err != nil {
			//log.Fatal("can't prepare query where statement : ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "can't prepare query where statement : " + err.Error()})
			c.Abort()
			return
		}
		rows, err = stmt.Query(t1)
		if err != nil {
			//log.Fatal("can't query where statement : ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "can't query where statement : " + err.Error()})
			c.Abort()
			return
		}
	}

	for rows.Next() {
		err := rows.Scan(&todoItem.ID, &todoItem.Title, &todoItem.Status)
		if err != nil {
			//log.Fatal("can't Scan row into variable : ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "can't Scan into variable : " + err.Error()})
			c.Abort()
			return
		}
		fmt.Println("id: ", todoItem.ID, "\n title: ", todoItem.Title, "\n status: ", todoItem.Status, "\n==========")
		todosSelected = append(todosSelected, todoItem)
	}
	fmt.Println("Console => End   getAllTodosHandler => query all todos success")
	c.JSON(http.StatusOK, todosSelected)
}

//Query string Parameters in Path
func getPathTodosHandler(c *gin.Context) {

	fmt.Println("Console => Start getPathTodosHandler")

	var todoItem Todo

	stmt, err = db.Prepare("SELECT id, title, status From todos where id=$1")
	if err != nil {
		//log.Fatal("can't prepare query one statement : ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't prepare query one statement : " + err.Error()})
		c.Abort()
		return
	}

	rowId := c.Param("id")
	row := stmt.QueryRow(rowId)
	err = row.Scan(&todoItem.ID, &todoItem.Title, &todoItem.Status)
	if err != nil {
		//log.Fatal("can't Scan row into variables : ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't prepare query one statement : " + err.Error()})
		c.Abort()
		return
	}

	fmt.Println("one row\n id: ", todoItem.ID, "\n title: ", todoItem.Title, "\n status: ", todoItem.Status)
	c.JSON(http.StatusOK, todoItem)
	fmt.Println("Console => End   getPathTodosHandler => quey one todos success")
}
func createTodosHandler(c *gin.Context) {
	fmt.Println("Console => Start createTodosHandler")

	var todoItem Todo
	err := c.ShouldBindJSON(&todoItem) // read "body" and "Unmarshal" to "struct t"
	if err != nil {
		//fmt.Println("error in puteTodosHandler")
		//log.Fatal("can't Bind Body to JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't Bind Body to JSON : " + err.Error()})
		c.Abort()
		return
	}
	//globalTestInsertCount += 1
	//todoItem.Title = "buy test insert" + strconv.Itoa(globalTestInsertCount)
	//todoItem.Status = "active"
	row := db.QueryRow("INSERT INTO todos (title, status) values ($1, $2 ) RETURNING id", todoItem.Title, todoItem.Status)
	err = row.Scan(&todoItem.ID)
	if err != nil {
		//fmt.Println("can't scan id : ", err)
		//log.Fatal("can't scan : ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't scan : " + err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "inserted"})
	fmt.Println("insert todo success id : ", todoItem.ID)
	fmt.Println("Console => End   createTodosHandler => insert one todos success")
}
func deleteTodosHandler(c *gin.Context) {
	fmt.Println("Console => Start deleteTodosHandler")
	rowId := c.Param("id")
	stmt, err = db.Prepare("DELETE FROM todos WHERE id=$1")

	numId, err := strconv.Atoi(rowId)
	if err != nil {
		//log.Fatal("can't convert to int -", rowId, "- : ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't convert to int -" + rowId + "- : " + err.Error()})
		c.Abort()
		return
	}

	//_, err = db.ExecContext(c, "DELETE FROM todos WHERE id=$1" , numId)
	_, err = stmt.ExecContext(c, numId)
	if err != nil {
		//fmt.Println("can't execContext : ", err)
		//log.Fatal("can't execContext : ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't execContext : " + err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	fmt.Println("delete todo success id : ", rowId)
	fmt.Println("Console => End   deleteTodosHandler => delete one todos success")
}
func putTodosHandler(c *gin.Context) {
	fmt.Println("Console => Start putTodosHandler")

	var todoItem Todo

	err := c.ShouldBindJSON(&todoItem) // read "body" and "Unmarshal" to "struct t"
	if err != nil {
		//fmt.Println("error in puteTodosHandler")
		//log.Fatal("can't Bind Body to JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't Bind Body to JSON : " + err.Error()})
		c.Abort()
		return
	}

	stmt, err = db.Prepare("UPDATE todos SET title=$2, status=$3 WHERE id=$1")
	if err != nil {
		//log.Fatal("can't prepare statement update : ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't prepare update statement: " + err.Error()})
		c.Abort()
		return
	}

	if _, err := stmt.Exec(todoItem.ID, todoItem.Title, todoItem.Status); err != nil {
		//log.Fatal("error execute update", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't execute update statement: " + err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
	fmt.Println("Console => End   putTodosHandler  => update todos success")
}

func createTable(c *gin.Context) {
	fmt.Println("Console => Start createTable")

	createTb := `
	CREATE TABLE IF NOT EXISTS todos (
		id SERIAL PRIMARY KEY,
		title TEXT,
		status TEXT
	);
	`

	_, err := db.Exec(createTb)

	if err != nil {
		//log.Fatal("can't create table : ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't create table : " + err.Error()})
		c.Abort()
		return
	}

	fmt.Println("Console => End   createTable  => create table success")
}

func authMiddleware(c *gin.Context) {
	fmt.Println("This is a middleware")
	token := c.GetHeader("Authorization")
	if token != "Bearer token123" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized."})
		c.Abort()
		return
	}
	c.Next()
	fmt.Println("after in middleware")
}

func main() {
	connectDB()
	r := gin.Default()    // return Engine ทีเป็น struct กลับมา
	r.Use(authMiddleware) //Using middleware
	r.GET("/todos", getAllTodosHandler)
	r.GET("/todos/:id", getPathTodosHandler)
	r.POST("/todos", createTodosHandler)
	r.DELETE("/todos/:id", deleteTodosHandler)
	r.PUT("/todos", putTodosHandler)
	r.POST("/create_todos_table", createTable)
	r.Run(":1234") // listen and serve on 0.0.0.0:8080
	defer db.Close()
}
