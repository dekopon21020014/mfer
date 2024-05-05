package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"log"
)

func loadMfer(mfer Mfer, path string) (Mfer, error) {
	var (
		length  byte
		tagCode byte
	)

	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

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
			mantissa, err := Binary2Uint32(mfer.Control.ByteOrder, bytes[i+2:i+int(length)]...)
			if err != nil {
				return mfer, err
			}
			mfer.Sampling.Interval.Mantissa = mantissa

		case SENSITIVITY: /* サンプリング解像度 */
			// 得られるのはコードだから文字列に変換したいなら別途処理が必要
			mfer.Sampling.Sensitivity.UnitCode = bytes[i]
			mfer.Sampling.Sensitivity.Exponent = int8(bytes[i+1])
			mantissa, err := Binary2Uint32(mfer.Control.ByteOrder, bytes[i+2:i+int(length)]...)
			if err != nil {
				return mfer, err
			}
			mfer.Sampling.Sensitivity.Mantissa = mantissa

		case DATA_TYPE:
			mfer.Sampling.DataTypeCode = bytes[i]

		case OFFSET:
			offset, err := Binary2Uint64(mfer.Control.ByteOrder, bytes[i:i+int(length)]...)
			if err != nil {
				return mfer, err
			}
			mfer.Sampling.Offset = offset

		case NULL:
			mfer.Sampling.Null = 0

		// for Mfer.Frame
		case BLOCK: /* ブロック長 1ブロックが何データから成るか */
			binaryBlockLength := bytes[i : i+int(length)]
			blockLength, err := Binary2Uint32(mfer.Control.ByteOrder, binaryBlockLength...)
			if err != nil {
				return mfer, err
			}
			mfer.Frame.BlockLength = blockLength

		case CHANNEL: /* チャンネル数 */
			binaryNumChannel := bytes[i : i+int(length)]
			num, err := Binary2Uint32(mfer.Control.ByteOrder, binaryNumChannel...)
			if err != nil {
				return mfer, err
			}
			mfer.Frame.NumChannel = num
			fmt.Printf("  chan = %d\n", mfer.Frame.NumChannel)

		case SEQUENCE: /* シーケンス数 */
			binaryNumSequence := bytes[i : i+int(length)]
			num, err := Binary2Uint32(mfer.Control.ByteOrder, binaryNumSequence...)
			if err != nil {
				return mfer, err
			}
			mfer.Frame.NumSequence = num

		case F_POINTER:
			/* do somthing */

		// for Mfer.WaveFrom
		case WAVE_FORM_TYPE:
			code, err := Binary2Uint16(mfer.Control.ByteOrder, bytes[i:i+int(length)]...)
			if err != nil {
				return mfer, err
			}
			mfer.WaveForm.Code = code
			// mfer.WaveForm.Code, err = Binary2Uint16(mfer.Control.ByteOrder, bytes[i:i+int(length)]...)

		case CHANNEL_ATTRIBUTE:
			fmt.Printf("  chan = %d\n", len(mfer.Frame.Channels))
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
			// tagやデータ長は無条件でビッグエンディアン
			// dataLength := Binary2Uint32(0, bytes[i : i+4]...)
			dataLength := binary.BigEndian.Uint32(bytes[i : i+4])
			// mfer.WaveForm.Data = bytes[i : i + int(dataLength)] // この行は必要だけど見にくく成るのでコメントアウトしておく
			i += 4
			i += int(dataLength)
			continue

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
			ageInDays, err := Binary2Uint32(mfer.Control.ByteOrder, bytes[i+1:i+3]...)
			if err != nil {
				return mfer, err
			}
			mfer.Helper.Patient.AgeInDays = ageInDays

			birthYear, err := Binary2Uint32(mfer.Control.ByteOrder, bytes[i+3:i+5]...)
			if err != nil {
				return mfer, err
			}
			mfer.Helper.Patient.BirthYear = birthYear
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
	return mfer, nil
}
