package main

import "errors"

// サンプリング関連のハンドラー
func handleInterval(mfer *Mfer, bytes []byte, length uint32) error {
	/*
		if len(bytes) < 6 {
			return errors.New("インターバル用のバイト数が不足")
		}
	*/
	mfer.Sampling.Interval.UnitCode = bytes[0]
	mfer.Sampling.Interval.Exponent = int8(bytes[1])
	mantissa, err := Binary2Uint32(mfer.Control.ByteOrder, bytes[2:length]...)
	if err != nil {
		return err
	}
	mfer.Sampling.Interval.Mantissa = mantissa
	return nil
}

func handleSensitivity(mfer *Mfer, bytes []byte, length uint32) error {
	/*
		if len(bytes) < 6 {
			return errors.New("感度用のバイト数が不足")
		}
	*/
	mfer.Sampling.Sensitivity.UnitCode = bytes[0]
	mfer.Sampling.Sensitivity.Exponent = int8(bytes[1])
	mantissa, err := Binary2Uint32(mfer.Control.ByteOrder, bytes[2:length]...)
	if err != nil {
		return err
	}
	mfer.Sampling.Sensitivity.Mantissa = mantissa
	return nil
}

func handleDataType(mfer *Mfer, bytes []byte, length uint32) error {
	if len(bytes) < 1 {
		return errors.New("データタイプ用のバイト数が不足")
	}
	mfer.Sampling.DataTypeCode = bytes[0]
	return nil
}

func handleOffset(mfer *Mfer, bytes []byte, length uint32) error {
	offset, err := Binary2Uint64(mfer.Control.ByteOrder, bytes...)
	if err != nil {
		return err
	}
	mfer.Sampling.Offset = offset
	return nil
}

func handleNull(mfer *Mfer, bytes []byte, length uint32) error {
	mfer.Sampling.Null = 0
	return nil
}
