package main

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
)

// フレーム関連のハンドラー
func handleBlock(mfer *Mfer, bytes []byte, length uint32) error {
	blockLength, err := Binary2Uint32(mfer.Control.ByteOrder, bytes...)
	if err != nil {
		return err
	}
	lastIndex := len(mfer.Frames) - 1
	mfer.Frames[lastIndex].BlockLength = blockLength
	color.Red("BLOCK")
	return nil
}

func handleChannel(mfer *Mfer, bytes []byte, length uint32) error {
	num, err := Binary2Uint32(mfer.Control.ByteOrder, bytes...)
	if err != nil {
		return err
	}
	lastIndex := len(mfer.Frames) - 1
	mfer.Frames[lastIndex].NumChannel = num
	fmt.Printf("  チャンネル数 = %d\n", mfer.Frames[0].NumChannel)
	return nil
}

func handleSequence(mfer *Mfer, bytes []byte, length uint32) error {
	num, err := Binary2Uint32(mfer.Control.ByteOrder, bytes...)
	if err != nil {
		return err
	}
	lastIndex := len(mfer.Frames) - 1
	mfer.Frames[lastIndex].NumSequence = num
	color.Red("シーケンス")
	return nil
}

func handleFPointer(fer *Mfer, bytes []byte, length uint32) error {
	color.Green("F_POINTER")
	return nil
}

func handleWaveFormType(mfer *Mfer, bytes []byte, length uint32) error {
	code, err := Binary2Uint16(mfer.Control.ByteOrder, bytes[0:2]...)
	if err != nil {
		return err
	}
	lastIndex := len(mfer.Frames) - 1
	mfer.Frames[lastIndex].WaveForm.Code = code
	mfer.Frames[lastIndex].WaveForm.Description = string(bytes[2:length])
	fmt.Printf("Wave From Code = %d, 説明 = %s\n", code, string(bytes[2:length]))
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
	fmt.Printf("%+v\n", mfer.Frames[0].Channels[channelNumber])
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
	fmt.Printf("length = %d\n", length)
	color.Green("DATA: No Imprementation")
	//lastIndex := len(mfer.Frames) - 1
	return nil
}
