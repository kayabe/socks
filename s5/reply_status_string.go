// Code generated by "stringer -type=ReplyStatus -linecomment -output reply_status_string.go"; DO NOT EDIT.

package s5

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ReplySuccess-0]
	_ = x[ReplyGeneralFailure-1]
	_ = x[ReplyConnectionNotAllowed-2]
	_ = x[ReplyNetworkUnreachable-3]
	_ = x[ReplyHostUnreachable-4]
	_ = x[ReplyConnectionRefused-5]
	_ = x[ReplyTTLExpired-6]
	_ = x[ReplyCommandNotSupported-7]
	_ = x[ReplyAddressTypeNotSupported-8]
}

const _ReplyStatus_name = "succeededgeneral SOCKS server failureconnection not allowed by rulesetNetwork unreachableHost unreachableConnection refusedTTL expiredCommand not supportedAddress type not supported"

var _ReplyStatus_index = [...]uint8{0, 9, 37, 70, 89, 105, 123, 134, 155, 181}

func (i ReplyStatus) String() string {
	if i >= ReplyStatus(len(_ReplyStatus_index)-1) {
		return "ReplyStatus(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ReplyStatus_name[_ReplyStatus_index[i]:_ReplyStatus_index[i+1]]
}
