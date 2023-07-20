package nicknames

import (
	"testing"

	"github.com/pkg/errors"
)

const nicknameSource = "Sodium, atomic number 11, was first isolated by " +
	"Humphry Davy in 1807. A chemical component of salt, he named it Na " +
	"in honor of the saltiest region on earth, North America."

// Tests that IsValid returns true for all usernames within the correct lengths.
func TestIsValid(t *testing.T) {
	for i := MinNicknameLength; i <= MaxNicknameLength; i++ {
		nick := nicknameSource[:i]
		if err := IsValid(nick); err != nil {
			t.Errorf("Error returned from nicknames.IsValid with valid "+
				"nickname %q of input of length %d: %+v", nick, i, err)
		}
	}
}

// Tests that IsValid return nil for an empty nickname.
func TestIsValid_Empty(t *testing.T) {
	if err := IsValid(""); err != nil {
		t.Errorf("Empty nickname should be valid, received: %+v", err)
	}
}

// Error path: Tests that IsValid returns the error ErrNicknameTooLong when the
// nickname is too long.
func TestIsValid_MaxLengthError(t *testing.T) {
	for i := MaxNicknameLength + 1; i < MaxNicknameLength*5; i++ {
		nick := nicknameSource[:i]
		err := IsValid(nick)
		if err == nil || !errors.Is(err, ErrNicknameTooLong) {
			t.Errorf("Wrong error returned from nicknames.IsValid with too "+
				"long input of length %d.\nexpected: %v\nreceived: %+v",
				i, ErrNicknameTooLong, err)
		}
	}
}

// Error path: Tests that IsValid returns the error ErrNicknameTooShort when the
// nickname is too short.
func TestIsValid_MinLengthError(t *testing.T) {
	for i := 1; i < MinNicknameLength; i++ {
		nick := nicknameSource[:i]

		err := IsValid(nick)
		if err == nil || !errors.Is(err, ErrNicknameTooShort) {
			t.Errorf("Wrong error returned from nicknames.IsValid with too "+
				"short input of length %d.\nexpected: %v\nreceived: %+v",
				i, ErrNicknameTooShort, err)
		}
	}
}
