package main

import (
    "bufio"
    "encoding/csv"
    "fmt"
    "io"
    "log"
    "os"
    "time"
    "strings"
    "strconv"
)

type DataPoint struct {
	timestamp time.Time
	data string
}

func main() {
	csvFile, _ := os.Open(os.Args[1])
    reader := csv.NewReader(bufio.NewReader(csvFile))


    var Solar []DataPoint
    var Mains []DataPoint
	var AClamp1 []DataPoint
	var AClamp2 []DataPoint

	solarWH := 0.0
	mainsWH := 0.0
	AClamp1WH := 0.0
	AClamp2WH := 0.0

    for {
    	line, error := reader.Read()
    	if error == io.EOF {
	        break
	    } else if error != nil {
	        log.Fatal(error)
	    }
	    t, err := time.Parse("Mon Jan 02 15:04:05", string(line[0]))
	    if err != nil {
		    fmt.Println(err)
		}

		s := strings.Split(string(line[1]), " ")
		if s[0] == "Solar" {

			//add next interval of energy produced
			solarWH += producedInInterval(Solar, s[1], t)

			//add data point to array
			Solar = append(Solar, DataPoint{
				timestamp: t, 
				data: s[1],
			})

		} else if s[0] == "Mains" {

			mainsWH += producedInInterval(Mains, s[2], t)

			Mains = append(Mains, DataPoint{
				timestamp: t, 
				data: s[2],
			})
		} else {
			if len(s) == 4 && s[2] == "1"{

				AClamp1WH += producedInInterval(AClamp1, s[3], t)

				AClamp1 = append(AClamp1, DataPoint{
					timestamp: t, 
					data: s[3],
				})
			} else if len(s) == 4  && s[2] == "2"{
				//fmt.Println(s[3])

				AClamp2WH += producedInInterval(AClamp2, s[3], t)

				AClamp2 = append(AClamp2, DataPoint{
					timestamp: t, 
					data: s[3],
				})
			}
		}
    }
    //fmt.Println(Solar)
    //fmt.Println(Mains)
    //fmt.Println(AClamp1)
    //fmt.Println(AClamp2)
    fmt.Printf("SolarWH: %g \n", solarWH)
    fmt.Printf("MainsWH: %g \n", mainsWH)
    fmt.Printf("AClamp1WH: %g \n", AClamp1WH)
    fmt.Printf("AClamp2WH: %g \n", AClamp2WH)
}

func producedInInterval(arr []DataPoint, d string, t time.Time) float64 {
	if len(arr) > 0  {
		lastTime := arr[len(arr) - 1].timestamp
		diff := t.Sub(lastTime)
		diffHours := diff.Hours()
		i, err := strconv.ParseFloat(d[:len(d) - 2],10)
		if err != nil {
		    fmt.Println(err)
		}
		return diffHours * i
	} else {
		return 0.0
	}
}
