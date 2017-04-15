package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"log"
	"os"
	"encoding/json"
)

var books map[string][]string
// TODO this is not good...
// maps BOOK to map of CHAPTERS to map of VERSES to VERSE
var MSG map[string]map[string]map[string]string

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	f, err := os.Open("./data/books.json")
	if err != nil {
		log.Fatal(err)
	}
	if err = json.NewDecoder(f).Decode(&books); err != nil {
		log.Fatal(err)
	}
	f, err = os.Open("./data/MSG.json")
	if err != nil {
		log.Fatal(err)
	}
	if err = json.NewDecoder(f).Decode(&MSG); err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := gin.Default()
	r.NoRoute(func (c *gin.Context) {
		c.JSON(404, gin.H{"error": "not found"})
	})
	r.GET("/books", func(c *gin.Context) {
		c.JSON(200, books)
	})
	r.GET("/books/:book", func(c *gin.Context) {
		book := MSG[c.Param("book")]
		chapters := []string{}
		for c := range book {
			chapters = append(chapters, c)
		}
		c.JSON(200, chapters)
	})
	r.GET("/books/:book/:chapter", func(c *gin.Context) {
		c.JSON(200, MSG[c.Param("book")][c.Param("chapter")])
	})
	port := "8111"
	log.Println("Serving http://localhost:" + port)
	r.Run("0.0.0.0:" + port)
}
