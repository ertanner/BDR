package main

import (
	"fmt"
	"time"
	"os"
	"log"
	"encoding/csv"
	"bufio"
	"io"

)


func seasonality()  {
	fmt.Println("Season")

	//	Pick time period (number of years)
	//	Pick season period (month, quarter)
	//	Calculate average price for season (Year or number of years)
	//	Calculate average price over time (month or week)
	//	Divide season average by over time average price x 100
	var filename string = "myfile.csv"
	startDate := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	endDate  := time.Date(2016, time.November, 10, 23, 0, 0, 0, time.UTC)
	p1 :=	calcAvg(filename, startDate, endDate)
	p2m := 1.0 //calcMonth(filename)
	p2w := 1.0 //calcWeek(filename)

	seasonM := (p1/p2m) * 100
	seasonW := (p1/p2w) * 100
	fmt.Println(seasonM)
	fmt.Println(seasonW)
}

func calcAvg(filename string, startDate time.Time, endDate time.Time )( float64){
	var amt float64 = 0.0
	var cnt int = 0
	var avg float64
	f, err := os.Open(filename) // For read access.
	if err != nil {log.Fatal(err)}

	r := csv.NewReader(bufio.NewReader(f))
	for {
		record, err := r.Read()
		if err == io.EOF {break}
		if (startDate > record[2] && endDate < record[2]) {
			//amt = amt + record[1]
			cnt = cnt + 1
		}

	}
	avg = amt /  float64(cnt)
	return avg
}

//func calcMonth(filename string)(float64){
//	ioutil.ReadFile(filename)
//
//
//	return
//}
//
//func calcWeek(filename string)(float64){
//	ioutil.ReadFile(filename)
//	return
//}