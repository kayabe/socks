package s5

import (
	"io"
	"unsafe"
)

// Request is the request for SOCKS V5.
type Request struct {
	Version     uint8       // unrequired, only used by Unpack
	Command     uint8       // required only by Pack
	Reserved    uint8       // unrequired, only used by Unpack
	AddressType uint8       // unrequired, only used by Unpack
	Destination Destination // required only by Pack
}

func (v *Request) Pack(w io.Writer) (err error) {
	buf := make([]byte, 4+v.Destination.Size())
	buf[0] = VERSION
	buf[1] = v.Command
	buf[3] = v.Destination.Kind()
	if err = v.Destination.Put(buf[4:]); err != nil {
		return
	}
	_, err = w.Write(buf)
	return
}

func (v *Request) Unpack(r io.Reader) (err error) {
	if _, err = io.ReadFull(r, (*[4]byte)(unsafe.Pointer(v))[:]); err != nil {
		return
	}
	if v.Version != VERSION {
		return ErrUnsupportedVersion
	}
	switch v.AddressType {
	case AddressTypeDomainName:
		v.Destination = new(RequestV5DestDomainName)
	case AddressTypeIPv4:
		v.Destination = new(RequestV5DestIPv4)
	case AddressTypeIPv6:
		v.Destination = new(RequestV5DestIPv6)
	default:
		return ErrUnsupportedAddressType
	}
	return v.Destination.Unpack(r)
}
