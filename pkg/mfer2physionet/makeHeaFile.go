package mfer2physionet

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	m "github.com/dekopon21020014/mfer/pkg/mfer"
)

func MakeHeaFile(file *os.File, physionetData []byte, mfer *m.Mfer) error {
	var heaContent string // こいつにheaファイルの中身を作る
	var fileName string
	var frequency, interval float64

	// 1. ファイルのベースネームから拡張子を除去したものを取得
	fileName = filepath.Base(file.Name())
	ext := filepath.Ext(fileName)
	fileName = strings.TrimSuffix(fileName, ext)

	// 2. 周波数を求める
	if mfer.Sampling.Interval.UnitCode != 1 {
		// 基本的に1のはずなんだけど，1じゃなかった時に気づけるように強制終了するように実装しておく
		fmt.Printf("Sampling Interval isn't 1\n")
		log.Fatalf("Sampling Interval isn't 1")
		return fmt.Errorf("Sampling Interval isn't 1")
	} else if mfer.Sampling.Interval.UnitCode == 1 {
		// これが1なら，記述されているデータの単位は秒
		mantissa := mfer.Sampling.Interval.Mantissa // 仮数部
		exponent := mfer.Sampling.Interval.Exponent // 指数部
		interval = float64(mantissa) * math.Pow(10, float64(exponent))

		// Hz を計算
		frequency = 1 / interval // これはi整数に直した方がよいのか・
	}

	// 3. 計測時間
	blockLength := mfer.Frames[0].BlockLength
	numSequence := mfer.Frames[0].NumSequence
	// timeLnegth = interval * float64(blockLength) * float64(numSequence) * 1000
	numSample := blockLength * numSequence

	// .heaの1行目
	heaContent = fmt.Sprintf("%s 12 %d %d\n", fileName, uint32(frequency), numSample)
	fileName = fileName + ".dat"
	// 各誘導に関するデータ

	var sensitivity float64
	if mfer.Sampling.Sensitivity.UnitCode != 0 {
		fmt.Printf("sampling unit isn's Volt")
		log.Fatal("sampling unit isnt's Volt")
		return fmt.Errorf("sampling unit isn't Volt")
	} else if mfer.Sampling.Sensitivity.UnitCode == 0 {
		exponent := mfer.Sampling.Sensitivity.Exponent
		mantissa := mfer.Sampling.Sensitivity.Mantissa
		sensitivity = float64(mantissa) * math.Pow(10, float64(exponent)) * 1000
	}

	var numBit uint8
	var maxVal, minVal int16
	if mfer.Sampling.DataTypeCode != 0 {
		fmt.Println("data type isn't signed 16 bit integer")
		log.Fatal("data type isn't signed 16 bit integer")
	} else if mfer.Sampling.DataTypeCode == 0 {
		numBit = 16
		dataLen := len(physionetData) / 2
		int16Data := make([]int16, dataLen)
		for i := 0; i < dataLen; i++ {
			int16Data[i] = int16(binary.LittleEndian.Uint16(physionetData[i*2 : i*2+2]))
		}
		maxVal, minVal = math.MinInt16, math.MaxInt16
		for _, v := range int16Data {
			if v > maxVal {
				maxVal = v
			}
			if v < minVal {
				minVal = v
			}
		}
	}

	var offset uint64
	if mfer.Sampling.Offset != 0 {
		fmt.Println("Offset isn't 0")
		log.Fatal("offset isn't 0")
	} else {
		offset = 0
	}

	for _, lead := range []string{
		"I", "II", "III",
		"AVR", "AVL", "AVF",
		"V1", "V2", "V3", "V4", "V5", "V6",
	} {
		heaContent += fmt.Sprintf(
			"%s %d %f(%d)/mV %d %d %d %d %d %s\n",
			fileName, numBit, sensitivity, offset, numBit, offset, minVal, maxVal, offset, lead,
		)
	}

	if _, err := file.Write([]byte(heaContent)); err != nil {
		log.Fatal("heaファイルの書き込みに失敗")
	}

	return nil
}
