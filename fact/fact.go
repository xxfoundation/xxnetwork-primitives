package fact

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

// marshal is for transmission for UDB, not a part of the fact interface
func (f Fact) Stringify() string {
	return f.T.Stringify() + f.Fact
}

func UnstringifyFact(s string) (Fact, error) {
	if len(s) < 1 {
		return Fact{}, errors.New("stringified facts must at least have a type at the start")
	}
	T := s[:1]
	fact := s[1:]
	ft, err := UnstringifyFactType(T)
	if err != nil {
		return Fact{}, err
	}

	return NewFact(ft, fact)
}
