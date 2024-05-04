package main

import (
	"database/sql"
	"net/http"
	"simple-web-app-with-db/config"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Book struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	ReleaseYear string `json:"releaseYear"`
	Pages       int    `json:"pages"`
}

var db = config.ConnectDB()

func main() {
	router := gin.Default()

	defer db.Close()

	bookRouter := router.Group("/books")
	{
		bookRouter.GET("/", getAllBooks)
		bookRouter.GET("/:id", getBookById)
		bookRouter.POST("/", createBook)
		bookRouter.PUT("/:id", updateBook)
		bookRouter.DELETE("/:id", deleteBook)
	}

	router.Run(":3000") //default port :8080
}

// Handler untuk membuat buku baru
func createBook(c *gin.Context) {
	var newBook Book

	err := c.ShouldBind(&newBook)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "INSERT INTO mst_book (title, author, release_year, pages) VALUES ($1, $2, $3, $4) RETURNING id"

	var bookId int
	err = db.QueryRow(query, newBook.Title, newBook.Author, newBook.ReleaseYear, newBook.Pages).Scan(&bookId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create book"})
		return
	}

	newBook.ID = bookId
	c.JSON(http.StatusCreated, newBook)
}

// Handler untuk mengupdate buku baru
func updateBook(c *gin.Context) {
	id := c.Param("id")

	bookId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	// find book by id
	findBookQuery := "SELECT * FROM mst_book WHERE id = $1"
	book := Book{}

	err = db.QueryRow(findBookQuery, bookId).Scan(&book.ID, &book.Title, &book.Author, &book.ReleaseYear, &book.Pages)
	// check jika id tidak ada
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	var updatedBook Book
	err = c.ShouldBind(&updatedBook)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if strings.TrimSpace(updatedBook.Title) == "" {
		updatedBook.Title = book.Title
	}
	if strings.TrimSpace(updatedBook.Author) == "" {
		updatedBook.Author = book.Author
	}
	if strings.TrimSpace(updatedBook.ReleaseYear) == "" {
		updatedBook.ReleaseYear = book.ReleaseYear
	}
	if updatedBook.Pages == 0 {
		updatedBook.Pages = book.Pages
	}

	// update buku
	queryUpdate := "UPDATE mst_book SET title = $1, author = $2, release_year = $3, pages = $4 WHERE id = $5"
	_, err = db.Exec(queryUpdate, updatedBook.Title, updatedBook.Author, updatedBook.ReleaseYear, updatedBook.Pages, bookId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update book"})
	} else {
		updatedBook.ID = bookId
		c.JSON(http.StatusOK, gin.H{
			"message": "Book updated successfully",
			"data":    updatedBook,
		})
	}
}

// Handler untuk menampilkan semua buku atau buku berdasarkan pencarian judul
func getAllBooks(c *gin.Context) {
	searchTitle := c.Query("title")

	query := "SELECT id, title, author, release_year, pages FROM mst_book"

	var rows *sql.Rows
	var err error

	if searchTitle != "" {
		query += " WHERE title ILIKE '%' || $1 ||  '%'"
		rows, err = db.Query(query, searchTitle)
	} else {
		rows, err = db.Query(query)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	defer rows.Close()

	var matchBook []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.ReleaseYear, &book.Pages)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error"})
			return
		}
		matchBook = append(matchBook, book)
	}

	if len(matchBook) > 0 {
		c.JSON(http.StatusOK, matchBook)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
	}
}

// Handler untuk menampilkan buku berdasarkan id
func getBookById(c *gin.Context) {
	id := c.Param("id")

	bookId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	query := "SELECT * FROM mst_book WHERE id = $1"

	book := Book{}

	err = db.QueryRow(query, bookId).Scan(&book.ID, &book.Title, &book.Author, &book.ReleaseYear, &book.Pages)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, book)
}

// Handler untuk hapus buku
func deleteBook(c *gin.Context) {
	id := c.Param("id")

	bookId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}
	// find book by id
	findBookQuery := "SELECT * FROM mst_book WHERE id = $1"
	book := Book{}

	err = db.QueryRow(findBookQuery, bookId).Scan(&book.ID, &book.Title, &book.Author, &book.ReleaseYear, &book.Pages)

	// check jika id tidak ada
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	// delete Query
	deleteQuery := "DELETE FROM mst_book WHERE id = $1"
	db.Exec(deleteQuery, bookId)

	c.JSON(http.StatusOK, gin.H{
		"message": "Book deleted successfully",
	})
}
