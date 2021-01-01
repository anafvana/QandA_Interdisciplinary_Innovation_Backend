package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	//"net/http"
	"os"
	
	"github.com/labstack/echo/v4"
	_ "github.com/go-sql-driver/mysql"
)

type server struct {
	e 	*echo.Echo
	db 	*sql.DB
}

type cred struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

/*func (s *server) getCategory(c echo.Context) error {
	return c.JSON(http.StatusOK, ???)
}

func (s *server) getKeyword(c echo.Context) error{
	return c.JSON(http.StatusOK, ???)
}*/

func creds(fn string) string{
	f, err := os.Open(fn)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	bytes, _ := ioutil.ReadAll(f)

	var c cred

	json.Unmarshal(bytes, &c)
	return fmt.Sprintf("%v:%v@/%v", c.Username, c.Password, c.Database)
}

func main() {
	db, err := sql.Open("mysql", creds("credentials.json"))
	s := &server{
		e:	echo.New(),
		db:	db,
	}
	fmt.Println(err)
	s.CreateTables();
	//s.e.POST("/tables", s.createTables)
	//s.e.GET("/cat", s.getCategory)
	//s.e.GET("/kw", s.getKeyword)
	//s.e.Logger.Fatal(s.e.Start(":1323"))
}

func (s *server) CreateTables() error {
	tableNames := []string {"questionsKeywords", "questionsCategories", "categories", "keywords", "questions"}
	for i := range tableNames{
		q := fmt.Sprintf("DROP TABLE IF EXISTS %v ;", tableNames[i])
		_, err := s.db.Exec(q)
		fmt.Println(err)
	}
	

	q1 := `
	CREATE TABLE questions(
		idQuest INT NOT NULL AUTO_INCREMENT, 
		question VARCHAR(600), 
		answer VARCHAR(4000),
		submission_date DATE,
		last_update_date DATE,
		PRIMARY KEY ( idQuest )
	);`

	q2 := `
	CREATE TABLE keywords(
		idKW INT NOT NULL AUTO_INCREMENT, 
		keyword VARCHAR(50),
		PRIMARY KEY ( idKW )
	);`

	q3 := `
	CREATE TABLE categories(
		idCat INT NOT NULL AUTO_INCREMENT, 
		category VARCHAR(50),
		PRIMARY KEY ( idCat )
	);`

	q4 := `
	CREATE TABLE questionsCategories(
		idQC INT NOT NULL AUTO_INCREMENT, 
		questionID INT NOT NULL,
		categoryID INT NOT NULL,
		PRIMARY KEY ( idQC ),
		FOREIGN KEY ( questionID ) REFERENCES questions( idQuest ),
		FOREIGN KEY ( categoryID ) REFERENCES categories( idCat )
	);`

	q5 := `
	CREATE TABLE questionsKeywords(
		idQKW INT NOT NULL AUTO_INCREMENT, 
		questionID INT NOT NULL,
		keywordID INT NOT NULL,
		PRIMARY KEY ( idQKW ),
		FOREIGN KEY ( questionID ) REFERENCES questions( idQuest ),
		FOREIGN KEY ( keywordID ) REFERENCES keywords( idKW )
	);`
		
	_, err1 := s.db.Exec(q1)
	fmt.Println(err1)
	_, err2 := s.db.Exec(q2)
	fmt.Println(err2)
	_, err3 := s.db.Exec(q3)
	fmt.Println(err3)
	_, err4 := s.db.Exec(q4)
	fmt.Println(err4)
	_, err5 := s.db.Exec(q5)
	fmt.Println(err5)

	return err2
}
