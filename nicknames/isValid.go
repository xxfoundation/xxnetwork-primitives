package nicknames

import (
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
)

const (
	MinNicknameLength = 3
	MaxNicknameLength = 24
)

var ErrNicknameTooShort = errors.Errorf("nicknames must be at least "+
	"%d characters in length", MinNicknameLength)
var ErrNicknameTooLong = errors.Errorf("nicknames must be %d "+
	"characters in length or less", MaxNicknameLength)

// IsValid checks if a nickname is valid.
//
// Rules:
//   - A nickname must not be longer than 24 characters.
//   - A nickname must not be shorter than 1 character.
//   - If a nickname is blank (empty string), then it will be treated by the
//     system as no nickname.
//
// TODO: Add character filtering.
func IsValid(nick string) error {
	if nick == "" {
		jww.INFO.Printf(
			"Empty nickname passed; treating it as if no nickname was set.")
		return nil
	}

	runeNick := []rune(nick)
	if len(runeNick) < MinNicknameLength {
		return errors.WithStack(ErrNicknameTooShort)
	}

	if len(runeNick) > MaxNicknameLength {
		return errors.WithStack(ErrNicknameTooLong)
	}

	return nil
}
