package packer

import "io"

type Packer interface {
	// Pack writes the structure to the given writer as bytes.
	Pack(w io.Writer) error

	// Unpack reads from the given reader into the structure.
	Unpack(r io.Reader) error
}
