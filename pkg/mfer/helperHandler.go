package mfer

import (
	b "bytes"
	"errors"
	"io"
	"unicode/utf8"

	"github.com/dekopon21020014/mfer/pkg/byteorder"
	"github.com/fatih/color"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

var helperHandlerMap = map[byte]tagHandler{
	// Mfer.Helper
	P_NAME:  handlePatientName,
	P_ID:    handlePatientID,
	P_AGE:   handlePatientAge,
	P_SEX:   handlePatientSex,
	TIME:    handleTime,
	MESSAGE: handleMessage,
	UID:     handleUID,
	MAP:     handleMap,
}

// ヘルパー関連のハンドラー
func handlePatientName(mfer *Mfer, bytes []byte, length uint32) error {
	if utf8.Valid(bytes) {
		mfer.Helper.Patient.Name = string(bytes)
	} else {
		reader := transform.NewReader(b.NewReader(bytes), japanese.ShiftJIS.NewDecoder())
		utf8Bytes, err := io.ReadAll(reader)
		if err != nil {
			return err
		}
		mfer.Helper.Patient.Name = string(utf8Bytes)
	}
	return nil
}

func handlePatientID(mfer *Mfer, bytes []byte, length uint32) error {
	mfer.Helper.Patient.Id = string(bytes)
	return nil
}

func handlePatientAge(mfer *Mfer, bytes []byte, length uint32) error {
	if len(bytes) < 7 {
		return errors.New("患者の年齢情報のバイト数が不足")
	}
	mfer.Helper.Patient.Age = bytes[0]
	ageInDays, err := byteorder.Binary2Uint32(mfer.Control.ByteOrder, bytes[1:3]...)
	if err != nil {
		return err
	}
	mfer.Helper.Patient.AgeInDays = ageInDays

	birthYear, err := byteorder.Binary2Uint32(mfer.Control.ByteOrder, bytes[3:5]...)
	if err != nil {
		return err
	}
	mfer.Helper.Patient.BirthYear = birthYear
	mfer.Helper.Patient.BirthMonth = bytes[5]
	mfer.Helper.Patient.BirthDay = bytes[6]
	return nil
}

func handlePatientSex(mfer *Mfer, bytes []byte, length uint32) error {
	if len(bytes) < 1 {
		return errors.New("バイト数がおかしい")
	}
	return nil
}

func handleTime(mfer *Mfer, bytes []byte, length uint32) error {
	year, err := byteorder.Binary2Uint16(mfer.Control.ByteOrder, bytes[0:2]...)
	if err != nil {
		return err
	}
	month := bytes[2]
	day := bytes[3]
	hour := bytes[4]
	minute := bytes[5]
	second := bytes[6]
	miliSec, err := byteorder.Binary2Uint16(mfer.Control.ByteOrder, bytes[7:9]...)
	if err != nil {
		return err
	}

	microSec, err := byteorder.Binary2Uint16(mfer.Control.ByteOrder, bytes[9:11]...)
	if err != nil {
		return err
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
	return nil
}

func handleMessage(mfer *Mfer, bytes []byte, length uint32) error {
	color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")
	color.Green("MESSAGE: No Imprementation\n")
	color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")
	return nil
}

func handleUID(mfer *Mfer, bytes []byte, length uint32) error {
	color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")
	color.Green("UID: No Imprementation\n")
	color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")
	return nil
}
func handleMap(mfer *Mfer, bytes []byte, length uint32) error {
	color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")
	color.Green("MAP: No Imprementation\n")
	color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")
	return nil
}
