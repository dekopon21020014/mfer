package main

import (
	"encoding/binary"
	"fmt"
	// "encoding/json"
)

func parseMfer(bytes []byte) Mfer {
	var (
		mfer    Mfer
		length  byte
		tagCode byte
	)

	for i := 0; i < len(bytes); {
		tagCode = bytes[i]
		i++
		length = bytes[i]
		i++
		// fmt.Printf("tagCode = %02x, length = %02x(%d), i = %d\n", tagCode, length, length, i)

		switch tagCode { /* tag名を出力する */
		// for Mfer.Sampling
		case INTERVAL: /* サンプリング間隔とか */
			mfer.Sampling.Interval.UnitCode = bytes[i]
			mfer.Sampling.Interval.Exponent = int8(bytes[i+1])
			binaryMantissa := bytes[i+2 : i+int(length)]
			mfer.Sampling.Interval.Mantissa = Binary2Uint32(mfer.Control.ByteOrder, binaryMantissa...)

		case SENSITIVITY: /* サンプリング解像度 */
			// 得られるのはコードだから文字列に変換したいなら別途処理が必要
			mfer.Sampling.Sensitivity.UnitCode = bytes[i]
			mfer.Sampling.Sensitivity.Exponent = int8(bytes[i+1])
			mfer.Sampling.Sensitivity.Mantissa = Binary2Uint32(mfer.Control.ByteOrder, bytes[i+2:i+int(length)]...)

		case DATA_TYPE:
			mfer.Sampling.DataTypeCode = bytes[i]

		case OFFSET:
			mfer.Sampling.Offset = Binary2Uint64(mfer.Control.ByteOrder, bytes[i:i+int(length)]...)

		case NULL:
			mfer.Sampling.Null = 0

		// for Mfer.Frame
		case BLOCK: /* ブロック長 1ブロックが何データから成るか */
			binaryBlockLength := bytes[i : i+int(length)]
			mfer.Frame.BlockLength = Binary2Uint32(mfer.Control.ByteOrder, binaryBlockLength...)

		case CHANNEL: /* チャンネル数 */
			binaryNumChannel := bytes[i : i+int(length)]
			mfer.Frame.NumChannel = Binary2Uint32(mfer.Control.ByteOrder, binaryNumChannel...)

		case SEQUENCE: /* シーケンス数 */
			binaryNumSequence := bytes[i : i+int(length)]
			mfer.Frame.NumSequence = Binary2Uint32(mfer.Control.ByteOrder, binaryNumSequence...)

		case POINTER:
			/* do somthing */

		// for Mfer.WaveFrom
		case WAVE_FORM_TYPE:
			mfer.WaveForm.Code = Binary2Uint16(mfer.Control.ByteOrder, bytes[i:i+int(length)]...)

		case CHANNEL_ATTRIBUTE:
			length = bytes[i] // チャンネル属性のタグコードは2バイト
			i++
			channnel := Channel{
				Num:     len(mfer.Frame.Channels),
				TagCode: bytes[i],
				Length:  bytes[i+1],
				Data:    (bytes[i+2 : (i + 2 + int(bytes[i+1]))]),
			}
			mfer.Frame.Channels = append(mfer.Frame.Channels, channnel)
		
		// for mfer.WaveForm
		case LDN:
			// 2byte 以上のこともある？
			mfer.WaveForm.Ldn = bytes[i]

		case INFORMATION:
			// 何かよくわかっていない
			mfer.WaveForm.Information = "there are some information"

		case FILTER:
			mfer.WaveForm.Filter = string(bytes[i : i+int(length)])

		case IDP:
			// 2byte以上のこともあり？
			mfer.WaveForm.IpdCode = bytes[i]

		case 0x17: /* 波形データの制作者 */
			mfer.Control.MchineInfo = string(bytes[i : i+int(length)])

		case 0x01: /* バイトオーダー */
			mfer.Control.ByteOrder = bytes[i]

		case DATA: /* 波形データ部 */
			fmt.Printf("length = %x\n", length)
			// tagやデータ長は無条件でビッグエンディアん
			dataLength := Binary2Uint32(0, bytes[i : i+4]...)
			fmt.Printf("DataLength = %d\n", dataLength)
			i = i+4
			// mfer.WaveForm.Data = bytes[i : i + int(dataLength)] // この行は必要だけど見にくく成るのでコメントアウトしておく

		case PREAMBLE: /* プリアンブル */
			mfer.Extensions.Preamble = string(bytes[i : i+int(length)])
		}

		if tagCode == DATA { /* 波形データ部まできたらループ終了 */
			break
		}

		i += int(length)
	}
	return mfer
}

func Binary2Uint16(byteOrder byte, bytes ...byte) uint16 {
	// fmt.Printf("wave code = %d\n", bytes)
	b := make([]byte, len(bytes))
	copy(b, bytes)
	padding := make([]byte, 2-len(b))
	if byteOrder == 0 { // big endian
		b = append(padding, b...)
		return binary.BigEndian.Uint16(b)
	} else if byteOrder == 1 { // little endian
		b = append(b, padding...)
		return binary.LittleEndian.Uint16(b)
	}
	return 111
}

// 4byte以内の文字列をuint32に変換する
func Binary2Uint32(byteOrder byte, bytes ...byte) uint32 {
	b := make([]byte, len(bytes))
	copy(b, bytes)
	padding := make([]byte, 4-len(b))
	if byteOrder == 0 { // big endian
		b = append(padding, b...)
		return binary.BigEndian.Uint32(b)
	} else if byteOrder == 1 { // little endian
		b = append(b, padding...)
		return binary.LittleEndian.Uint32(b)
	}
	return 0
}

func Binary2Uint64(byteOrder byte, bytes ...byte) uint64 {
	b := make([]byte, len(bytes))
	copy(b, bytes)
	padding := make([]byte, 8-len(b))
	if byteOrder == 0 { // big endian
		b = append(padding, b...)
		return binary.BigEndian.Uint64(b)
	} else if byteOrder == 1 { // little endian
		b = append(b, padding...)
		return binary.LittleEndian.Uint64(b)
	}
	return 0
}
