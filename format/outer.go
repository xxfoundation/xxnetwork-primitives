package format

// An buffer containing the data from the associated data
type Fingerprint [AD_KEYFP_LEN]byte

// Describes the encryption stage of a message
const (
	Unecnrypted = iota
	E2E
	Garbled
	Error
	RekeyTrigger
	Rekey
)
