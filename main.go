package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	var path string
	var err error

	if len(os.Args) < 2 { // コマンドライン引数なしならとりあえずECG01を対象にする(開発用)
		path = "sample-data/ECG01.mwf"
	} else {
		path = os.Args[1]
	}

	mfer := newMfer()
	mfer, err = loadMfer(mfer, path)
	if err != nil {
		log.Fatal(err)
		m, _ := json.MarshalIndent(mfer, "", "    ")
		fmt.Println(string(m))
		return
	}
	
	m, _ := json.MarshalIndent(mfer, "", "    ")
	fmt.Println(string(m))
	// fmt.Printf("%+v", mfer)
}
