package mfer2physionet

func Convert(leads map[string][]byte) []byte {
	var physioNetData []byte
	for i := 0; i < len(leads["I"]); i += 2 {
		for _, name := range []string{
			"I", "II", "III",
			"aVR", "aVL", "aVF",
			"V1", "V2", "V3", "V4", "V5", "V6",
		} {
			physioNetData = append(physioNetData, leads[name][i:i+2]...)
		}
	}
	return physioNetData
}
