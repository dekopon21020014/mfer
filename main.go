package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	path := "sample-data/ECG01.mwf"
	bytes, err := os.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	mfer := parseMfer(bytes)
	m, _ := json.MarshalIndent(mfer, "", "    ")
	fmt.Println(string(m))
}
