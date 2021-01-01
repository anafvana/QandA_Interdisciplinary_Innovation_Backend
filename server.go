package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	"math/rand"
	
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

type categories []category

type category struct {
	Name string `json:"cat"`
}

type keywords []keyword

type keyword struct {
	Name string `json:"kw"`
}

type entry struct {
	Question string
	Answer string
	SubmissionDate time.Time
	LastUpdateDate time.Time
	Categories categories
	KeyWords keywords
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

func (s *server) createTables() error {
	tableNames := []string {"entriesKeywords", "entriesCategories", "categories", "keywords", "entries"}
	for i := range tableNames{
		q := fmt.Sprintf("DROP TABLE IF EXISTS %v ;", tableNames[i])
		_, err := s.db.Exec(q)
		fmt.Println(err)
	}
	

	q1 := `
	CREATE TABLE entries(
		idEntry INT NOT NULL AUTO_INCREMENT, 
		question VARCHAR(600), 
		answer VARCHAR(4000),
		submission_date DATE,
		last_update_date DATE,
		PRIMARY KEY ( idEntry )
	);`

	q2 := `
	CREATE TABLE keywords(
		keyword NOT NULL VARCHAR(50),
		PRIMARY KEY ( keyword )
	);`

	q3 := `
	CREATE TABLE categories(
		category NOT NULL VARCHAR(50),
		PRIMARY KEY ( category )
	);`

	q4 := `
	CREATE TABLE entriesCategories(
		idEC INT NOT NULL AUTO_INCREMENT, 
		entryID INT NOT NULL,
		categoryID INT NOT NULL,
		PRIMARY KEY ( idEC ),
		FOREIGN KEY ( entryID ) REFERENCES entries( idEntry ),
		FOREIGN KEY ( categoryID ) REFERENCES categories( idCat )
	);`

	q5 := `
	CREATE TABLE entriesKeywords(
		idEKW INT NOT NULL AUTO_INCREMENT, 
		entryID INT NOT NULL,
		keywordID INT NOT NULL,
		PRIMARY KEY ( idEKW ),
		FOREIGN KEY ( entryID ) REFERENCES entries( idEntry ),
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

func (s *server) getEntries() error{
}

func (s *server) getCategoryNames(c echo.Context) error{
	rows, _ := s.db.Query("SELECT * FROM categories")

	var categoryName string
	var cats categories
	
	for rows.Next() {
		err := rows.Scan(&categoryName)
		if err != nil {
			fmt.Println(err)
		}
		cats = append(cats, category{Name: categoryName})
	}
	
	return c.JSON(http.StatusOK, cats)
}

func (s *server) getKeywordList (c echo.Context) error{
	rows, _ := s.db.Query("SELECT * FROM keywords")

	var kw string
	var kws keywords
	
	for rows.Next() {
		err := rows.Scan(&kw)
		if err != nil {
			fmt.Println(err)
		}
		kws = append(kws, keyword{Name: kw})
	}
	
	return c.JSON(http.StatusOK, kws)
}

func dummyCats() categories{
	return categories {
		{"cats"},
		{"dogs"},
		{"birds"},
		{"horses"}, 
		{"llamas"},
	}
}

func dummyKeywords() keywords{
	return keywords {
		{"hello"},
		{"darkness"},
		{"old"},
		{"friend"}, 
		{"world"}, 
		{"this"}, 
		{"dog"},
	}
}

func dummyQueries() {
	var data []entry
	var cc categories
	var kkww keywords

	kws := dummyKeywords()
	cats := dummyCats()


	for i := 0; i < 10; i++ {
		var e entry
		e = entry {
			Question: "nec ultrices dui sapien eget mi proin sed libero enim sed faucibus turpis in eu mi bibendum neque egestas congue quisque egestas diam in arcu cursus euismod quis viverra nibh cras pulvinar mattis nunc sed blandit libero volutpat sed cras ornare arcu dui vivamus arcu felis bibendum ut tristique et",
			Answer: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ullamcorper morbi tincidunt ornare massa eget. Sit amet volutpat consequat mauris nunc congue nisi vitae. Rhoncus mattis rhoncus urna neque viverra justo. Porttitor lacus luctus accumsan tortor posuere. Cras ornare arcu dui vivamus. Ultrices vitae auctor eu augue. Est placerat in egestas erat imperdiet. Ut eu sem integer vitae justo. Velit egestas dui id ornare arcu odio ut sem nulla. Et molestie ac feugiat sed lectus vestibulum. Nam aliquam sem et tortor consequat id. Aliquam vestibulum morbi blandit cursus. Magna eget est lorem ipsum dolor sit amet. Ac auctor augue mauris augue. Auctor neque vitae tempus quam pellentesque nec nam. Fermentum dui faucibus in ornare quam viverra orci sagittis eu. Metus dictum at tempor commodo ullamcorper a lacus vestibulum. Nisi vitae suscipit tellus mauris a diam maecenas. Quis lectus nulla at volutpat diam ut venenatis tellus in. Arcu non sodales neque sodales ut etiam sit amet. Lacus suspendisse faucibus interdum posuere. Et molestie ac feugiat sed lectus vestibulum mattis ullamcorper. Tristique senectus et netus et malesuada fames ac turpis. Duis ut diam quam nulla porttitor massa id neque. Cras fermentum odio eu feugiat. Sit amet massa vitae tortor condimentum. Sit amet est placerat in egestas erat imperdiet sed. Pellentesque sit amet porttitor eget dolor morbi. Non consectetur a erat nam at lectus.",			
		}
		c := rand.Intn(5-1) + 1
		kw := rand.Intn(7-0) + 0

		for j := 0; j < c; j++ {
			cc = append(cc, cats[j])
		}

		for k := 0; k < kw; k++{
			kkww = append(kkww, kws[k])
		}

		e.Categories = cc
		e.KeyWords = kkww
	}
}


func main() {
db, err := sql.Open("mysql", creds("credentials.json"))
	s := &server{
		e:	echo.New(),
		db:	db,
	}
	fmt.Println(err)
	s.createTables();
	//s.e.POST("/tables", s.createTables)
	//s.e.GET("/cat", s.getCategory)
	//s.e.GET("/kw", s.getKeyword)
	//s.e.Logger.Fatal(s.e.Start(":1323"))
}