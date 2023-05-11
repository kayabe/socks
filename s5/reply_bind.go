package s5

import (
	"encoding/binary"
	"io"
	"net/netip"
	"unsafe"
)

// Bind structure for IPv4 or IPv6
type ReplyBind struct {
	Address netip.Addr
	Port    uint16
}

func (b *ReplyBind) Pack(buf []byte) {
	size := b.Size()

	// reinterpret address structure as [2]uint64, // 0 is high, 1 is low
	v := (*[2]uint64)(unsafe.Pointer(&b.Address))

	if size == 18 {
		// convert IPv6 long to [16]byte and add it to buf
		binary.BigEndian.PutUint64(buf, v[0])
		binary.BigEndian.PutUint64(buf[8:], v[1])
	} else {
		// convert IPv4 long to [4]byte and add it to buf
		binary.BigEndian.PutUint32(buf, uint32(v[1]))
	}

	// convert uint16 port to [2]byte
	buf[size-2] = byte(b.Port >> 8)
	buf[size-1] = byte(b.Port)
}

func (b *ReplyBind) Unpack(r io.Reader, v6 bool) (err error) {
	var size uint8
	if v6 {
		size = 16 + 2
	} else {
		size = 4 + 2
	}

	// allocate size based on address type
	buf := make([]byte, size)
	if _, err = io.ReadFull(r, buf); err != nil {
		return
	}
	b.Address, _ = netip.AddrFromSlice(buf[:size-2])
	b.Port = uint16(buf[size-2])<<8 | uint16(buf[size-1])
	return
}

// Kind ...
func (b *ReplyBind) Kind() uint8 {
	if b.Address.Is6() {
		return AddressTypeIPv6
	}
	return AddressTypeIPv4
}

func (b *ReplyBind) Size() uint8 {
	if b.Address.Is6() {
		return 16 + 2
	}
	return 4 + 2
}
