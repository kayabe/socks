package s5

import (
	"bytes"
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
)

var addrIPv4 = netip.MustParseAddr("127.0.0.1")
var addrIPv6 = netip.MustParseAddr("2001:0db8:85a3:0000:0000:8a2e:0370:7334")

var testIPv4Buf = []byte{127, 0, 0, 1, 0x07, 0x5c}
var testIPv6Buf = []byte{0x20, 0x01, 0x0d, 0xb8, 0x85, 0xa3, 0x00, 0x00, 0x00, 0x00, 0x8a, 0x2e, 0x03, 0x70, 0x73, 0x34, 0x07, 0x5c}

func TestReplyBindIPv4Pack(t *testing.T) {
	buf := make([]byte, 4+2)
	(&ReplyBind{Address: addrIPv4, Port: 1884}).Pack(buf)
	assert.Equal(t, testIPv4Buf, buf)
}

func TestReplyBindIPv6Pack(t *testing.T) {
	buf := make([]byte, 16+2)
	(&ReplyBind{Address: addrIPv6, Port: 1884}).Pack(buf)
	assert.Equal(t, testIPv6Buf, buf)
}

func TestReplyBindIPv4Unpack(t *testing.T) {
	var buf = bytes.NewBuffer(testIPv4Buf)
	var r ReplyBind

	if err := r.Unpack(buf, false); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, ReplyBind{Address: addrIPv4, Port: 1884}, r)
}

func TestReplyBindIPv6Unpack(t *testing.T) {
	var buf = bytes.NewBuffer(testIPv6Buf)
	var r ReplyBind

	if err := r.Unpack(buf, true); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, ReplyBind{Address: addrIPv6, Port: 1884}, r)
}

func BenchmarkBindIPv4Pack(b *testing.B) {
	buf := make([]byte, 4+2)
	for i := 0; i < b.N; i++ {
		(&ReplyBind{Address: addrIPv4, Port: 1884}).Pack(buf)
	}
}

func BenchmarkBindIPv4Unpack(b *testing.B) {
	var r ReplyBind

	for i := 0; i < b.N; i++ {
		var buf = bytes.NewBuffer(testIPv4Buf)
		_ = r.Unpack(buf, false)
		buf.Reset()
	}
}

func BenchmarkBindIPv6Pack(b *testing.B) {
	buf := make([]byte, 16+2)
	for i := 0; i < b.N; i++ {
		(&ReplyBind{Address: addrIPv6, Port: 1884}).Pack(buf)
	}
}

func BenchmarkBindIPv6Unpack(b *testing.B) {
	var r ReplyBind

	for i := 0; i < b.N; i++ {
		var buf = bytes.NewBuffer(testIPv6Buf)
		_ = r.Unpack(buf, true)
		buf.Reset()
	}
}
