package std12lead

import (
	"bytes"
	"encoding/binary"
	"fmt"

	m "github.com/dekopon21020014/mfer/pkg/mfer"
)

const (
	bytesPerBlock   = 2
	signedInt16Type = 0 // データタイプのコード
)

// LeadCalculator は12誘導心電図の計算に必要な情報を保持する構造体です
type LeadCalculator struct {
	blockLength     int
	data            []byte
	blockTotalBytes int
}

// NewLeadCalculator は LeadCalculator の新しいインスタンスを生成します
func NewLeadCalculator(mfer *m.Mfer) (*LeadCalculator, error) {
	if mfer.Sampling.DataTypeCode != signedInt16Type {
		return nil, fmt.Errorf("unsupported data type: expected signed 16-bit integer")
	}

	blockLength := mfer.Frames[0].BlockLength
	return &LeadCalculator{
		blockLength:     int(blockLength),
		data:            mfer.Frames[0].WaveForm.Data,
		blockTotalBytes: int(blockLength) * bytesPerBlock,
	}, nil
}

// Convert8To12Lead は8誘導から12誘導心電図への変換を行います
func (lc *LeadCalculator) Convert8To12Lead() (map[string][]byte, error) {
	lead := make(map[string][]byte)

	// 基本的な8誘導の設定
	basicLeadNames := []string{"I", "II", "V1", "V2", "V3", "V4", "V5", "V6"}
	for i, name := range basicLeadNames {
		begin := i * lc.blockTotalBytes
		end := begin + lc.blockTotalBytes
		if begin >= 0 && end <= len(lc.data) && begin < end {
			lead[name] = lc.data[begin:end]
		} else {
			return nil, fmt.Errorf("index out of bounds for lead %s", name)
		}
	}

	// 追加の誘導の計算
	additionalLeads := []string{"III", "aVR", "aVL", "aVF"}
	for _, name := range additionalLeads {
		calculatedLead, err := lc.calculateLead(name, lead)
		if err != nil {
			return nil, err
		}
		lead[name] = calculatedLead
	}

	return lead, nil
}

// calculateLead は特定の誘導を計算します
func (lc *LeadCalculator) calculateLead(targetLeadName string, lead map[string][]byte) ([]byte, error) {
	operand1, operand2, err := lc.getOperands(targetLeadName, lead)
	if err != nil {
		return nil, err
	}

	targetLead := make([]byte, 0, lc.blockTotalBytes)

	for i := 0; i < lc.blockTotalBytes; i += bytesPerBlock {
		sum := lc.calculateLeadValue(operand1[i:i+bytesPerBlock], operand2[i:i+bytesPerBlock], targetLeadName)

		buf := new(bytes.Buffer)
		if err := binary.Write(buf, binary.LittleEndian, sum); err != nil {
			return nil, err
		}
		targetLead = append(targetLead, buf.Bytes()...)
	}

	return targetLead, nil
}

// getOperands は対象の誘導に応じたオペランドを取得します
func (lc *LeadCalculator) getOperands(targetLeadName string, lead map[string][]byte) ([]byte, []byte, error) {
	switch targetLeadName {
	case "III": // III = II - I
		return lead["II"], lead["I"], nil
	case "aVR": // aVR = -(I + II) / 2
		return lead["I"], lead["II"], nil
	case "aVL": // aVL = (I - III) / 2
		return lead["I"], lead["III"], nil
	case "aVF": // aVF = (II + III) / 2
		return lead["II"], lead["III"], nil
	default:
		return nil, nil, fmt.Errorf("unknown lead: %s", targetLeadName)
	}
}

// calculateLeadValue は誘導の値を計算します
func (lc *LeadCalculator) calculateLeadValue(operand1, operand2 []byte, targetLeadName string) int16 {
	var val_operand1, val_operand2 int16
	binary.Read(bytes.NewReader(operand1), binary.LittleEndian, &val_operand1)
	binary.Read(bytes.NewReader(operand2), binary.LittleEndian, &val_operand2)

	switch targetLeadName {
	case "III": // III = II - I
		return val_operand1 - val_operand2
	case "aVR": // aVR = -(I + II) / 2
		return -(val_operand1 + val_operand2) / 2
	case "aVL": // aVL = (I - III) / 2
		return (val_operand1 - val_operand2) / 2
	case "aVF": // aVF = (II + III) / 2
		return (val_operand1 + val_operand2) / 2
	default:
		return 0
	}
}

// Convert8To12Lead は従来の関数インターフェースを維持するためのヘルパー関数
func Convert8To12Lead(mfer *m.Mfer) map[string][]byte {
	calculator, err := NewLeadCalculator(mfer)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	leads, err := calculator.Convert8To12Lead()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return leads
}
