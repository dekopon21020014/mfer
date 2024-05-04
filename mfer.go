package main

type Mfer struct {
	Sampling   Sampling
	Frame      Frame
	WaveForm   WaveForm
	Control    Control
	Extensions Extensions
	Helper     Helper
}

type Sampling struct {
	Interval     Interval
	Sensitivity  Sensitivity
	DataTypeCode byte /* 0 to 9 */
	Offset       uint64
	Null         int
}

type Interval struct {
	UnitCode byte /* 0: Hz, 1: second(time), 2: m(distance) */
	Exponent int8
	Mantissa uint32
}

type Sensitivity struct {
	UnitCode byte /* 0: Vold, 1: mmHg(torr), 2: pa, 3: cmH2O, 4: mmHg .... 22: cd */
	Exponent int8
	Mantissa uint32
}

type Frame struct {
	BlockLength uint32
	NumChannel  uint32
	NumSequence uint32
	Pointer     int
	Channels    []Channel
}

type Channel struct {
	Num     int
	TagCode byte
	Length  byte
	//Ldn     int
	Data []byte
}

type WaveForm struct {
	Code uint16 /* 0: unknown, 1: standard 12, ... , */
	//Attribute   Attribute
	Ldn         byte /* 波形コード，stringの波形情報も取得する必要あり */
	Information string
	Filter      string /* 広域，低域とか */
	IpdCode     byte   /* 無条件間引きとか，線形補完とか stringの説明もある？ */
	Data        []byte /* 波形データどうやって保存しようかな */
}

type Attribute struct {
	length int
}

type Control struct {
	ByteOrder       byte   /* 0: big endian, 1: little endian */
	Version         []byte /* 3byte必要 */
	CharCode        string /* 文字コード, デフォルトはASCII */
	Zero            int
	Comment         string
	MachineInfo     string
	CompressionCode uint16 /* 圧縮コード データ長　圧縮データ の3つがある */
}

type Extensions struct {
	Preamble string
	Event    Event /* イベントコード 開始時刻，　持続時間　イベント情報の4つがある */

	Value            int    /* 値コード，時刻ポイント，値の3つがある */
	Condition        int    /* 記録条件，説明コード1, 説明コード２，開始ポイント，接続時間，説明内容 */
	Error            int    /* *波形変換誤差 */
	Group            int    /* グループ定義 */
	ReferencePointer string /* URLなどによって参照するらしい */
	Signiture        int    /* 署名方法と署名データがある */
}

type Event struct {
	Code     byte
	Begin    uint32
	Duration uint32
	Info     string
}

type Helper struct {
	Patient Patient
	Time    Time
	Message string
	UID     int
	Map     int /* ファイルをマップする？？ */
	End     int
}

/* 患者情報 */
type Patient struct {
	Id         string
	Name       string
	Age        byte
	AgeInDays  uint32
	BirthYear  uint32
	BirthMonth byte
	BirthDay   byte
	Sex        byte
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
