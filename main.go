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
	fmt.Printf("beginning of len(bytes) = %d\n", len(bytes))

	if err != nil {
		log.Fatal(err)
	}

	mfer := parseMfer(bytes)
	m, _ := json.MarshalIndent(mfer, "", "    ")
	fmt.Println(string(m))
}
// 80 38 01 00