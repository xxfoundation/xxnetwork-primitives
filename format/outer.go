package format

// Describes the encryption formatting of a message
type OuterType uint32

const (
	Unencrypted OuterType = iota
	E2E
	Garbled
	Error
	RekeyTrigger
	Rekey
)
