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
		path = "ECG/K-Heart/2024-11-06_17-29-22/0a5b84fc657983af4c80b3ae599b878f670140fd834ca9f27b63d8fbeef17f65_20241106170856.mwf"
	} else {
		path = os.Args[1]
	}

	mfer := newMfer()
	mfer, err = loadMfer(mfer, path)
	if err != nil {
		print("error *************************************************")
		log.Fatal(err)
		m, _ := json.MarshalIndent(mfer, "", "    ")
		fmt.Println(string(m))
		return
	}

	m, _ := json.MarshalIndent(mfer, "", "    ")
	fmt.Println(string(m))
	// fmt.Printf("%+v", mfer)
}
