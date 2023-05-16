package s5

import (
	"io"
	"net"
	"unsafe"
)

// AuthMethod ...
//
//	X'03' to X'7F' IANA ASSIGNED
//	X'80' to X'FE' RESERVED FOR PRIVATE METHODS
type AuthMethod uint8

const (
	MethodAuthNone           AuthMethod = iota // no authentication required
	MethodAuthGSSAPI                           // GSSAPI
	MethodAuthUserPW                           // username/password
	MethodAuthCHAP                             // Challenge-Handshake Authentication Protocol
	_                                          // Unassigned
	MerhodAuthCRAM                             // Challenge-Response Authentication Method
	MethodAuthSSL                              // Secure Sockets Layer
	MethodAuthNDS                              // NDS Authentication
	MethodAuthMAF                              // Multi-Authentication Framework
	MethodAuthJSB                              // JSON Parameter Block
	MethodAuthNoneAcceptable AuthMethod = 0xFF // no acceptable methods
)

// AuthReply is the reply after SOCKS5 authentication method.
type AuthReply struct {
	Version  uint8 // required, should be the version of the authentication method
	Status   ReplyStatus
	Validate func(uint8) error
}

// Pack for server-side
func (t *AuthReply) Pack(w io.Writer) (err error) {
	_, err = w.Write([]byte{t.Version, byte(t.Status)})
	return
}

// Unpack for client-side
func (t *AuthReply) Unpack(r io.Reader) (err error) {
	// casting the t pointer to an array of 2 bytes so we can directly read only into Version and Status.
	if _, err = io.ReadFull(r, (*[2]byte)(unsafe.Pointer(t))[:]); err != nil {
		return
	}
	if t.Validate != nil {
		err = t.Validate(t.Version)
	}
	return
}

// AuthReplyFromConn for client-side
func AuthReplyFromConn(conn net.Conn, validateCb func(authVersion uint8) error) (_ *AuthReply, err error) {
	var reply = &AuthReply{
		Validate: validateCb,
	}
	if err = reply.Unpack(conn); err != nil {
		return nil, err
	}
	return reply, nil
}

// AuthHandler used to determine if the authentication is valid or not
type AuthHandler func(authMethod AuthMethod, authData any) bool
