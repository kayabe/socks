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

// AuthReply is the reply to a SOCKS5 authentication method.
type AuthReply struct {
	Version uint8 // unrequired, only used by Unpack
	Status  ReplyStatus
}

// Pack for server-side
func (t *AuthReply) Pack(w io.Writer) (err error) {
	_, err = w.Write([]byte{VERSION, byte(t.Status)})
	return
}

// Unpack for client-side
func (t *AuthReply) Unpack(r io.Reader) (err error) {
	// casting the t pointer to an array of 2 bytes so we can directly read only into Version and Status.
	if _, err = io.ReadFull(r, (*[2]byte)(unsafe.Pointer(t))[:]); err != nil {
		return
	}

	// Check if the version matches
	if t.Version != VERSION {
		return ErrUnsupportedVersion
	}

	return
}

// AuthReplyFromConn for client-side
func AuthReplyFromConn(conn net.Conn) (_ *AuthReply, err error) {
	var reply = new(AuthReply)
	if err = reply.Unpack(conn); err != nil {
		return nil, err
	}
	return reply, nil
}

// AuthHandler used to determine if the authentication is valid or not
type AuthHandler func(authMethod AuthMethod, authData any) bool
