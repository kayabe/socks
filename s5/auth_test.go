package s5

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthPack(t *testing.T) {
	var buf = bytes.NewBuffer(nil)
	_ = (&AuthReply{Status: ReplySuccess}).Pack(buf)

	assert.Equal(t, []byte{VERSION, byte(ReplySuccess)}, buf.Bytes())
}

func TestAuthUnpack(t *testing.T) {
	var buf = bytes.NewBuffer([]byte{VERSION, byte(ReplySuccess)})
	var r AuthReply
	_ = r.Unpack(buf)

	assert.Equal(t, AuthReply{VERSION, ReplySuccess}, r)
}

func BenchmarkAuthPack(b *testing.B) {
	var buf = bytes.NewBuffer(nil)
	var r = AuthReply{Status: ReplySuccess}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = r.Pack(buf)
		buf.Reset()
	}
}

func BenchmarkAuthUnpack(b *testing.B) {
	var data = []byte{VERSION, byte(ReplySuccess)}
	var r AuthReply

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var buf = bytes.NewBuffer(data)
		_ = r.Unpack(buf)
		buf.Reset()
	}
}
