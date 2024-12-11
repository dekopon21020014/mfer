package main

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/fatih/color"
)

// tagHandler は異なるタグコードを処理する関数型
type tagHandler func(mfer *Mfer, bytes []byte, length uint32) error

// loadMferはリファクタリングされ、可読性とメンテナンス性が向上
func loadMfer(mfer Mfer, path string) (Mfer, error) {
	// ファイル全体を読み込む
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return mfer, fmt.Errorf("ファイルの読み込みに失敗: %w", err)
	}

	// タグハンドラーのマップを作成
	handlers := map[byte]tagHandler{
		// Mfer.Sampling
		INTERVAL:    handleInterval,
		SENSITIVITY: handleSensitivity,
		DATA_TYPE:   handleDataType,
		OFFSET:      handleOffset,
		NULL:        handleNull,

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

		// Mfer.Control
		BYTE_ORDER:   handleByteOrder,
		VERSION:      handleVersion,
		CHAR_CODE:    handleCharCode,
		ZERO:         handleZero,
		COMMENT:      handleComment,
		MACHINE_INFO: handleMachineInfo,
		COMPRESSION:  handleCompression,

		// Mfer.Extensions
		PREAMBLE:  handlePreamble,
		EVENT:     handleEvent,
		VALUE:     handleValue,
		CONDITION: handleCondition,
		ERROR:     handleError,
		GROUP:     handleGroup,
		R_POINTER: handleRPointer,
		SIGNITURE: handleSignature,

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

	// バイト列を走査
	for i := 0; i < len(fileBytes); {
		tagCode := fileBytes[i]
		i++

		// ZEROをスキップ、ENDで終了
		if tagCode == ZERO {
			continue
		} else if tagCode == END {
			break
		} else if tagCode == CHANNEL_ATTRIBUTE {
			fmt.Println("======CHANNEL ATTRIBUTE========")
			tmp := fileBytes[i]
			fileBytes[i] = fileBytes[i+1] + 1
			fileBytes[i+1] = tmp
			fmt.Println(fileBytes[i], ", ", fileBytes[i+1])
		}

		// 拡張長さ対応でlengthをパース
		length, bytesRead, err := parseLength(fileBytes[i:])
		if err != nil {
			return mfer, err
		}
		i += bytesRead

		// 適切なハンドラを検索して実行
		handler, exists := handlers[tagCode]
		if exists {
			err = handler(&mfer, fileBytes[i:i+int(length)], length)
			if err != nil {
				return mfer, fmt.Errorf("タグ %02x の処理中にエラー: %w", tagCode, err)
			}
		} else {
			// デバッグ用に未処理のタグをログ出力
			color.Yellow("未処理のタグコード: %02x", tagCode)
		}

		// 次のセクションに移動
		i += int(length)
	}

	return mfer, nil
}

// parseLengthは標準および拡張長さのエンコーディングを処理
func parseLength(bytes []byte) (uint32, int, error) {
	length := uint32(bytes[0])
	if length <= 0x7f {
		return length, 1, nil
	}

	numBytes := length - 0x80
	if numBytes > 4 {
		return 0, 0, fmt.Errorf("無効な長さ: %d", numBytes)
	}

	paddedBytes := append(make([]byte, 4-numBytes), bytes[1:1+numBytes]...)
	extendedLength := binary.BigEndian.Uint32(paddedBytes)
	return extendedLength, int(numBytes + 1), nil
}
