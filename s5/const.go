package s5

const VERSION uint8 = 0x05

const (
	CommandConnect uint8 = iota + 1
	CommandBind
	CommandAssociate
)

const (
	AddressTypeIPv4 uint8 = iota + 1
	_
	AddressTypeDomainName
	AddressTypeIPv6
)
