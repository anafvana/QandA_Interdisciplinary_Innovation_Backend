package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func (s *server) createTables() {
	qMain := `
	CREATE TABLE questions(
		idQuest INT NOT NULL AUTO_INCREMENT, 
		question VARCHAR(600), 
		answer VARCHAR(4000)
		submission_date DATE,
		last_update_date DATE,
		PRIMARY KEY ( idQuest )
	);`
	
	res1, err1 := s.db.Exec(qMain)

	qKeyWords := `
	CREATE TABLE keywords(
		idKW INT NOT NULL AUTO_INCREMENT, 
		keyword VARCHAR(50),
		PRIMARY KEY ( idKW )
	);`

	res2, err2 := s.db.Exec(qKeyWords)

	qCategories := `
	CREATE TABLE categories(
		idCat INT NOT NULL AUTO_INCREMENT, 
		category VARCHAR(50),
		PRIMARY KEY ( idCat )
	);`

	res3, err3 := s.db.Exec(qCategories)

	qQuestCat := `
	CREATE TABLE questionsCategories(
		idQC INT NOT NULL AUTO_INCREMENT, 
		questionID INT NOT NULL,
		categoryID INT NOT NULL,
		PRIMARY KEY ( idQC ),
		FOREIGN KEY ( questionID ) REFERENCES questions( idQuest ),
		FOREIGN KEY ( categoryID ) REFERENCES categories( idCat )
	);`

	res4, err4 := s.db.Exec(qQuestCat)

	qQuestKW := `
	CREATE TABLE questionsKeywords(
		idQKW INT NOT NULL AUTO_INCREMENT, 
		questionID INT NOT NULL,
		keywordID INT NOT NULL,
		PRIMARY KEY ( idQKW ),
		FOREIGN KEY ( questionID ) REFERENCES questions( idQuest ),
		FOREIGN KEY ( keywordID ) REFERENCES keywords( idKW )
	);`

	res5, err5 := s.db.Exec(qQuestKW)
}
