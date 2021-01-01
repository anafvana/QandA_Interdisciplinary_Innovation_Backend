package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	
	"github.com/labstack/echo/v4"
	_ "github.com/go-sql-driver/mysql"
)

type server struct {
	e 	*echo.Echo
	db 	*sql.DB
}

type cred struct {
	username string `json:"username"`
	password string `json:"password"`
	database string `json:"database"`
}

func (s *server) getCategory(c echo.Context) error {
	return c.JSON(http.StatusOK, ???)
}

func (s *server) getKeyword(c echo.Context) error{
	return c.JSON(http.StatusOK, ???)
}

func creds(fn string) string{
	f, err := os.Open(fn)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	bytes, _ := ioutil.ReadAll(f)

	var c cred

	json.Unmarshal(bytes, &c)
	return fmt.Sprintf("%v:%v/%v", c.username, c.password, c.database)
}

func main() {
	db, err := sql.Open("mysql", creds("credentials.json"))
	s := &server{
		e:	echo.New(),
		db:	db,
	}
	s.e.GET("/cat", s.getCategory)
	s.e.GET("/kw", s.getKeyword)
	s.e.Logger.Fatal(s.e.Start(":1323"))
}