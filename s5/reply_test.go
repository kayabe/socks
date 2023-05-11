package s5

import (
	"bytes"
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testReplyIPv4 = netip.MustParseAddr("127.0.0.1")
var testReplyIPv4Buf = []byte{VERSION, byte(ReplySuccess), 0x00, AddressTypeIPv4, 127, 0, 0, 1, 0x7, 0x5c}

func TestReplyV5Status(t *testing.T) {
	var status = ReplySuccess
	assert.Equal(t, "succeeded", status.String(), "the status should be equal")
}

func TestReplyUnpack(t *testing.T) {
	var buf = bytes.NewBuffer(testReplyIPv4Buf)
	var r Reply
	if err := r.Unpack(buf); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, Reply{
		Version:     VERSION,
		Status:      ReplySuccess,
		AddressType: AddressTypeIPv4,
		Bind: ReplyBind{
			Address: testReplyIPv4,
			Port:    1884,
		},
	}, r)
}

func TestReplyPack(t *testing.T) {
	var buf = bytes.NewBuffer(nil)
	if err := (&Reply{
		Status: ReplySuccess,
		Bind: ReplyBind{
			Address: testReplyIPv4,
			Port:    1884,
		},
	}).Pack(buf); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, testReplyIPv4Buf, buf.Bytes())
}

func BenchmarkReplyPack(b *testing.B) {
	var buf = bytes.NewBuffer(nil)

	for i := 0; i < b.N; i++ {
		_ = (&Reply{
			Status: ReplySuccess,
			Bind: ReplyBind{
				Address: testReplyIPv4,
				Port:    1884,
			},
		}).Pack(buf)
		buf.Reset()
	}
}

func BenchmarkReplyUnpack(b *testing.B) {
	data := []byte{VERSION, byte(ReplySuccess), 0x00, AddressTypeIPv4, 127, 0, 0, 1, 0xe, 0x0b}
	var r Reply

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var buf = bytes.NewBuffer(data)
		_ = r.Unpack(buf)
		buf.Reset()
	}
}
