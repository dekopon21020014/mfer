package main

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
)

// 拡張関連のハンドラー
func handlePreamble(mfer *Mfer, bytes []byte, length uint32) error {
	mfer.Extensions.Preamble = string(bytes)
	return nil
}

func handleEvent(mfer *Mfer, bytes []byte, length uint32) error {
	if len(bytes) < 10 {
		return errors.New("イベント情報のバイト数が不足")
	}
	eventCode, err := Binary2Uint16(mfer.Control.ByteOrder, bytes[0:2]...)
	if err != nil {
		return err
	}
	mfer.Extensions.Event.Code = eventCode

	begin, err := Binary2Uint32(mfer.Control.ByteOrder, bytes[2:6]...)
	if err != nil {
		return err
	}
	mfer.Extensions.Event.Begin = begin

	duration, err := Binary2Uint32(mfer.Control.ByteOrder, bytes[6:10]...)
	if err != nil {
		return err
	}
	mfer.Extensions.Event.Duration = duration
	mfer.Extensions.Event.Info = string(bytes[10:])
	return nil
}

func handleValue(mfer *Mfer, bytes []byte, length uint32) error {
	code, err := Binary2Uint16(mfer.Control.ByteOrder, bytes[0:2]...) //binary.BigEndian.Uint16(bytes[i : i+2])
	if err != nil {
		return err
	}

	time, err := Binary2Uint32(mfer.Control.ByteOrder, bytes[2:6]...) // binary.BigEndian.Uint32(bytes[i+2 : i+6])
	if err != nil {
		return err
	}

	val := string(bytes[6:int(length)])
	fmt.Sprintf("Value\n  code = %d, time = %d, val = %s\n", code, time, val)
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
