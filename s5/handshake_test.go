package s5

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandshakeRequestPack(t *testing.T) {
	var buf = bytes.NewBuffer(nil)

	if err := (&HandshakeRequest{Methods: []AuthMethod{MethodAuthNone}}).Pack(buf); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, []byte{0x05, 0x01, 0x00}, buf.Bytes(), "the handshake request should be equal")
}

func TestHandshakeRequestUnpack(t *testing.T) {
	var r HandshakeRequest

	if err := r.Unpack(bytes.NewBuffer([]byte{0x05, 0x01, 0x00})); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, HandshakeRequest{Version: 0x05, NMethods: 0x01, Methods: []AuthMethod{MethodAuthNone}}, r, "the handshake reply should be equal")
}

func TestHandshakeRequestUnpack2(t *testing.T) {
	var r HandshakeRequest

	if err := r.Unpack(bytes.NewBuffer([]byte{0x05, 0x00})); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, HandshakeRequest{Version: 0x05, NMethods: 0x00}, r, "the handshake reply should be equal")
}

func TestHandshakeRequestUnpackFail(t *testing.T) {
	var r HandshakeRequest

	if err := r.Unpack(bytes.NewBuffer([]byte{'H', 'T', 'T', 'P'})); err != nil {
		assert.Equal(t, ErrUnsupportedVersion, err)
	}
}

func TestHandshakeRequestUnpackFail2(t *testing.T) {
	var r HandshakeRequest

	if err := r.Unpack(bytes.NewBuffer([]byte{})); err != nil {
		assert.Equal(t, io.EOF, err)
	}
}

func TestHandshakeRequestUnpackFail3(t *testing.T) {
	var r HandshakeRequest

	if err := r.Unpack(bytes.NewBuffer([]byte{0x05, 0x01})); err != nil {
		assert.Equal(t, io.EOF, err)
	}
}

func TestHandshakeReplyUnpack(t *testing.T) {
	var r HandshakeReply

	if err := r.Unpack(bytes.NewBuffer([]byte{0x05, 0x00})); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, HandshakeReply{Version: 0x05, Method: MethodAuthNone}, r, "the handshake reply should be equal")
}

func TestHandshakeReplyPack(t *testing.T) {
	var r = HandshakeReply{Version: 0x05, Method: MethodAuthNone}
	var buf = bytes.NewBuffer(nil)

	if err := r.Pack(buf); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, []byte{0x05, 0x00}, buf.Bytes(), "the handshake reply should be equal")
}

func BenchmarkHandshakeRequestPack(b *testing.B) {
	var buf = bytes.NewBuffer(nil)

	var methods []AuthMethod
	for i := 0; i < 0xff; i++ {
		methods = append(methods, MethodAuthNone)
	}
	var r = HandshakeRequest{Methods: methods}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = r.Pack(buf)
		buf.Reset()
	}
}

func BenchmarkHandshakeRequestUnpack(b *testing.B) {
	data := []byte{0x05, 0x01, 0x00}
	var r HandshakeRequest

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(data)
		_ = r.Unpack(buf)
		buf.Reset()
	}
}

func BenchmarkHandshakeReplyUnpack(b *testing.B) {
	var data = []byte{0x05, 0x00}
	var r HandshakeReply

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(data)
		_ = r.Unpack(buf)
		buf.Reset()
	}
}

func BenchmarkHandshakeReplyPack(b *testing.B) {
	var buf = bytes.NewBuffer(nil)
	var r = HandshakeReply{Version: VERSION, Method: AuthUserPWVersion}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = r.Pack(buf)
		buf.Reset()
	}
}
