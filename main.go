package main

import (
	"github.com/gin-gonic/gin"
)

type Book struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	ReleaseYear string `json:"releaseYear"`
	Pages       int    `json:"pages"`
}

func main() {
	router := gin.Default()

	// Membuat buku baru
	router.POST("/books", createBook)

	// Menampilkan semua buku
	router.GET("/books", getAllBooks)

	router.Run(":8080")
}

// Handler untuk membuat buku baru
func createBook(c *gin.Context) {

}

// Handler untuk menampilkan semua buku atau buku berdasarkan pencarian judul
func getAllBooks(c *gin.Context) {

}
