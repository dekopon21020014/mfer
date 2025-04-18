package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	m "github.com/dekopon21020014/mfer/pkg/mfer"
	"github.com/dekopon21020014/mfer/pkg/mfer2physionet"
	"github.com/dekopon21020014/mfer/pkg/std12lead"
)

func main() {
	// コマンドラインオプション
	outputDir := flag.String("d", ".", "出力先ディレクトリ")
	parallel := flag.Int("p", 4, "並列処理の数")
	flag.Parse()

	// 入力パスを取得
	if flag.NArg() < 1 {
		fmt.Println("エラー: 入力ファイルまたはディレクトリを指定してください。")
		fmt.Println("使用方法: コマンド [オプション] <入力ファイルまたはディレクトリ>")
		fmt.Println("オプション:")
		flag.PrintDefaults()
		os.Exit(1)
	}
	inputPath := flag.Arg(0)

	// 出力先ディレクトリを作成
	if err := os.MkdirAll(*outputDir, os.ModePerm); err != nil {
		log.Fatalf("出力ディレクトリの作成に失敗しました: %v", err)
	}

	// 入力がディレクトリかファイルかを判定
	info, err := os.Stat(inputPath)
	if err != nil {
		log.Fatalf("指定された入力パスが無効です: %v", err)
	}

	if info.IsDir() {
		// ディレクトリ内の.mwfファイルを並列で処理
		processDirectory(inputPath, *outputDir, *parallel)
	} else {
		// 単一ファイルを処理
		processFile(inputPath, *outputDir)
	}
}

func processDirectory(inputDir, outputDir string, parallel int) {
	// .mwfファイルを収集
	var files []string
	err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".mwf") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("ディレクトリ内のファイル処理に失敗しました: %v", err)
	}

	// 並列処理用のワーカーグループ
	var wg sync.WaitGroup
	fileChan := make(chan string, len(files))

	// ファイルをチャネルに送信
	for _, file := range files {
		fileChan <- file
	}
	close(fileChan)

	// ワーカーを起動
	for i := 0; i < parallel; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range fileChan {
				processFile(file, outputDir)
			}
		}()
	}

	// 全てのワーカーが終了するのを待機
	wg.Wait()
}

func processFile(inputPath, outputDir string) {
	// MFERファイルをロード
	mfer := m.NewMfer()
	mfer, err := m.LoadMfer(mfer, inputPath)
	if err != nil {
		fmt.Printf("エラー: MFERファイルのロードに失敗しました: %s\n", inputPath)
		return
	}

	if len(mfer.Frames[0].WaveForm.Data) == 0 {
		fmt.Printf("エラー: MFERファイルにフレームがありません: %s\n", inputPath)
		return
	}

	// 12誘導に変換
	calculator, err := std12lead.NewLeadCalculator(&mfer)
	if err != nil {
		fmt.Printf("エラー: 12誘導への変換に失敗しました: %s\n", inputPath)
		return
	}
	leads, err := calculator.Convert8To12Lead()
	if err != nil {
		fmt.Printf("エラー: リード変換に失敗しました: %s\n", inputPath)
		return
	}

	// PhysioNet形式のデータに変換
	physionetData := mfer2physionet.Convert(leads)

	// 入力ファイルの拡張子を.datに変更して出力ファイル名を決定
	outputFileName := filepath.Base(inputPath)

	// datファイルの作成
	datFileName := outputFileName[:len(outputFileName)-len(filepath.Ext(outputFileName))] + ".dat"
	datFilePath := filepath.Join(outputDir, datFileName)
	datFile, err := os.Create(datFilePath)
	if err != nil {
		fmt.Printf("エラー: ファイル作成に失敗しました: %s\n", datFilePath)
		return
	}
	defer datFile.Close()

	// データを書き込む
	_, err = datFile.Write(physionetData)
	if err != nil {
		fmt.Printf("エラー: データ書き込みに失敗しました: %s\n", datFilePath)
		return
	}

	heaFileName := outputFileName[:len(outputFileName)-len(filepath.Ext(outputFileName))] + ".hea"
	heaFilePath := filepath.Join(outputDir, heaFileName)
	heaFile, err := os.Create(heaFilePath)
	if err != nil {
		fmt.Printf("エラー: .heaファイルの作成に失敗しました: %s\n", heaFilePath)
		return
	}
	defer heaFile.Close()

	if err := mfer2physionet.MakeHeaFile(heaFile, physionetData, &mfer); err != nil {
		fmt.Printf("エラー: .heaファイルの作成に失敗しました: %s\n", err)
		return
	}
}
