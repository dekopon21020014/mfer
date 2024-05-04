package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	var path string
	/*
	if len(os.Args) < 1 {
		path = "sample-data/ECG01.mwf"
	} else {
		path = os.Args[0]
	}
	*/
	path = "sample-data/ECG01.mwf"
	
	bytes, err := os.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	mfer, err := parseMfer(bytes)
	if (err != nil) {
		log.Fatal(err)
		m, _ := json.MarshalIndent(mfer, "", "    ")
		fmt.Println(string(m))
		return 
	}
	m, _ := json.MarshalIndent(mfer, "", "    ")
	fmt.Println(string(m))
}
