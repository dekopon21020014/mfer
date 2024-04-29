package main

import (
	"fmt"
	"log"
	"os"
)

type Mfer struct {
	Interval         Interval
	Sensitivity      Sensitivity
	Patient          Patient
	DataType         int /* 0 to 9*/
	Offset           int
	Null             int
	BlockLength      int
	NumChannel       int
	NumSequence      int
	Pointer          int
	Classification   int    /* 0: unknown, 1: standard 12, ... , */
	Ldn              int    /* 波形コード，stringの波形情報も取得する必要あり */
	Filter           string /* 広域，低域とか */
	IPD              int    /* 無条件間引きとか，線形補完とか stringの説明もある？*/
	Wave             int    /* 波形データどうやって保存しようかな */
	ByteOrder        int    /* 0: big endian, 1: little endian */
	Version          int    /* 3byte必要 */
	CharCode         string /* 文字コード, デフォルトはASCII */
	Zero             int
	Comment          string
	MchineInfo       string
	Compression      int /* 圧縮コード データ長　圧縮データ の3つがある */
	Preamble         string
	Event            int    /* イベントコード 開始時刻，　持続時間　イベント情報の4つがある */
	Val              int    /* 値コード，時刻ポイント，値の3つがある */
	Condition        int    /* 記録条件，説明コード1, 説明コード２，開始ポイント，接続時間，説明内容*/
	Error            int    /* *波形変換誤差 */
	Set              int    /* グループ定義 */
	ReferencePointer string /* URLなどによって参照するらしい */
	Signiture        int    /* 署名方と署名データがある */
	Time             Time
	Message          string
	UID              int
	Map              int /* ファイルをマップする？？ */
	End              int
}

type Interval struct {
	Unit     int /* 0: Hz, 1: second(time), 2: m(distance) */
	Exponent int
	Integer  int
}

type Sensitivity struct {
	Unit     int /* 0: Vold, 1: mmHg(torr), 2: pa, 3: cmH2O, 4: mmHg .... 22: cd */
	Exponent int
	Integer  int
}

/* 患者情報 */
type Patient struct {
	Id         int
	Name       string
	Age        int
	AgeInDays  int
	BirthYear  int
	BirthMonth int
	BirthDay   int
	Sex        int
}

type Time struct {
	Year     int
	Month    int
	Day      int
	Hour     int
	Minute   int
	Second   int
	MiliSec  int
	MicroSec int
}

func main() {
	path := "sample-data/ECG02.mwf"
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(bytes); {
		tag := bytes[i]
		i++
		if tag == 0x3f { // チャンネル属性は可変長
			i++
		}

		length := bytes[i]
		i++

		if tag == 0x1e {
			fmt.Printf("length: %02x\n", length)
			fmt.Print("data_length: ")
			for j := 0; j < 4; j++ {
				fmt.Printf("%02x ", bytes[i+j])
			}
			fmt.Println("")
			fmt.Println("len(bytes): ", len(bytes))
			fmt.Println("i: ", i+4)
			break
		}

		for j := i; j < i+int(length); j++ {
			break
		}

		i += int(length)
	}
}
