package main

import (
	"bytes"
	"encoding/csv"
	"io/ioutil"
	"log"
	"os"
)

func main() {

    if len(os.Args) < 2 {
	println("usage: csvparer filename.csv")
        os.Exit(2)
    }

    fileName := os.Args[1]
    file, err := os.Open(fileName)
    defer file.Close()
    if err != nil {
    	log.Printf("error opening file: %s, ", err.Error())
	}

	var jsonFile bytes.Buffer
	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
			log.Fatal( err )
		}

	jsonFile.WriteString("[\n")
	for i, _ := range lines {
		if i != 0 {

			jsonFile.WriteString("\t{\n")
			for ii, l := range lines[i] {
				if (ii+1) == len(lines[i]) {
					jsonFile.WriteString("\t\t\"" + lines[0][ii] + "\": " + "\"" + l + "\"\n")
				} else {
					jsonFile.WriteString("\t\t\"" + lines[0][ii] + "\": " + "\"" + l + "\",\n")
				}
			}
			if (i+1) == len(lines) {
				jsonFile.WriteString("\t}\n")
			} else {
				jsonFile.WriteString("\t}, \n")
			}
		}

	}
	jsonFile.WriteString("]\n")
	ioutil.WriteFile("out.json", jsonFile.Bytes(), 0644)
	println("saved output to out.json")
}
