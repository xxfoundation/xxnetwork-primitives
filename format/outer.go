package format

// Describes the encryption formatting of a message
type CryptoType uint32

const (
	None CryptoType = iota
	Unencrypted
	E2E
	Garbled
	Error
	Rekey
)

// Return string representation
// of CryptoType
// Might be useful in tests and debug prints
func (ct CryptoType) String() string {
	var ret string
	switch ct {
	case None:
		ret = "None"
	case Unencrypted:
		ret = "Unencrypted"
	case E2E:
		ret = "E2E"
	case Garbled:
		ret = "Garbled"
	case Error:
		ret = "Error"
	case Rekey:
		ret = "Rekey"
	}
	return ret
}
