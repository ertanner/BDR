package main

import (
	"log"
	"fmt"
	"strconv"
)

func getProdCat(id string) map[int]string {
	var ID int
	var Name string
	prodCat1 := make(map[int]string)
	rows, err := db.Query("SELECT distinct CatID" +id+ ", Name FROM pcat"+id)
	if err != nil {
		log.Println("Erorr getting PCat" + id)
	}

	for rows.Next() {
		err := rows.Scan(&ID,&Name)
		if err != nil {
			log.Fatal(err)
		}
		prodCat1[ID] = Name
	}
	defer rows.Close()
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range prodCat1{
		fmt.Println( strconv.Itoa(k),v)
	}
	return prodCat1
}
