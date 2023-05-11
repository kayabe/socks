package s5

import (
	"io"
	"unsafe"
)

const (
	DestIPv4Size = 4 + 2  // IPv4 + Port
	DestIPv6Size = 16 + 2 // IPv6 + Port
)

type Destination interface {
	// Put has to be used appropriately together with Size
	Put([]byte) error

	// Size returns binary length of the struct
	Size() uint8

	// Unpack ...
	Unpack(r io.Reader) error

	// Kind returns AddressType
	Kind() uint8
}

// RequestV5DestDomainName is a domaim name destination.
type RequestV5DestDomainName struct {
	Address []byte
	Port    uint16
}

// Kind returns AddressType
func (h *RequestV5DestDomainName) Kind() uint8 { return AddressTypeDomainName }

// Size returns binary length of struct
func (h *RequestV5DestDomainName) Size() uint8 { return 1 + uint8(len(h.Address)) + 2 }

func (h *RequestV5DestDomainName) Put(buf []byte) error {
	size := h.Size()

	if len(buf) < int(size) {
		return ErrInvalidBufferSize
	}

	buf[0] = size - 3        // address length
	copy(buf[1:], h.Address) // address

	// port
	buf[size-2] = byte(h.Port >> 8)
	buf[size-1] = byte(h.Port)

	return nil
}

// ReadFrom reads from the given reader into the structure.
func (h *RequestV5DestDomainName) Unpack(r io.Reader) (err error) {
	var length uint8
	if _, err = io.ReadFull(r, (*[1]byte)(unsafe.Pointer(&length))[:]); err != nil {
		return
	}
	buf := make([]byte, length+2)
	// read domain and port
	if _, err = io.ReadFull(r, buf); err != nil {
		return
	}
	h.Address = buf[:length]
	h.Port = uint16(buf[length+1]) | uint16(buf[length])<<8
	return
}

///

// RequestV5DestinationIPv4 is a IPv4 destination.
type RequestV5DestIPv4 struct {
	Address [4]byte
	Port    uint16
}

// Kind returns AddressType
func (h *RequestV5DestIPv4) Kind() uint8 { return AddressTypeIPv4 }

// Size returns binary length of struct
func (h *RequestV5DestIPv4) Size() uint8 { return DestIPv4Size }

// Put ..
func (d *RequestV5DestIPv4) Put(buf []byte) error {
	if len(buf) < DestIPv4Size {
		return ErrInvalidBufferSize
	}

	// adress
	copy(buf, d.Address[:])
	// port
	buf[4] = byte(d.Port >> 8)
	buf[5] = byte(d.Port)
	return nil
}

func (d *RequestV5DestIPv4) Unpack(r io.Reader) (err error) {
	_, err = io.ReadFull(r, (*[6]byte)(unsafe.Pointer(d))[:])
	return
}

// RequestV5DestinationIPv6 is a IPv6 destination.
type RequestV5DestIPv6 struct {
	Address [16]byte
	Port    uint16
}

// Kind returns AddressType
func (h *RequestV5DestIPv6) Kind() uint8 { return AddressTypeIPv6 }

// Size returns binary length of struct
func (h *RequestV5DestIPv6) Size() uint8 { return DestIPv6Size }

func (h *RequestV5DestIPv6) Put(buf []byte) error {
	if len(buf) < DestIPv6Size {
		return ErrInvalidBufferSize
	}

	// address
	copy(buf, h.Address[:])

	// port
	buf[16] = byte(h.Port >> 8)
	buf[17] = byte(h.Port)
	return nil
}

func (h *RequestV5DestIPv6) Unpack(r io.Reader) (err error) {
	_, err = io.ReadFull(r, (*[18]byte)(unsafe.Pointer(h))[:])
	return
}
