package format

// Describes the encryption formatting of a message
type CryptoType uint32

const (
	None CryptoType = iota
	Unencrypted
	E2E
	Garbled
	Error
	RekeyTrigger
	Rekey
)
