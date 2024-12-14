package mfer

import (
	"errors"
	"fmt"

	"github.com/dekopon21020014/mfer/pkg/byteorder"
	"github.com/fatih/color"
)

var frameHandlerMap = map[byte]tagHandler{
	// Mfer.Frame
	BLOCK:             handleBlock,
	CHANNEL:           handleChannel,
	SEQUENCE:          handleSequence,
	F_POINTER:         handleFPointer,
	WAVE_FORM_TYPE:    handleWaveFormType,
	CHANNEL_ATTRIBUTE: handleChannelAttribute,
	LDN:               handleLdn,
	INFORMATION:       handleInformation,
	FILTER:            handleFilter,
	IPD:               handleIpd,
	DATA:              handleData,
}

// フレーム関連のハンドラー
func handleBlock(mfer *Mfer, bytes []byte, length uint32) error {
	blockLength, err := byteorder.Binary2Uint32(mfer.Control.ByteOrder, bytes...)
	if err != nil {
		return err
	}
	lastIndex := len(mfer.Frames) - 1
	mfer.Frames[lastIndex].BlockLength = blockLength
	return nil
}

func handleChannel(mfer *Mfer, bytes []byte, length uint32) error {
	num, err := byteorder.Binary2Uint32(mfer.Control.ByteOrder, bytes...)
	if err != nil {
		return err
	}
	lastIndex := len(mfer.Frames) - 1
	mfer.Frames[lastIndex].NumChannel = num
	return nil
}

func handleSequence(mfer *Mfer, bytes []byte, length uint32) error {
	num, err := byteorder.Binary2Uint32(mfer.Control.ByteOrder, bytes...)
	if err != nil {
		return err
	}
	lastIndex := len(mfer.Frames) - 1
	mfer.Frames[lastIndex].NumSequence = num
	return nil
}

func handleFPointer(fer *Mfer, bytes []byte, length uint32) error {
	color.Green("F_POINTER")
	return nil
}

func handleWaveFormType(mfer *Mfer, bytes []byte, length uint32) error {
	code, err := byteorder.Binary2Uint16(mfer.Control.ByteOrder, bytes[0:2]...)
	if err != nil {
		return err
	}
	lastIndex := len(mfer.Frames) - 1
	mfer.Frames[lastIndex].WaveForm.Code = code
	mfer.Frames[lastIndex].WaveForm.Description = string(bytes[2:length])
	return nil
}

func handleChannelAttribute(mfer *Mfer, bytes []byte, length uint32) error {
	channelNumber := int(bytes[0])
	lastIndex := len(mfer.Frames) - 1
	mfer.Frames[lastIndex].Channels[channelNumber] = Channel{
		Num:     channelNumber,
		TagCode: bytes[1],
		Length:  bytes[2],
		Data:    bytes[3:],
	}
	return nil
}

func handleLdn(mfer *Mfer, bytes []byte, length uint32) error {
	if len(bytes) < 1 {
		return errors.New("LDNのバイト数が不足")
	}
	lastIndex := len(mfer.Frames) - 1
	mfer.Frames[lastIndex].WaveForm.Ldn = bytes[0]
	return nil
}

func handleInformation(mfer *Mfer, bytes []byte, length uint32) error {
	color.Green("情報: 未実装")
	lastIndex := len(mfer.Frames) - 1
	mfer.Frames[lastIndex].WaveForm.Information = "情報が存在します"
	return nil
}

func handleFilter(mfer *Mfer, bytes []byte, length uint32) error {
	lastIndex := len(mfer.Frames) - 1
	mfer.Frames[lastIndex].WaveForm.Filter = string(bytes)
	return nil
}

func handleIpd(mfer *Mfer, bytes []byte, length uint32) error {
	if len(bytes) < 1 {
		return errors.New("IPDのバイト数が不足")
	}
	lastIndex := len(mfer.Frames) - 1
	mfer.Frames[lastIndex].WaveForm.IpdCode = bytes[0]
	return nil
}

func handleData(mfer *Mfer, bytes []byte, length uint32) error {
	fmt.Println("length = ", length)
	lastIndex := len(mfer.Frames) - 1
	mfer.Frames[lastIndex].WaveForm.Data = bytes
	return nil
}
