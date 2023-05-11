package s5

import "errors"

// General
var (
	ErrUnimplemented          = errors.New("unimplemented")
	ErrUnsupportedVersion     = errors.New("unsupported version")
	ErrUnsupportedAddressType = errors.New("unsupported address type")
	ErrUnknownNetwork         = errors.New("unknown network")
	ErrAuthNoneAcceptable     = errors.New("auth none acceptable")
	ErrInvalidHostnameLength  = errors.New("invalid hostname length")
)

//
var ErrMethodsLimit = errors.New("reached the limit of methods")
var ErrInvalidBufferSize = errors.New("invalid buffer size")

// Reply
var (
	ErrReplyGeneralFailure          = errors.New("general SOCKS server failure")
	ErrReplyConnectionNotAllowed    = errors.New("connection not allowed by ruleset")
	ErrReplyNetworkUnreachable      = errors.New("network unreachable")
	ErrReplyHostUnreachable         = errors.New("host unreachable")
	ErrReplyConnectionRefused       = errors.New("connection refused")
	ErrReplyTTLExpired              = errors.New("TTL expired")
	ErrReplyCommandNotSupported     = errors.New("command not supported")
	ErrReplyAddressTypeNotSupported = errors.New("address type not supported")
)

// Auth UserPW
var (
	ErrAuthVersion = errors.New("invalid version")
	ErrAuthUnknown = errors.New("unknown authentication method")
)
