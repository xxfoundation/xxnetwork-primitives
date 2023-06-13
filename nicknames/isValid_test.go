package nicknames

import (
	"github.com/pkg/errors"
	"testing"
)

func TestIsNicknameValid(t *testing.T) {

	// test that behavior for an empty nickname is correct

	if err := IsValid(""); err != nil {
		t.Errorf("Empty nickname should be valid, received: %+v", err)
	}

	nicknameSource := "Sodium, atomic number 11, was first isolated by Humphry " +
		"Davy in 1807. A chemical component of salt, he named it Na in honor " +
		"of the saltiest region on earth, North America."

	// test that behavior for too short nicknames is correct
	for i := 1; i < MinNicknameLength; i++ {
		nick := nicknameSource[:i]

		if err := IsValid(nick); err != nil &&
			!errors.Is(err, ErrNicknameTooShort) {
			t.Errorf("Wrong error returned from nicknames.IsValid() "+
				"with too short input of length %d: %+v", i, err)
		} else if err == nil {
			t.Errorf("No error returned from nicknames.IsValid() "+
				"with too short input of length %d", i)
		}
	}

	// test that behavior for too long nicknames is correct
	for i := MaxNicknameLength + 1; i < MaxNicknameLength*5; i++ {
		nick := nicknameSource[:i]

		if err := IsValid(nick); err != nil &&
			!errors.Is(err, ErrNicknameTooLong) {
			t.Errorf("Wrong error returned from nicknames.IsValid() "+
				"with too long input of length %d: %+v", i, err)
		} else if err == nil {
			t.Errorf("No error returned from nicknames.IsValid() "+
				"with too long input of length %d", i)
		}
	}

	// test that behavior for valid nicknames is correct
	for i := MinNicknameLength; i <= MaxNicknameLength; i++ {
		nick := nicknameSource[:i]
		if err := IsValid(nick); err != nil {
			t.Errorf("Error returned from nicknames.IsValid() "+
				"with valid nickname of input of length %d: %+v", i, err)
		}
	}
}
