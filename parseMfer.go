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
		fmt.Printf("tagCode = %02x, length = %02x(%d), i = %d\n", tagCode, length, length, i)

		switch tagCode { /* tag名を出力する */
		// for Mfer.Sampling
		case INTERVAL: /* サンプリング間隔とか */
			mfer.Sampling.Interval.UnitCode = bytes[i]
			mfer.Sampling.Interval.Exponent = int8(bytes[i+1])
			mfer.Sampling.Interval.Mantissa = Binary2Uint32(mfer.Control.ByteOrder, bytes[i+2:i+int(length)]...)

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

		case F_POINTER:
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

		case DATA: /* 波形データ部 */
			// tagやデータ長は無条件でビッグエンディアん
			// dataLength := Binary2Uint32(0, bytes[i : i+4]...)
			i = i + 4
			// mfer.WaveForm.Data = bytes[i : i + int(dataLength)] // この行は必要だけど見にくく成るのでコメントアウトしておく

		// for Mfer.Control
		case BYTE_ORDER:
			mfer.Control.ByteOrder = bytes[i]

		case VERSION:
			// なんだこれ？
			mfer.Control.Version = bytes[i : i+3]

		case CHAR_CODE:
			mfer.Control.CharCode = string(bytes[i : i+int(length)])

		case ZERO:
			i--

		case COMMENT:
			mfer.Control.Comment += string(bytes[i : i+int(length)])

		case MACHINE_INFO:
			mfer.Control.MachineInfo = string(bytes[i : i+int(length)])

		case COMPRESSION:
			fmt.Printf("COMPRESSTION\n")
			// ここもいまいちよくわからない
			i--
			mfer.Control.CompressionCode = binary.BigEndian.Uint16(bytes[i : i+2])
			// length = binary.BigEndian.Uint32(bytes[i+2 : i+6])

		// for Mfer.Extensions
		case PREAMBLE: /* プリアンブル */
			mfer.Extensions.Preamble = string(bytes[i : i+int(length)])

		case EVENT:
			/* よくわかっていない */
			i--
			mfer.Extensions.Event.Code = bytes[i]
			// ここはビッグエンディアンなのかリトルエンディアンなのか...
			mfer.Extensions.Event.Begin = binary.BigEndian.Uint32(bytes[i+2 : i+6])
			mfer.Extensions.Event.Duration = binary.BigEndian.Uint32(bytes[i+6 : i+10])
			mfer.Extensions.Event.Info = string(bytes[i+10 : i+10+256])

		// この辺のフォーマットがよくわかっていない
		case VALUE:
		case CONDITION:
		case ERROR:
		case GROUP:
		case R_POINTER:
		case SIGNITURE:

		// for Mfer.Helper
		case P_NAME:
			mfer.Helper.Patient.Name = string(bytes[i : i+int(length)])

		case P_ID:
			mfer.Helper.Patient.Id = string(bytes[i : i+int(length)])

		case P_AGE:
			mfer.Helper.Patient.Age = bytes[i]
			mfer.Helper.Patient.AgeInDays = Binary2Uint32(mfer.Control.ByteOrder, bytes[i+1 : i+3]...)
			mfer.Helper.Patient.BirthYear = Binary2Uint32(mfer.Control.ByteOrder, bytes[i+3 : i+5]...)
			mfer.Helper.Patient.BirthMonth = bytes[i+5]
			mfer.Helper.Patient.BirthDay = bytes[i+6]

		case P_SEX:
			mfer.Helper.Patient.Sex = bytes[i]

		case TIME:
			// よくわかっていない

		case MESSAGE:
		case UID:

		case MAP:
		case END:
			fmt.Printf("END\n")
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
