package main

import (
	"net/http"
	"fmt"
	"log"
	"os"
	"io"
	"strings"
	"net/url"
)

var books = []string{
	"Genesis", "Exodus", "Leviticus", "Numbers", "Deuteronomy", "Joshua", "Judges", "Ruth", "1 Samuel", "2 Samuel", "1 Kings", "2 Kings", "1 Chronicles", "2 Chronicles", "Ezra", "Nehemiah", "Esther", "Job", "Psalm", "Proverbs", "Ecclesiastes", "Song of Solomon", "Isaiah", "Jeremiah", "Lamentations", "Ezekiel", "Daniel", "Hosea", "Joel", "Amos", "Obadiah", "Jonah", "Micah", "Nahum", "Habakkuk", "Zephaniah", "Haggai", "Zechariah", "Malachi",
	"Matthew", "Mark", "Luke", "John", "Acts", "Romans", "1 Corinthians", "2 Corinthians", "Galatians", "Ephesians", "Philippians", "Colossians", "1 Thessalonians", "2 Thessalonians", "1 Timothy", "2 Timothy", "Titus", "Philemon", "Hebrews", "James", "1 Peter", "2 Peter", "1 John", "2 John", "3 John", "Jude", "Revelation",
}
var chapters = []int{
	50, 40, 27, 36, 34, 24, 21, 4, 31, 24, 22, 25, 29, 36, 10, 13, 10, 42, 150, 31, 12, 8, 66, 52, 5, 48, 12, 14, 3, 9, 1, 4, 7, 3, 3, 3, 2, 14, 4,
	28, 16, 24, 21, 28, 16, 16, 13, 6, 6, 4, 4, 5, 3, 6, 4, 3, 1, 13, 5, 5, 3, 5, 1, 1, 1, 22,
}
var abbreviations = []string{
	"Gen","Exod","Lev","Num","Deut","Josh","Judg","Ruth","1Sam","2Sam","1Kgs","2Kgs","1Chr","2Chr","Ezra","Neh","Esth","Job","Ps","Prov","Eccl","Song","Isa","Jer","Lam","Ezek","Dan","Hos","Joel","Amos","Obad","Jonah","Mic","Nah","Hab","Zeph","Hag","Zech","Mal",
	"Matt","Mark","Luke","John","Acts","Rom","1Cor","2Cor","Gal","Eph","Phil","Col","1Thess","2Thess","1Tim","2Tim","Titus","Phlm","Heb","Jas","1Pet","2Pet","1John","2John","3John","Jude","Rev",
}
var u = "tqW6fep4K42h8anbCDhMIabiKWK3rH6nG9JjtUS7"
var p = "jackal11"

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	getMsgAll()
}

func getMsgAll() {
	for i, abbr := range abbreviations {
		maxChapter := chapters[i]
		for j := 1; j <= maxChapter; j++ {
			getMsg(abbr, books[i], j)
		}
	}
}

func getMsg(abbr, book string, chapter int) {
	lowerBook := strings.ToLower(book)
	filename := fmt.Sprint("data/msg/", lowerBook, "-", chapter, ".json")
	if _, err := os.Stat(filename); err == nil {
		log.Println("Skipping", filename)
		return
	}
	log.Println("Fetching", filename)

	r, err := http.Get(fmt.Sprint("https://", u, ":", p, "@bibles.org/v2/chapters/eng-MSG:", abbr, ".", chapter, "/verses.js"))
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = io.Copy(f, r.Body)
	if err != nil {
		log.Fatal(err)
	}
}

func getAllEsv() {
	for i, book := range books {
		maxChapter := chapters[i]
		for j := 1; j <= maxChapter; j++ {
			getEsv(book, j)
		}
	}
}

func getEsv(book string, chapter int) {
	lowerBook := strings.ToLower(book)
	filename := fmt.Sprint("data/esv/", lowerBook, "-", chapter, ".xml")
	if _, err := os.Stat(filename); err == nil {
		log.Println("Skipping", filename)
		return
	}
	log.Println("Fetching", filename)

	uri, err := url.Parse("http://www.esvapi.org/v2/rest/passageQuery")
	uri.Query().Add("key", "IP")
	uri.Query().Add("passage", fmt.Sprint(book, " ", chapter))
	uri.Query().Add("output-format", "crossway-xml-1.0")
	if err != nil {
		log.Fatal(err)
	}
	r, err := http.Get(uri.String())
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = io.Copy(f, r.Body)
	if err != nil {
		log.Fatal(err)
	}
}
