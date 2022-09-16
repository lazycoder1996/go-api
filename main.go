package main

import (
	"database/sql"
	"fmt"
	"go_dev/internal/controllers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db = initDb()

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "justtheoldpassword"
	dbname   = "data"
)

func initDb() *sql.DB {
	psqlconn := fmt.Sprintf("host = %s port = %d user = %s password = %s dbname = %s sslmode = disable", host, port, user, password, dbname)
	_db, err := sql.Open("postgres", psqlconn)

	CheckError(err)

	if e := _db.Ping(); e != nil {
		log.Fatalf("Error: %v", e)
	}

	return _db
}
func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	port := os.Getenv("PORT")
	defer db.Close()
	router := gin.Default()
	router.GET("/products", controllers.GetProducts(db))
	router.GET("/products/:guid", controllers.GetProduct(db))
	router.POST("/products", controllers.AddProduct(db))
	router.PUT("/products/:guid", controllers.UpdateProduct(db))
	router.DELETE("/products/:guid", controllers.DeleteProduct(db))

	log.Fatal(router.Run(":" + port))
}

// package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// 	_ "github.com/lib/pq"
// )

// type album struct {
// 	ID     string  `json:"id"`
// 	Title  string  `json:"title"`
// 	Artist string  `json:"artist"`
// 	Price  float64 `json:"price"`
// }

// var db = initDb()

// func initDb() *sql.DB {
// 	psqlconn := fmt.Sprintf("host = %s port = %d user = %s password = %s dbname = %s sslmode = disable", host, port, user, password, dbname)
// 	_db, err := sql.Open("postgres", psqlconn)

// 	CheckError(err)
// 	return _db
// }

// func getAlbums(c *gin.Context) {
// 	stmt := `select * from albums`
// 	rows, err := db.Query(stmt)
// 	switch err {
// 	case sql.ErrNoRows:
// 		c.IndentedJSON(http.StatusBadRequest, gin.H{
// 			"message": "no records found",
// 		})
// 	case nil:
// 		defer rows.Close()
// 		a := make([]album, 0)
// 		var rowsReadErr bool
// 		for rows.Next() {
// 			var id int
// 			var title string
// 			var artist string
// 			var price float64
// 			err = rows.Scan(&id, &title, &artist, &price)
// 			if err != nil {
// 				rowsReadErr = true
// 				break
// 			}
// 			a = append(a, album{ID: fmt.Sprintf("%d", id), Title: title, Artist: artist, Price: price})
// 		}
// 		if rowsReadErr {
// 			return
// 		}
// 		c.IndentedJSON(http.StatusOK, a)
// 	default:
// 		defer rows.Close()
// 	}
// }

// func postAlbum(c *gin.Context) {
// 	var newAlbum album
// 	if err := c.BindJSON(&newAlbum); err != nil {
// 		return
// 	}
// 	stmt := `insert into albums ("title", "artist", "price") values($1, $2, $3) returning id`
// 	lastInsertId := 0
// 	db.QueryRow(stmt, newAlbum.Title, newAlbum.Artist, newAlbum.Price).Scan(&lastInsertId)
// 	newAlbum.ID = fmt.Sprintf("%d", lastInsertId)
// 	c.IndentedJSON(http.StatusCreated, newAlbum)
// }

// func deleteAll(c *gin.Context) {
// 	stmt := `delete from albums`
// 	res, err := db.Exec(stmt)
// 	if err != nil {
// 		c.IndentedJSON(http.StatusInternalServerError, gin.H{
// 			"message": "error deleting rows",
// 		})
// 	}
// 	n, err := res.RowsAffected()
// 	if err != nil {
// 		c.IndentedJSON(http.StatusInternalServerError, gin.H{
// 			"message": "error occurred while checking the returned result from database after deletion",
// 		})
// 		return
// 	}

// 	if n == 0 {
// 		c.IndentedJSON(http.StatusBadRequest, gin.H{
// 			"message": "could not delete recored, there might be no records in table",
// 		})
// 		return
// 	}
// 	c.IndentedJSON(http.StatusOK, gin.H{
// 		"message": "successfully deleted all records",
// 	})
// }

// func deleteAlbumById(c *gin.Context) {
// 	id, _ := strconv.Atoi(c.Param("id"))
// 	stmt := `delete from albums where id = $1`
// 	res, err := db.Exec(stmt, id)
// 	if err != nil {
// 		c.IndentedJSON(http.StatusInternalServerError, gin.H{
// 			"message": "error occurred while deleting artist record",
// 		})
// 		return
// 	}
// 	n, err := res.RowsAffected()
// 	if err != nil {
// 		c.IndentedJSON(http.StatusInternalServerError, gin.H{
// 			"message": "error occurred while checking the returned result from database after deletion",
// 		})
// 		return
// 	}

// 	if n == 0 {
// 		c.IndentedJSON(http.StatusBadRequest, gin.H{
// 			"message": "could not delete recored, there might be no records for such id",
// 		})
// 		return
// 	}
// 	c.IndentedJSON(http.StatusOK, gin.H{
// 		"message": "successfully deleted the record",
// 	})

// }

// func getAlbumByID(c *gin.Context) {
// 	id, _ := strconv.Atoi(c.Param("id"))
// 	stmt := `select * from albums where id = $1`
// 	row := db.QueryRow(stmt, id)
// 	var artist, title string
// 	var price float64
// 	err := row.Scan(&id, &title, &artist, &price)
// 	switch err {
// 	case sql.ErrNoRows:
// 		c.IndentedJSON(http.StatusBadRequest, gin.H{
// 			"message": "no records found",
// 		})
// 	case nil:
// 		c.IndentedJSON(http.StatusOK, album{
// 			ID: fmt.Sprintf("%d", id), Artist: artist, Title: title, Price: price,
// 		})
// 	default:
// 		c.IndentedJSON(http.StatusOK, gin.H{"message": "album not found"})
// 	}
// }

// func updateAlbum(c *gin.Context) {
// 	var modifiedAlbum album
// 	if err := c.ShouldBindJSON(&modifiedAlbum); err != nil {
// 		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "invalid JSON body"})
// 		return
// 	}
// 	stmt := `update albums set title=$1, artist=$2, price=$3 where id=$4`
// 	res, err := db.Exec(stmt, modifiedAlbum.Title, modifiedAlbum.Artist, modifiedAlbum.Price, modifiedAlbum.ID)
// 	if err != nil {
// 		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error occurred while updating album"})
// 		return
// 	}
// 	n, _ := res.RowsAffected()
// 	if n == 0 {
// 		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "could not update please after some time"})
// 		return
// 	}
// 	c.IndentedJSON(http.StatusOK, gin.H{"message": "update made successfully"})
// }

// const (
// 	host     = "localhost"
// 	port     = 5432
// 	user     = "postgres"
// 	password = "justtheoldpassword"
// 	dbname   = "albums_db"
// )

// func main() {

// 	defer db.Close()

// 	router := gin.Default()
// 	router.GET("/albums", getAlbums)
// 	router.POST("/albums", postAlbum)
// 	router.GET("/albums/:id", getAlbumByID)
// 	router.DELETE("albums/delete/:id", deleteAlbumById)
// 	router.PUT("albums/update/:id", updateAlbum)
// 	router.DELETE("albums", deleteAll)
// 	router.Run("192.168.9.101:9090")
// }
