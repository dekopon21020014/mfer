package main

import (
	b "bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"os"
	"unicode/utf8"
	"github.com/fatih/color"
)

func loadMfer(mfer Mfer, path string) (Mfer, error) {
	var (
		tagCode byte
		length  uint32
	)

	bytes, err := os.ReadFile(path)
	if err != nil { /* cannnot open the file from path */
		return mfer, err
	}

	for i := 0; i < len(bytes); {
		tagCode = bytes[i]
		i++
		if tagCode == ZERO {
			continue
		} else if tagCode == END {
			break
		}

		length = uint32(bytes[i])
		i++
		if length > 0x7f { /* MSBが1ならば */
			numBytes := length - 0x80
			if numBytes > 4 {
				fmt.Printf("length = %d, numBytes = %d, bytes = %d\n", length, numBytes, bytes[i-1])
				return mfer, errors.New("error nbytes")
			}
			length = binary.BigEndian.Uint32(append(make([]byte, 4-numBytes), bytes[i:i+int(numBytes)]...))
			i += int(numBytes)
		}

		// fmt.Printf("tagCode = %02x, length = %02x(%d), i = %d\n", tagCode, length, length, i)

		switch tagCode {
		/*
		 * for Mfer.Sampling
		 */
		case INTERVAL: /* サンプリング間隔とか */
			mfer.Sampling.Interval.UnitCode = bytes[i]
			mfer.Sampling.Interval.Exponent = int8(bytes[i+1])
			mantissa, err := Binary2Uint32(mfer.Control.ByteOrder, bytes[i+2:i+int(length)]...)
			if err != nil {
				return mfer, err
			}
			mfer.Sampling.Interval.Mantissa = mantissa

		case SENSITIVITY: /* サンプリング解像度 */
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

		/*
		 * for Mfer.Frame
		 */
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
			fmt.Printf("  num chan = %d\n", mfer.Frame.NumChannel)

		case SEQUENCE: /* シーケンス数 */
			binaryNumSequence := bytes[i : i+int(length)]
			num, err := Binary2Uint32(mfer.Control.ByteOrder, binaryNumSequence...)
			if err != nil {
				return mfer, err
			}
			mfer.Frame.NumSequence = num

		case F_POINTER:
			/* do somthing */

		/*
		 * for Mfer.WaveFrom
		 */
		case WAVE_FORM_TYPE:
			code, err := Binary2Uint16(mfer.Control.ByteOrder, bytes[i:i+2]...)
			if err != nil {
				return mfer, err
			}
			mfer.WaveForm.Code = code
			mfer.WaveForm.Description = string(bytes[i+2 : i+int(length)])

		case CHANNEL_ATTRIBUTE:
			channelNumber := int(bytes[i-1])
			length = uint32(bytes[i])
			i++
			mfer.Frame.Channels[channelNumber] = Channel{
				Num:     channelNumber,
				TagCode: bytes[i],
				Length:  bytes[i+1],
				Data:    (bytes[i+2 : ((i + 2) + int(length))]),
			}

		/*
		 * for mfer.WaveForm
		 */
		case LDN:
			// 2byte 以上のこともある？
			mfer.WaveForm.Ldn = bytes[i]

		case INFORMATION:
			color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")
			color.Green("INFORMATION: No Imprementation\n")
			color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")
			// 何かよくわかっていない
			mfer.WaveForm.Information = "there are some information"

		case FILTER:
			mfer.WaveForm.Filter = string(bytes[i : i+int(length)])

		case IPD:
			// 2byte以上のこともあり？
			mfer.WaveForm.IpdCode = bytes[i]

		case DATA: /* 波形データ部 */
			color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")
			color.Green("DATA: No Imprementation\n")
			color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")

		/*
		 * for Mfer.Control
		 */
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
			compressionCode, err := Binary2Uint16(0, bytes[i:i+2]...)
			if err != nil {
				return mfer, err
			}
			mfer.Control.CompressionCode = compressionCode

			dataLength, err := Binary2Uint32(0, bytes[i+2:i+6]...)
			if err != nil {
				return mfer, err
			}
			fmt.Printf("code = %d, length = %d, dataLength = %d\n", compressionCode, length, dataLength)

		/*
		 * for Mfer.Extensions
		 */
		case PREAMBLE: /* プリアンブル */
			mfer.Extensions.Preamble = string(bytes[i : i+int(length)])

		case EVENT:
			eventCode, err := Binary2Uint16(mfer.Control.ByteOrder, bytes[i:i+2]...) //bytes[i : i+2]
			if err != nil {
				return mfer, err
			}
			mfer.Extensions.Event.Code = eventCode

			begin, err := Binary2Uint32(mfer.Control.ByteOrder, bytes[i+2:i+6]...)
			if err != nil {
				return mfer, err
			}
			mfer.Extensions.Event.Begin = begin

			duration, err := Binary2Uint32(mfer.Control.ByteOrder, bytes[i+6:i+10]...)
			if err != nil {
				return mfer, err
			}
			mfer.Extensions.Event.Duration = duration
			mfer.Extensions.Event.Info = string(bytes[i+10 : i+int(length)])

		case VALUE: /* 測定値など波形の値に関する情報を記述する */
			code, err := Binary2Uint16(mfer.Control.ByteOrder, bytes[i:i+2]...) //binary.BigEndian.Uint16(bytes[i : i+2])
			if err != nil {
				return mfer, err
			}

			time, err := Binary2Uint32(mfer.Control.ByteOrder, bytes[i+2:i+6]...) // binary.BigEndian.Uint32(bytes[i+2 : i+6])
			if err != nil {
				return mfer, err
			}

			val := string(bytes[i+6 : i+int(length)])
			fmt.Sprintf("Value\n  code = %d, time = %d, val = %s\n", code, time, val)

		case CONDITION:
			color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")
			color.Green("CONDITION: No Imprementation\n")
			color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")

		case ERROR:
			color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")
			color.Green("ERROR: No Imprementation\n")
			color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")

		case GROUP:
			color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")
			color.Green("GROUP: No Imprementation\n")
			color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")

		case R_POINTER:
			color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")
			color.Green("R_POINTER: No Imprementation\n")
			color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")

		case SIGNITURE:
			color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")
			color.Green("SIGNITURE: No Imprementation\n")
			color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")

		/*
		 * for Mfer.Helper
		 */
		case P_NAME: /* 患者の名前 */
			if utf8.Valid(bytes[i : i+int(length)]) {
				mfer.Helper.Patient.Name = string(bytes[i : i+int(length)])
			} else {
				reader := transform.NewReader(b.NewReader(bytes[i:i+int(length)]), japanese.ShiftJIS.NewDecoder())
				utf8Bytes, err := ioutil.ReadAll(reader)
				if err != nil {
					log.Fatal(err)
				}
				mfer.Helper.Patient.Name = string(utf8Bytes)
			}

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
			year, err := Binary2Uint16(mfer.Control.ByteOrder, bytes[i : i+2]...)
			if err != nil {
				return mfer, err
			}
			month := bytes[i+2]
			day := bytes[i+3]
			hour := bytes[i+4]
			minute := bytes[i+5]
			second := bytes[i+6]
			miliSec, err := Binary2Uint16(mfer.Control.ByteOrder, bytes[i+7 : i+9]...)
			if err != nil {
				return mfer, err
			}

			microSec, err := Binary2Uint16(mfer.Control.ByteOrder, bytes[i+9 : i+11]...)
			if err != nil {
				return mfer, err
			}

			time := Time{
				Year:     year,
				Month:    month,
				Day:      day,
				Hour:     hour,
				Minute:   minute,
				Second:   second,
				MiliSec:  miliSec,
				MicroSec: microSec,
			}
			mfer.Helper.Time = time

		case MESSAGE:
			color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")
			color.Green("MESSAGE: No Imprementation\n")
			color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")

		case UID:
			color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")
			color.Green("UID: No Imprementation\n")
			color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")

		case MAP:
			color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")
			color.Green("MAP: No Imprementation\n")
			color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")

		}

		i += int(length)
	}
	return mfer, nil
}
