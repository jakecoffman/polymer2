package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"log"
	"os"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"github.com/jmoiron/sqlx"
	"strconv"
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

type Book struct {
	BookNumber int `db:"book_number"`
	ShortName string `db:"short_name"`
	LongName string `db:"long_name"`
}

func main() {
	esvdb := mustOpen("./data/ESV.SQLite3")
	msgdb := mustOpen("./data/MSG.SQLite3")
	defer esvdb.Close()
	defer msgdb.Close()

	versionMap := map[string]*sqlx.DB{
		"esv": esvdb,
		"msg": msgdb,
	}

	r := gin.Default()
	r.NoRoute(func (c *gin.Context) {
		c.JSON(404, gin.H{"error": "not found"})
	})
	// list books
	r.GET("/books", func(c *gin.Context) {
		var books []Book
		err := esvdb.Select(&books, "select book_number, short_name, long_name from books order by book_number")
		if err != nil {
			c.JSON(500, map[string]string{"error": err.Error()})
			return
		}
		c.JSON(200, books)
	})
	// list chapters
	r.GET("/versions/:version/:book", func(c *gin.Context) {
		version := c.Param("version")
		selectedDb, ok := versionMap[version]
		if !ok {
			c.JSON(422, gin.H{"error": "bad version string"})
			return
		}
		book, err := strconv.Atoi(c.Param("book"))
		if err != nil {
			c.JSON(422, gin.H{"error": "book needs to be an int"})
			return
		}
		chapters := []string{}
		err = selectedDb.Select(&chapters, "select distinct chapter from verses where book_number=?", book)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, chapters)
	})
	// list verses of chapter
	r.GET("/versions/:version/:book/:chapter", func(c *gin.Context) {
		version := c.Param("version")
		chapter := c.Param("chapter")
		book, err := strconv.Atoi(c.Param("book"))
		if err != nil {
			c.JSON(422, gin.H{"error": "book needs to be an int"})
			return
		}
		selectedDb, ok := versionMap[version]
		if !ok {
			c.JSON(422, gin.H{"error": "bad version string"})
			return
		}
		var verses []Verse
		err = selectedDb.Select(&verses, "select verse, text from verses where book_number=? and chapter=?", book, chapter)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, verses)
	})
	port := "8111"
	log.Println("Serving http://localhost:" + port)
	r.Run("0.0.0.0:" + port)
}

type Verse struct {
	Verse int `db:"verse"`
	Text string `db:"text"`
}

func mustOpen(dbname string) *sqlx.DB {
	ddb, err := sqlx.Open("sqlite3", dbname)
	if err != nil {
		log.Fatal(err)
	}
	return ddb
}
