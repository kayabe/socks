package s5

// ReplyStatus ...
//
//	X'09' to X'FF' unassigned
type ReplyStatus uint8

// Error ...
func (s ReplyStatus) Error() error {
	switch s {
	case ReplyGeneralFailure:
		return ErrReplyGeneralFailure
	case ReplyConnectionNotAllowed:
		return ErrReplyConnectionNotAllowed
	case ReplyNetworkUnreachable:
		return ErrReplyNetworkUnreachable
	case ReplyHostUnreachable:
		return ErrReplyHostUnreachable
	case ReplyConnectionRefused:
		return ErrReplyConnectionRefused
	case ReplyTTLExpired:
		return ErrReplyTTLExpired
	case ReplyCommandNotSupported:
		return ErrReplyCommandNotSupported
	case ReplyAddressTypeNotSupported:
		return ErrReplyAddressTypeNotSupported
	}
	return nil
}

const (
	ReplySuccess                 ReplyStatus = iota // succeeded
	ReplyGeneralFailure                             // general SOCKS server failure
	ReplyConnectionNotAllowed                       // connection not allowed by ruleset
	ReplyNetworkUnreachable                         // Network unreachable
	ReplyHostUnreachable                            // Host unreachable
	ReplyConnectionRefused                          // Connection refused
	ReplyTTLExpired                                 // TTL expired
	ReplyCommandNotSupported                        // Command not supported
	ReplyAddressTypeNotSupported                    // Address type not supported
)

//go:generate stringer -type=ReplyStatus -linecomment -output reply_status_string.go
