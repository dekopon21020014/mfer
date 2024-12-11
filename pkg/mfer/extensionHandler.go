package mfer

import (
	"errors"
	"fmt"

	"github.com/dekopon21020014/mfer/pkg/byteorder"
	"github.com/fatih/color"
)

var extensionHandlerMap = map[byte]tagHandler{
	// Mfer.Extensions
	PREAMBLE:  handlePreamble,
	EVENT:     handleEvent,
	VALUE:     handleValue,
	CONDITION: handleCondition,
	ERROR:     handleError,
	GROUP:     handleGroup,
	R_POINTER: handleRPointer,
	SIGNITURE: handleSignature,
}

// 拡張関連のハンドラー
func handlePreamble(mfer *Mfer, bytes []byte, length uint32) error {
	mfer.Extensions.Preamble = string(bytes)
	return nil
}

func handleEvent(mfer *Mfer, bytes []byte, length uint32) error {
	if len(bytes) < 10 {
		return errors.New("イベント情報のバイト数が不足")
	}
	eventCode, err := byteorder.Binary2Uint16(mfer.Control.ByteOrder, bytes[0:2]...)
	if err != nil {
		return err
	}
	mfer.Extensions.Event.Code = eventCode

	begin, err := byteorder.Binary2Uint32(mfer.Control.ByteOrder, bytes[2:6]...)
	if err != nil {
		return err
	}
	mfer.Extensions.Event.Begin = begin

	duration, err := byteorder.Binary2Uint32(mfer.Control.ByteOrder, bytes[6:10]...)
	if err != nil {
		return err
	}
	mfer.Extensions.Event.Duration = duration
	mfer.Extensions.Event.Info = string(bytes[10:])
	return nil
}

func handleValue(mfer *Mfer, bytes []byte, length uint32) error {
	code, err := byteorder.Binary2Uint16(mfer.Control.ByteOrder, bytes[0:2]...) //binary.BigEndian.Uint16(bytes[i : i+2])
	if err != nil {
		return err
	}

	time, err := byteorder.Binary2Uint32(mfer.Control.ByteOrder, bytes[2:6]...) // binary.BigEndian.Uint32(bytes[i+2 : i+6])
	if err != nil {
		return err
	}

	val := string(bytes[6:int(length)])
	_ = fmt.Sprintf("Value\n  code = %d, time = %d, val = %s\n", code, time, val)
	return nil
}

func handleCondition(mfer *Mfer, bytes []byte, length uint32) error {
	color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")
	color.Green("CONDITION: No Imprementation\n")
	color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")
	return nil
}

func handleError(mfer *Mfer, bytes []byte, length uint32) error {
	color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")
	color.Green("ERROR: No Imprementation\n")
	color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")
	return nil
}

func handleGroup(mfer *Mfer, bytes []byte, length uint32) error {
	color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")
	color.Green("GROUP: No Imprementation\n")
	color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")
	return nil
}

func handleRPointer(mfer *Mfer, bytes []byte, length uint32) error {
	color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")
	color.Green("R_POINTER: No Imprementation\n")
	color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")
	return nil
}

func handleSignature(mfer *Mfer, bytes []byte, length uint32) error {
	color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")
	color.Green("SIGNITURE: No Imprementation\n")
	color.Green("!!!!!!!!!!!!!!!!!!!!!!!\n")
	return nil
}
