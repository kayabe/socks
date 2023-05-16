package s5

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthPack(t *testing.T) {
	var buf = bytes.NewBuffer(nil)
	_ = (&AuthReply{Version: AuthUserPWVersion, Status: ReplyGeneralFailure}).Pack(buf)
	assert.Equal(t, []byte{AuthUserPWVersion, byte(ReplyGeneralFailure)}, buf.Bytes())
}

func TestAuthUnpack(t *testing.T) {
	var buf = bytes.NewBuffer([]byte{AuthUserPWVersion, byte(ReplyGeneralFailure)})
	var r AuthReply
	r.Validate = func(Version uint8) error {
		if Version != AuthUserPWVersion {
			return ErrAuthVersion
		}
		return nil
	}
	_ = r.Unpack(buf)
	r.Validate = nil
	assert.Equal(t, AuthReply{Version: AuthUserPWVersion, Status: ReplyGeneralFailure}, r)
}

func BenchmarkAuthPack(b *testing.B) {
	var buf = bytes.NewBuffer(nil)
	var r = AuthReply{Status: ReplyGeneralFailure}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = r.Pack(buf)
		buf.Reset()
	}
}

func BenchmarkAuthUnpack(b *testing.B) {
	var data = []byte{AuthUserPWVersion, byte(ReplyGeneralFailure)}
	var r AuthReply

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var buf = bytes.NewBuffer(data)
		_ = r.Unpack(buf)
		buf.Reset()
	}
}
