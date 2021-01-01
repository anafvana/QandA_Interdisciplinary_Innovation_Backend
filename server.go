package main

import (
	//"net/http"
	"database/sql"
	
	"github.com/labstack/echo/v4"
	_ "github.com/go-sql-driver/mysql"
)

type server struct {
	e 	*echo.Echo
	db 	*sql.DB
}

func (s *server) getCategory(c echo.Context) error {
	return c.JSON(http.StatusOK, ???)
}

func (s *server) getKeyword(c echo.Context) error{
	return c.JSON(http.StatusOK, ???)
}

func main() {
	db, err := sql.Open()
	s := &server{
		e:	echo.New(),
		db:	db,
	}
	s.e.GET("/cat", s.getCategory)
	s.e.GET("/kw", s.getKeyword)
	s.e.Logger.Fatal(s.e.Start(":1323"))
}