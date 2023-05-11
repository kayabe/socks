package s5

import (
	"io"
	"unsafe"
)

// HandshakeRequest is the handshake request for SOCKS V5.
type HandshakeRequest struct {
	Version  uint8 // unrequired, only used by Unpack
	NMethods uint8 // unrequired, only used by Unpack
	Methods  []AuthMethod
}

func (h *HandshakeRequest) Pack(w io.Writer) (err error) {
	// Number of methods
	nMethods := len(h.Methods)
	if nMethods > 0xFF {
		return ErrMethodsLimit
	}

	buf := make([]byte, 2+nMethods)
	buf[0] = VERSION
	buf[1] = uint8(nMethods)

	if nMethods > 0 {
		copy(buf[2:], *(*[]byte)(unsafe.Pointer(&h.Methods)))
	}

	// Write the entire buffer to the provided io.Writer
	_, err = w.Write(buf)
	return
}

func (h *HandshakeRequest) Unpack(r io.Reader) (err error) {
	// Always read the first two bytes from the reader into the struct's Version, NMethods
	if _, err = io.ReadFull(r, (*[2]byte)(unsafe.Pointer(h))[:]); err != nil {
		return
	}

	// Check if the version matches
	if h.Version != VERSION {
		return ErrUnsupportedVersion
	}

	if h.NMethods > 0 {
		// Allocation needed for unsafe.Pointer
		h.Methods = make([]AuthMethod, h.NMethods)

		// Read NMethods of bytes from the reader into the struct's Methods
		if _, err = io.ReadFull(r, *(*[]byte)(unsafe.Pointer(&h.Methods))); err != nil {
			return
		}
	}

	return
}

// HandshakeReply is the handshake reply for SOCKS V5.
type HandshakeReply struct {
	Version uint8
	Method  AuthMethod
}

func (t *HandshakeReply) Pack(w io.Writer) (err error) {
	// Write the entire buffer to the provided io.Writer
	_, err = w.Write([]byte{t.Version, byte(t.Method)})
	return
}

func (t *HandshakeReply) Unpack(r io.Reader) (err error) {
	// Read the first two bytes from the reader
	if _, err = io.ReadFull(r, (*[2]byte)(unsafe.Pointer(t))[:]); err != nil {
		return
	}

	// Check if the version matches
	if t.Version != VERSION {
		return ErrUnsupportedVersion
	}

	return
}
