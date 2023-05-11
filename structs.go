package socks

type (
	ProtocolV4  uint8
	ProtocolV4A uint8
	ProtocolV5  uint8
)

const (
	V4  ProtocolV4  = 0x04
	V4A ProtocolV4A = 0x04
	V5  ProtocolV5  = 0x05
)
