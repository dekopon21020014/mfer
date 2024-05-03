package main

const (
	// for Mfer.Sampling
	INTERVAL    = 0x0b
	SENSITIVITY = 0x0c
	DATA_TYPE   = 0x0a
	OFFSET      = 0x0d
	NULL        = 0x12

	// for Mfer.Frame
	BLOCK    = 0x04
	CHANNEL  = 0x05
	SEQUENCE = 0x06
	POINTER  = 0x07

	// for Mfer.WaveForm
	WAVE_FORM_TYPE    = 0x08
	CHANNEL_ATTRIBUTE = 0x3f
	LDN               = 0x09
	INFORMATION       = 0x15
	FILTER            = 0x11
	IDP = 0x0f
	DATA = 0x1e

	PREAMBLE = 0x40
)
