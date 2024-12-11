package mfer

var groupMap = map[byte]string{
	// Mfer.Sampling
	INTERVAL:    "Sampling",
	SENSITIVITY: "Sampling",
	DATA_TYPE:   "Sampling",
	OFFSET:      "Sampling",
	NULL:        "Sampling",

	// Mfer.Frame
	BLOCK:             "Frame",
	CHANNEL:           "Frame",
	SEQUENCE:          "Frame",
	F_POINTER:         "Frame",
	WAVE_FORM_TYPE:    "Frame",
	CHANNEL_ATTRIBUTE: "Frame",
	LDN:               "Frame",
	INFORMATION:       "Frame",
	FILTER:            "Frame",
	IPD:               "Frame",
	DATA:              "Frame",

	// Mfer.Control
	BYTE_ORDER:   "Control",
	VERSION:      "Control",
	CHAR_CODE:    "Control",
	ZERO:         "Control",
	COMMENT:      "Control",
	MACHINE_INFO: "Control",
	COMPRESSION:  "Control",

	// Extensions
	PREAMBLE:  "Extensions",
	EVENT:     "Extensions",
	VALUE:     "Extensions",
	CONDITION: "Extensions",
	ERROR:     "Extensions",
	GROUP:     "Extensions",
	R_POINTER: "Extensions",
	SIGNITURE: "Extensions",

	// Mfer.Helper
	P_NAME:  "Helper",
	P_ID:    "Helper",
	P_AGE:   "Helper",
	P_SEX:   "Helper",
	TIME:    "Helper",
	MESSAGE: "Helper",
	UID:     "Helper",
	MAP:     "Helper",
	END:     "Helper",
}
