package s5

import (
	"io"
	"unsafe"
)

// Reply is a reply structure for SOCKS V5.
type Reply struct {
	Version     uint8 // unrequired, only used by Unpack
	Status      ReplyStatus
	Reserved    uint8 // unrequired, only used by Unpack
	AddressType uint8
	Bind        ReplyBind
}

// Pack writes the structure to the given writer as bytes.
func (t *Reply) Pack(w io.Writer) (err error) {
	var size uint8 = 4 // Version, Status, Reserved, AddressType
	var buf = make([]byte, size+t.Bind.Size())
	buf[0] = VERSION
	buf[1] = byte(t.Status)
	buf[3] = t.Bind.Kind()
	t.Bind.Pack(buf[4:])

	_, err = w.Write(buf)
	return
}

// Unpack reads from the given reader into the structure.
func (t *Reply) Unpack(r io.Reader) (err error) {
	// casting the t pointer to an array of 4 bytes
	// so we can directly read only into Version, Status, Reserved and AddressType.
	if _, err = io.ReadFull(r, (*[4]byte)(unsafe.Pointer(t))[:]); err != nil {
		return
	}
	if t.Version != VERSION {
		return ErrUnsupportedVersion
	}
	switch t.AddressType {
	case AddressTypeIPv4, AddressTypeIPv6:
	default:
		return ErrUnsupportedAddressType
	}
	return t.Bind.Unpack(r, t.AddressType == AddressTypeIPv6)
}
