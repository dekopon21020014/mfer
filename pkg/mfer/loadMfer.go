package mfer

import (
	"encoding/binary"
	"fmt"
	"os"
)

// tagHandler は異なるタグコードを処理する関数型
type tagHandler func(mfer *Mfer, bytes []byte, length uint32) error

// loadMferはリファクタリングされ、可読性とメンテナンス性が向上
func LoadMfer(mfer Mfer, path string) (Mfer, error) {
	// ファイル全体を読み込む
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return mfer, fmt.Errorf("ファイルの読み込みに失敗: %w", err)
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
			tmp := fileBytes[i]
			fileBytes[i] = fileBytes[i+1] + 1
			fileBytes[i+1] = tmp
		}

		// 拡張長さ対応でlengthをパース
		length, bytesRead, err := parseLength(fileBytes[i:])
		if err != nil {
			return mfer, err
		}
		i += bytesRead

		handler := func() tagHandler {
			structType := groupMap[tagCode]
			switch structType {
			case "Sampling":
				return samplingHandlerMap[tagCode]
			case "Frame":
				if len(mfer.Frames) > 0 && mfer.Frames[len(mfer.Frames)-1].WaveForm.Data != nil {
					mfer.Frames = append(mfer.Frames, newFrame())
				}
				return frameHandlerMap[tagCode]
			case "Control":
				return controlHandlerMap[tagCode]
			case "Extensions":
				return extensionHandlerMap[tagCode]
			case "Helper":
				return helperHandlerMap[tagCode]
			default:
				return nil
			}
		}()

		err = handler(&mfer, fileBytes[i:i+int(length)], length)
		if err != nil {
			return mfer, fmt.Errorf("タグ %02x の処理中にエラー: %w", tagCode, err)
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
