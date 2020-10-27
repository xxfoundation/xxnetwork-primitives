package contact

import "errors"

type Fact struct {
	Fact string
	T    FactType
}

func NewFact(ft FactType, fact string) (Fact, error) {
	//todo: filter the fact string
	return Fact{
		Fact: fact,
		T:    ft,
	}, nil
}

func (f Fact) Get() string {
	return f.Fact
}

func (f Fact) Type() int {
	return int(f.T)
}

// marshal is for transmission for UDB, not a part of the fact interface
func (f Fact) Stringify() string {
	return f.T.Stringify() + f.Fact
}

func UnstringifyFact(s string) (Fact, error) {
	ft, err := UnstringifyFactType(s)
	if err != nil {
		return Fact{}, err
	}
	if len(s) < 1 {
		return Fact{}, errors.New("cannot unstringify a fact that's just its type")
	}

	return NewFact(ft, s[1:])
}
