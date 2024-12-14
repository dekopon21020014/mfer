package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	m "github.com/dekopon21020014/mfer/pkg/mfer"
	"github.com/dekopon21020014/mfer/pkg/mfer2physionet"
	"github.com/dekopon21020014/mfer/pkg/std12lead"
)

func main() {
	var path string
	var err error

	if len(os.Args) < 2 { // コマンドライン引数なしならとりあえずECG01を対象にする(開発用)
		path = "ECG/K-Heart/2024-11-06_17-29-22/0a5b84fc657983af4c80b3ae599b878f670140fd834ca9f27b63d8fbeef17f65_20241106170856.mwf"
	} else {
		path = os.Args[1]
	}

	mfer := m.NewMfer()
	mfer, err = m.LoadMfer(mfer, path)
	if err != nil {
		print("********************* ERROR **********************")
		log.Fatal(err)
		m, _ := json.MarshalIndent(mfer, "", "    ")
		fmt.Println(string(m))
		return
	}

	m, _ := json.MarshalIndent(mfer, "", "    ")
	fmt.Println(string(m))

	calculator, err := std12lead.NewLeadCalculator(&mfer)
	if err != nil {
		fmt.Println(err)
		return
	}

	leads, err := calculator.Convert8To12Lead()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v", leads)

	physionetData := mfer2physionet.Convert(leads)

	file, err := os.Create("tmp/hoge.dat")
	if err != nil {
		log.Fatalf("ファイル作成に失敗しました: %v", err)
	}
	defer file.Close() // 処理が終わったらファイルを閉じる

	// データを書き込む
	_, err = file.Write(physionetData)
	if err != nil {
		log.Fatalf("データ書き込みに失敗しました: %v", err)
	}
}
