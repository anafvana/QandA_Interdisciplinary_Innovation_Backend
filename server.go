package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

type server struct {
	e  *echo.Echo
	db *sql.DB
}

type cred struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

//type categories []category

type category struct {
	Name string `json:"cat"`
}

//type keywords []keyword

type keyword struct {
	Name string `json:"kw"`
}

type entry struct {
	ID             string     `json:"_id"`
	Question       string     `json:"question"`
	Answer         string     `json:"answer"`
	SubmissionDate time.Time  `json:"submissionDate"`
	LastUpdate     time.Time  `json:"lastUpdate"`
	Categories     []category `json:"categories"`
	KeyWords       []keyword  `json:"keywords"`
}

/*func (s *server) getCategory(c echo.Context) error {
	return c.JSON(http.StatusOK, ???)
}

func (s *server) getKeyword(c echo.Context) error{
	return c.JSON(http.StatusOK, ???)
}*/

/*-----------------------	DATABASE	-----------------------*/
//Fetches credentials to log into database
func creds(fn string) string {
	f, err := os.Open(fn)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	bytes, _ := ioutil.ReadAll(f)

	var c cred

	json.Unmarshal(bytes, &c)
	return fmt.Sprintf("%v:%v@/%v", c.Username, c.Password, c.Database)
}

//Database set-up
func (s *server) createTables() error {
	//TODO delete this before official deployment
	tableNames := []string{"entriesKeywords", "entriesCategories", "categories", "keywords", "entries"}
	for i := range tableNames {
		q := fmt.Sprintf("DROP TABLE IF EXISTS %v ;", tableNames[i])
		_, err := s.db.Exec(q)
		if err != nil {
			log.Println(err)
		}
	}

	q1 := `
	CREATE TABLE entries(
		entryID INT NOT NULL AUTO_INCREMENT, 
		question VARCHAR(600), 
		answer VARCHAR(4000),
		submission_date DATETIME,
		last_update DATETIME,
		PRIMARY KEY ( entryID )
	);`

	q2 := `
	CREATE TABLE keywords(
		keyword VARCHAR(50) NOT NULL,
		PRIMARY KEY ( keyword )
	);`

	q3 := `
	CREATE TABLE categories(
		category VARCHAR(50) NOT NULL,
		PRIMARY KEY ( category )
	);`

	q4 := `
	CREATE TABLE entriesCategories(
		idEC INT NOT NULL AUTO_INCREMENT, 
		entryID INT NOT NULL,
		category VARCHAR(50) NOT NULL,
		PRIMARY KEY ( idEC ),
		FOREIGN KEY ( entryID ) REFERENCES entries( entryID ),
		FOREIGN KEY ( category ) REFERENCES categories( category )
	);`

	q5 := `
	CREATE TABLE entriesKeywords(
		idEKW INT NOT NULL AUTO_INCREMENT, 
		entryID INT NOT NULL,
		keyword VARCHAR(50) NOT NULL,
		PRIMARY KEY ( idEKW ),
		FOREIGN KEY ( entryID ) REFERENCES entries( entryID ),
		FOREIGN KEY ( keyword ) REFERENCES keywords( keyword )
	);`

	_, err := s.db.Exec(q1)
	if err != nil {
		log.Println(err)
	}
	_, err = s.db.Exec(q2)
	if err != nil {
		log.Println(err)
	}
	_, err = s.db.Exec(q3)
	if err != nil {
		log.Println(err)
	}
	_, err = s.db.Exec(q4)
	if err != nil {
		log.Println(err)
	}
	_, err = s.db.Exec(q5)
	if err != nil {
		log.Println(err)
	}

	return err
}

//Inserting entry in Database
func (s *server) newEntryDB(e entry) error {
	_, err := s.db.Exec(`
	INSERT INTO entries (question, answer, submission_date, last_update)
	VALUES (
		?, ?, ?, ?
	);`, e.Question, e.Answer, e.SubmissionDate, e.LastUpdate)
	if err != nil {
		log.Println(err)
	}

	var entryID string

	err = s.db.QueryRow(`
		SELECT LAST_INSERT_ID()
		FROM entries;
	`).Scan(&entryID)
	if err != nil {
		log.Println(err)
	}

	for i := range e.Categories {
		b, err := s.checkCategory(e.Categories[i])
		if err != nil {
			log.Println(err)
		}
		if !b {
			_, err := s.db.Exec(`
			INSERT INTO categories
			VALUES (
				?
			);`, e.Categories[i].Name)
			if err != nil {
				log.Println(err)
			}
		}

		_, err = s.db.Exec(`
			INSERT INTO entriesCategories
			SET entryID = (
					SELECT entryID
					FROM entries
					WHERE entryID = ?
			),
				category = (
					SELECT category
					FROM categories
					WHERE category = ?
				);`, entryID, e.Categories[i].Name)
		if err != nil {
			log.Println(err)
		}
	}

	for i := range e.KeyWords {
		b, err := s.checkKeyword(e.KeyWords[i])
		if err != nil {
			log.Println(err)
		}
		if !b {
			_, err := s.db.Exec(`
			INSERT INTO keywords
			VALUES (
				?
			);`, e.KeyWords[i].Name)
			if err != nil {
				log.Println(err)
			}
		}

		_, err = s.db.Exec(`
			INSERT INTO entriesKeywords
			SET entryID = (
					SELECT entryID
					FROM entries
					WHERE entryID = ?
			),
				keyword = (
					SELECT keyword
					FROM keywords
					WHERE keyword = ?
				);`, entryID, e.KeyWords[i].Name)
		if err != nil {
			log.Println(err)
		}
	}
	return err
}

