package main

import (
	"errors"
	"fmt"
)

// コントロール関連のハンドラー
func handleByteOrder(mfer *Mfer, bytes []byte, length uint32) error {
	if len(bytes) < 1 {
		return errors.New("バイトオーダーのバイト数が不足")
	}
	mfer.Control.ByteOrder = bytes[0]
	return nil
}

func handleVersion(mfer *Mfer, bytes []byte, length uint32) error {
	mfer.Control.Version = bytes
	return nil
}

func handleCharCode(mfer *Mfer, bytes []byte, length uint32) error {
	mfer.Control.CharCode = string(bytes)
	return nil
}

func handleZero(mfer *Mfer, bytes []byte, length uint32) error {
	return nil
}

func handleComment(mfer *Mfer, bytes []byte, length uint32) error {
	mfer.Control.Comment += string(bytes)
	return nil
}

func handleMachineInfo(mfer *Mfer, bytes []byte, length uint32) error {
	mfer.Control.MachineInfo = string(bytes)
	return nil
}

func handleCompression(mfer *Mfer, bytes []byte, length uint32) error {
	/*
		if len(bytes) < 6 {
			return errors.New("圧縮情報のバイト数が不足")
		}
	*/
	compressionCode, err := Binary2Uint16(0, bytes[0:2]...)
	if err != nil {
		return err
	}
	fmt.Printf("圧縮コード = %d\n", compressionCode)
	mfer.Control.CompressionCode = compressionCode

	dataLength, err := Binary2Uint32(0, bytes[2:6]...)
	if err != nil {
		return err
	}
	fmt.Printf("コード = %d, データ長 = %d\n", compressionCode, dataLength)
	return nil
}
