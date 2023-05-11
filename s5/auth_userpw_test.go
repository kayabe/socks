package s5

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var authBuf = []byte{0x01, 0x04, 't', 'e', 's', 't', 0x04, 't', 'e', 's', 't'}

func TestAuthUserPWPack(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	auth := AuthUserPW{Username: []byte("test"), Password: []byte("test")}

	if err := auth.Pack(buf); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, authBuf, buf.Bytes())
}

func TestAuthUserPWUnpack(t *testing.T) {
	buf := bytes.NewBuffer(authBuf)

	var auth AuthUserPW
	if err := auth.Unpack(buf); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, AuthUserPW{Version: AuthUserPWVersion, UsernameLength: 4, Username: []uint8("test"), PasswordLength: 4, Password: []uint8("test")}, auth)
}

func BenchmarkAuthUserPWPack(b *testing.B) {
	buf := bytes.NewBuffer(nil)
	auth := AuthUserPW{Username: []byte("test"), Password: []byte("test")}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = auth.Pack(buf)
		buf.Reset()
	}
}

func BenchmarkAuthUserPWUnpack(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(authBuf)
		_ = (&AuthUserPW{}).Unpack(buf)
	}
}
