package socks

import (
	"io"
	"net"
	"net/netip"
	"strconv"
	"unsafe"

	s5 "github.com/kayabe/socks/s5"
)

type Authentication interface {
	Pack(io.Writer) error
	Unpack(io.Reader) error

	// Returns authentication method
	Method() uint8
}

type Client struct {
	net.Dialer
	ProtocolVersion any
	Authentication  Authentication
	ProxyAddr       string
}

// NewClient creates a new SOCKS client.
func NewClient(proxyAddr string, proxyVersion any, options ...func(*Client)) (client *Client, err error) {
	client = &Client{
		ProtocolVersion: proxyVersion,
		ProxyAddr:       proxyAddr,
	}
	for _, o := range options {
		o(client)
	}
	return
}

// DialTCP connects to the given TCPAddr via the proxy server.
func (c *Client) DialTCP(network string, laddr, raddr *net.TCPAddr) (conn *net.TCPConn, err error) {
	proxyAddr, err2 := net.ResolveTCPAddr("tcp", c.ProxyAddr)
	if err2 != nil {
		return nil, err2
	}
	if conn, err = net.DialTCP(network, laddr, proxyAddr); err != nil {
		return
	}
	switch c.ProtocolVersion.(type) {
	case ProtocolV5:
		if err = c.HandshakeV5(conn); err != nil {
			return
		}
		if _, err = c.ConnectV5(conn, raddr); err != nil {
			return
		}
	default:
		return nil, s5.ErrUnsupportedVersion
	}
	return conn, nil
}

// Dial connects to the given address via the proxy server.
func (c *Client) Dial(network string, address string) (conn net.Conn, err error) {
	switch network {
	case "udp", "udp4", "udp6":
		return nil, s5.ErrUnimplemented
	}

	if conn, err = c.Dialer.Dial(network, c.ProxyAddr); err != nil {
		return
	}

	switch c.ProtocolVersion.(type) {
	case ProtocolV5:
		if err = c.HandshakeV5(conn); err != nil {
			return
		}
		if _, err = c.ConnectV5(conn, address); err != nil {
			return
		}
	default:
		return nil, s5.ErrUnsupportedVersion
	}

	return conn, nil
}

func (c *Client) HandshakeV5(conn net.Conn) (err error) {
	defer func() {
		if err != nil {
			conn.Close()
		}
	}()

	var handshake = new(s5.HandshakeRequest)

	if c.Authentication == nil {
		handshake.Methods = []s5.AuthMethod{s5.MethodAuthNone}
	} else {
		handshake.Methods = []s5.AuthMethod{s5.MethodAuthNone, s5.AuthMethod(c.Authentication.Method())}
	}

	if err = handshake.Pack(conn); err != nil {
		return
	}

	var reply s5.HandshakeReply

	if err = reply.Unpack(conn); err != nil {
		return
	}

	switch reply.Method {
	case s5.MethodAuthNoneAcceptable:
		return s5.ErrAuthNoneAcceptable
	case s5.MethodAuthUserPW:
		if c.Authentication != nil {
			if err = c.Authentication.Pack(conn); err != nil {
				return
			}

			if _, err = s5.AuthReplyFromConn(conn); err != nil {
				return
			}
		}
	}

	return nil
}

// ConnectV5 target can be used as string for both ip and domain or *net.TCPAddr for ip:port only
func (c *Client) ConnectV5(conn net.Conn, target any) (reply *s5.Reply, err error) {
	var req = &s5.Request{
		Version: s5.VERSION,
		Command: s5.CommandConnect,
	}

	var ap netip.AddrPort

	switch target := target.(type) {
	case string:
		if tc, err := netip.ParseAddrPort(target); err == nil {
			ap = tc
			break
		}

		var hostname, port string
		var castPort uint64

		if hostname, port, err = net.SplitHostPort(target); err != nil {
			break
		}

		if castPort, err = strconv.ParseUint(port, 10, 16); err != nil {
			break
		}

		req.AddressType = s5.AddressTypeDomainName
		req.Destination = &s5.RequestV5DestDomainName{Address: unsafe.Slice(unsafe.StringData(hostname), len(hostname)), Port: uint16(castPort)}
	case *net.TCPAddr:
		ap = target.AddrPort()
	case *net.UDPAddr:
		ap = target.AddrPort()
	}

	if ap.IsValid() {
		if ap.Addr().Is4() {
			req.AddressType = s5.AddressTypeIPv4
			req.Destination = &s5.RequestV5DestIPv4{Address: ap.Addr().As4(), Port: ap.Port()}
		} else if ap.Addr().Is6() {
			req.AddressType = s5.AddressTypeIPv6
			req.Destination = &s5.RequestV5DestIPv6{Address: ap.Addr().As16(), Port: ap.Port()}
		}
	}

	if req.Destination == nil {
		return nil, s5.ErrUnsupportedAddressType
	}

	if err = req.Pack(conn); err != nil {
		return
	}

	reply = new(s5.Reply)
	if err = reply.Unpack(conn); err != nil {
		return nil, err
	} else if reply.Status != s5.ReplySuccess {
		return nil, reply.Status.Error()
	}
	return
}
