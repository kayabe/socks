package s5

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestV5HostnamePack(t *testing.T) {
	r := &RequestV5DestDomainName{Address: []byte("google.com"), Port: 80}
	buf := make([]byte, r.Size())
	if err := r.Put(buf); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, []byte{0x0A, 'g', 'o', 'o', 'g', 'l', 'e', '.', 'c', 'o', 'm', 0x00, 0x50}, buf, "the destination should be equal")

}

func TestRequestV5IPv4Pack(t *testing.T) {
	r := RequestV5DestIPv4{Address: [4]byte{127, 0, 0, 1}, Port: 80}
	buf := make([]byte, r.Size())

	if err := r.Put(buf); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, []byte{0x7f, 0x00, 0x00, 0x01, 0x00, 0x50}, buf, "the destination should be equal")
}

func TestRequestV5IPv6Pack(t *testing.T) {
	r := RequestV5DestIPv6{Address: [16]byte{0xe, 0xe}, Port: 80}
	buf := make([]byte, r.Size())

	if err := r.Put(buf); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, []byte{0x0e, 0x0e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x50}, buf, "the destination should be equal")
}

func TestRequestV5Pack(t *testing.T) {
	var buf = bytes.NewBuffer(nil)

	if err := (&Request{
		Command:     CommandConnect,
		Destination: &RequestV5DestDomainName{Address: []byte("google.com"), Port: 80},
	}).Pack(buf); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, []byte{0x5, 0x1, 0x0, 0x3, 0xa, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x0, 0x50}, buf.Bytes())

	buf.Reset()

	if err := (&Request{
		Command:     CommandConnect,
		Destination: &RequestV5DestIPv4{Address: [4]byte{127, 0, 0, 1}, Port: 80},
	}).Pack(buf); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, []byte{0x5, 0x1, 0x0, 0x1, 0x7f, 0x0, 0x0, 0x1, 0x0, 0x50}, buf.Bytes())

	buf.Reset()

	if err := (&Request{
		Command:     CommandConnect,
		Destination: &RequestV5DestIPv6{Address: [16]byte{0xe, 0xe}, Port: 80},
	}).Pack(buf); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, []byte{0x5, 0x1, 0x0, 0x4, 0x0e, 0x0e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x50}, buf.Bytes())
}

func TestRequestV5Unpack(t *testing.T) {
	var (
		bufAddrDomainName = bytes.NewBuffer([]byte{0x05, 0x01, 0x00, 0x03, 0x0A, 'g', 'o', 'o', 'g', 'l', 'e', '.', 'c', 'o', 'm', 0x00, 0x50})
		bufAddrIPv4       = bytes.NewBuffer([]byte{0x05, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x50})
		bufAddrIPv6       = bytes.NewBuffer([]byte{0x05, 0x01, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x50})
	)

	//

	var req1 Request
	if err := req1.Unpack(bufAddrDomainName); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, Request{
		Version:     VERSION,
		Command:     CommandConnect,
		AddressType: AddressTypeDomainName,
		Destination: &RequestV5DestDomainName{Address: []byte("google.com"), Port: 80},
	}, req1, "the request should be equal")

	//

	var req2 Request
	if err := req2.Unpack(bufAddrIPv4); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, Request{
		Version:     VERSION,
		Command:     CommandConnect,
		AddressType: AddressTypeIPv4,
		Destination: &RequestV5DestIPv4{Address: [4]byte{}, Port: 80},
	}, req2, "the request should be equal")

	//

	var req3 Request
	if err := req3.Unpack(bufAddrIPv6); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, Request{
		Version:     VERSION,
		Command:     CommandConnect,
		AddressType: AddressTypeIPv6,
		Destination: &RequestV5DestIPv6{Address: [16]byte{}, Port: 80},
	}, req3, "the request should be equal")
}

///

func BenchmarkDestDomainPack(b *testing.B) {
	var buf = make([]byte, 10+3)
	var r = RequestV5DestDomainName{Address: []byte("google.com"), Port: 80}

	for i := 0; i < b.N; i++ {
		r.Put(buf)
	}
}

func BenchmarkDestDomainUnpack(b *testing.B) {
	var data = []byte{0x0A, 'g', 'o', 'o', 'g', 'l', 'e', '.', 'c', 'o', 'm', 0x00, 0x50}
	var r RequestV5DestDomainName
	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(data)
		r.Unpack(buf)
		buf.Reset()
	}
}

func BenchmarkRequestPack(b *testing.B) {
	buf := bytes.NewBuffer(nil)
	var r = Request{
		Version:     VERSION,
		Command:     CommandConnect,
		Destination: &RequestV5DestDomainName{Address: []byte("google.com"), Port: 80},
	}
	for i := 0; i < b.N; i++ {
		r.Pack(buf)
		buf.Reset()
	}
}

func BenchmarkRequestUnpack(b *testing.B) {
	data := []byte{0x05, 0x01, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x50}
	var r Request
	for i := 0; i < b.N; i++ {
		var buf = bytes.NewBuffer(data)
		r.Unpack(buf)
		buf.Reset()
	}
}
