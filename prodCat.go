package main

import (
	"log"
	"strconv"
)

func getCat(cat string) (map[string] string){
	var ID int
	var Name string

	pCat := make(map[string]string)
	rows, err := db.Query("SELECT distinct ID, Name FROM " + cat)
	if err != nil {
		log.Println("Erorr getting " + cat)
	}

	for rows.Next() {
		err := rows.Scan(&ID,&Name)
		if err != nil {
			log.Fatal(err)
		}
		pCat[ strconv.Itoa(ID)] = Name
	}
	defer rows.Close()
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return pCat
}
