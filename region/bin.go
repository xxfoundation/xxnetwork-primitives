package region

import (
	"encoding/json"
	"github.com/pkg/errors"
	"math"
	"strconv"
)

// GeoBin is the numerical representation of a geographical regional.
type GeoBin uint8

const (
	Americas GeoBin = iota
	WesternEurope
	CentralEurope
	EasternEurope
	MiddleEast
	Africa
	Russia
	Asia
)

// Error messages.
const (
	invalidRegionErr = "invalid region %q"
	jsonUnmarshalErr = "could not parse JSON: %+v"
)

// String returns the string representation of the GeoBin. This functions
// satisfies the fmt.Stringer interface.
func (b GeoBin) String() string {
	switch b {
	case Americas:
		return "Americas"
	case WesternEurope:
		return "WesternEurope"
	case CentralEurope:
		return "CentralEurope"
	case EasternEurope:
		return "EasternEurope"
	case MiddleEast:
		return "MiddleEast"
	case Africa:
		return "Africa"
	case Russia:
		return "Russia"
	case Asia:
		return "Asia"
	default:
		return "INVALID BIN " + strconv.Itoa(int(b))
	}
}

// GetRegion converts the region to a numerical representation.
func GetRegion(region string) (GeoBin, error) {
	switch region {
	case "Americas":
		return Americas, nil
	case "WesternEurope", "WestEurope":
		return WesternEurope, nil
	case "CentralEurope":
		return CentralEurope, nil
	case "EasternEurope", "EastEurope":
		return EasternEurope, nil
	case "MiddleEast":
		return MiddleEast, nil
	case "Africa":
		return Africa, nil
	case "Russia":
		return Russia, nil
	case "Asia":
		return Asia, nil
	default:
		return math.MaxUint8, errors.Errorf(invalidRegionErr, region)
	}
}

// Bytes returns the byte representation of the GeoBin.
func (b GeoBin) Bytes() []byte {
	return []byte{byte(b)}
}

// MarshalJSON allows a GeoBin to be marshaled into JSON. This functions
// satisfies the json.Marshaler interface.
func (b GeoBin) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

// UnmarshalJSON and allows a GeoBin to be marshaled into JSON. This functions
// satisfies the json.Unmarshaler interface.
func (b *GeoBin) UnmarshalJSON(data []byte) error {
	var region string

	err := json.Unmarshal(data, &region)
	if err != nil {
		return errors.Errorf(jsonUnmarshalErr, err)
	}

	bin, err := GetRegion(region)
	if err != nil {
		return err
	}

	*b = bin

	return nil
}
