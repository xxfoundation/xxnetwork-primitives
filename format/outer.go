package format

// Describes the encryption formatting of a message
type OuterType uint32

const (
	None OuterType = iota
	Unencrypted
	E2E
	Garbled
	Error
	RekeyTrigger
	Rekey
)