func (s *server) checkCategory(c category) (bool, error) {
	var cat string
	err := s.db.QueryRow(`
		SELECT COUNT(category) FROM categories WHERE category=? ;
	`, c.Name).Scan(&cat)
	if err != nil {
		log.Println(err)
	}

	catNr, err := strconv.ParseInt(cat, 6, 12)

	var b bool
	if catNr == 0 {
		b = false
	} else {
		b = true
	}

	return b, err
}

func (s *server) checkKeyword(kw keyword) (bool, error) {
	var k string
	err := s.db.QueryRow(`
		SELECT COUNT(keyword) FROM keywords WHERE keyword=? ;
	`, kw.Name).Scan(&k)
	if err != nil {
		log.Println(err)
	}

	kNr, err := strconv.ParseInt(k, 6, 12)

	var b bool
	if kNr == 0 {
		b = false
	} else {
		b = true
	}

	return b, err
}

//Fetching entry from Database
//Read as a struct
func (s *server) fetchEntry(id string) entry {
	var e entry

	//Fetches data from entries
	var sdSTR string
	var luSTR string
	err := s.db.QueryRow(`
		SELECT * FROM entries WHERE entryID=? ;
	`, id).Scan(&e.ID, &e.Question, &e.Answer, &sdSTR, &luSTR)
	if err != nil {
		log.Println(err)
	}
	layout := "2006-01-02 15:04:05"
	e.SubmissionDate, err = time.Parse(layout, sdSTR)
	if err != nil {
		log.Println(err)
	}

	e.LastUpdate, err = time.Parse(layout, luSTR)
	if err != nil {
		log.Println(err)
	}

	//Fetches categories
	e.Categories = s.fetchEntryCategory(e.ID)

	//Fetches keywords
	e.KeyWords = s.fetchEntryKeywords(e.ID)

	fmt.Println(e)
	return e
}

func (s *server) fetchEntryCategory(id string) []category {
	var cats []category

	//Fetches categories
	rows, err := s.db.Query(`
		SELECT category FROM entriesCategories WHERE entryID=? ;
	`, id)
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		var cat string
		err := rows.Scan(&cat)
		if err != nil {
			log.Println(err)
		}
		cats = append(cats, category{Name: cat})
	}

	return cats
}

func (s *server) fetchEntryKeywords(id string) []keyword {
	var kws []keyword

	//Fetches categories
	rows, err := s.db.Query(`
		SELECT keyword FROM entriesKeywords WHERE entryID=? ;
	`, id)
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		var kw string
		err := rows.Scan(&kw)
		if err != nil {
			log.Println(err)
		}
		kws = append(kws, keyword{Name: kw})
	}

	return kws
}

//Fetching data from Database
//Read into JSON
/* func (s *server) getEntries() error{
} */

func (s *server) getCategoryNames(c echo.Context) error {
	rows, _ := s.db.Query("SELECT * FROM categories;")

	var categoryName string
	var cats []category

	for rows.Next() {
		err := rows.Scan(&categoryName)
		if err != nil {
			log.Println(err)
		}
		cats = append(cats, category{Name: categoryName})
	}

	return c.JSON(http.StatusOK, cats)
}

func (s *server) getKeywordList(c echo.Context) error {
	rows, _ := s.db.Query("SELECT * FROM keywords;")

	var kw string
	var kws []keyword

	for rows.Next() {
		err := rows.Scan(&kw)
		if err != nil {
			log.Println(err)
		}
		kws = append(kws, keyword{Name: kw})
	}

	return c.JSON(http.StatusOK, kws)
}

/*-----------------------	JSON	-----------------------*/
//Reads all entries from a JSON file
func readEntriesJSON(fn string) []entry {
	f, err := os.Open(fn)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	bytes, _ := ioutil.ReadAll(f)

	var allEntries []entry

	err = json.Unmarshal(bytes, &allEntries)
	if err != nil {
		log.Println(err)
	}

	//TODO: Erase test
	/* for i := range allEntries {
		fmt.Println(allEntries[i])
	} */

	return allEntries
}

//Handles unmarshalling of categories
func (c *category) UnmarshalJSON(data []byte) error {
	var v string
	err := json.Unmarshal(data, &v)
	if err != nil {
		log.Println(err)
	}
	c.Name = v
	return err
}

//Handles unmarshalling of keywords
func (kw *keyword) UnmarshalJSON(data []byte) error {
	var v string
	err := json.Unmarshal(data, &v)
	if err != nil {
		log.Println(err)
	}
	kw.Name = v
	return err
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	db, err := sql.Open("mysql", creds("credentials.json"))
	s := &server{
		e:  echo.New(),
		db: db,
	}
	if err != nil {
		log.Println(err)
	}
	s.createTables()
	entries := readEntriesJSON("dummy.json")
	for i := range entries {
		s.newEntryDB(entries[i])
	}
	s.fetchEntry("2")
	//s.e.POST("/tables", s.createTables)
	//s.e.GET("/cat", s.getCategory)
	//s.e.GET("/kw", s.getKeyword)
	//s.e.Logger.Fatal(s.e.Start(":1323"))
}
