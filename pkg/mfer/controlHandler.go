package mfer

import (
	"errors"
	"fmt"

	"github.com/dekopon21020014/mfer/pkg/byteorder"
)

var controlHandlerMap = map[byte]tagHandler{
	// Mfer.Control
	BYTE_ORDER:   handleByteOrder,
	VERSION:      handleVersion,
	CHAR_CODE:    handleCharCode,
	ZERO:         handleZero,
	COMMENT:      handleComment,
	MACHINE_INFO: handleMachineInfo,
	COMPRESSION:  handleCompression,
}

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
	compressionCode, err := byteorder.Binary2Uint16(0, bytes[0:2]...)
	if err != nil {
		return err
	}
	mfer.Control.CompressionCode = compressionCode

	dataLength, err := byteorder.Binary2Uint32(0, bytes[2:6]...)
	if err != nil {
		return err
	}
	_ = fmt.Sprintf("コード = %d, データ長 = %d\n", compressionCode, dataLength)
	return nil
}
