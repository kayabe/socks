package s5

import (
	"io"
	"unsafe"
)

const AuthUserPWVersion = 0x1

// AuthUserPW is the username/password authentication method.
type AuthUserPW struct {
	Version        uint8 // unrequired, only used by Unpack
	UsernameLength uint8 // unrequired, only used by Unpack
	Username       []byte
	PasswordLength uint8 // unrequired, only used by Unpack
	Password       []byte
}

// Pack writes the authentication method to the given writer.
func (auth *AuthUserPW) Pack(w io.Writer) (err error) {
	usernameLength := uint8(len(auth.Username))
	passwordLength := uint8(len(auth.Password))

	buf := make([]byte, 1+usernameLength+1+passwordLength+1)
	buf[0] = AuthUserPWVersion
	buf[1] = usernameLength                     // username length
	copy(buf[2:], auth.Username)                // copy username to the buffer
	buf[2+usernameLength] = passwordLength      // password length
	copy(buf[3+usernameLength:], auth.Password) // copy password to the buffer

	// Write the entire buffer to the provided io.Writer
	_, err = w.Write(buf)
	return
}

// Unpack reads the authentication method from the given reader.
func (auth *AuthUserPW) Unpack(r io.Reader) (err error) {
	// casting the auth pointer to an array of 2 bytes so we can directly read only into Version and UsernameLength.
	if _, err = io.ReadFull(r, (*[2]byte)(unsafe.Pointer(auth))[:]); err != nil {
		return
	}
	// Check if the read version matches with the protocol's version
	if auth.Version != AuthUserPWVersion {
		return ErrAuthVersion
	}
	auth.Username = make([]byte, auth.UsernameLength)
	if _, err = io.ReadFull(r, auth.Username); err != nil {
		return
	}
	// casting the auth.PasswordLength pointer to an array of one byte so we can read into it.
	if _, err = io.ReadFull(r, (*[1]byte)(unsafe.Pointer(&auth.PasswordLength))[:]); err != nil {
		return
	}
	auth.Password = make([]byte, auth.PasswordLength)
	_, err = io.ReadFull(r, auth.Password)
	return
}

func (AuthUserPW) Method() uint8 {
	return uint8(MethodAuthUserPW)
}
