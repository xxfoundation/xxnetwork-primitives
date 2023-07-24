////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package region

import (
	"encoding/json"
	"math"
	"strconv"

	"github.com/pkg/errors"
)

// GeoBin is the numerical representation of a geographical regional.
type GeoBin uint8

const (
	NorthAmerica GeoBin = iota
	SouthAndCentralAmerica
	WesternEurope
	CentralEurope
	EasternEurope
	MiddleEast
	NorthernAfrica
	SouthernAfrica
	Russia
	EasternAsia
	WesternAsia
	Oceania
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
	case NorthAmerica:
		return "NorthAmerica"
	case SouthAndCentralAmerica:
		return "SouthAndCentralAmerica"
	case WesternEurope:
		return "WesternEurope"
	case CentralEurope:
		return "CentralEurope"
	case EasternEurope:
		return "EasternEurope"
	case MiddleEast:
		return "MiddleEast"
	case NorthernAfrica:
		return "NorthernAfrica"
	case SouthernAfrica:
		return "SouthernAfrica"
	case Russia:
		return "Russia"
	case EasternAsia:
		return "EasternAsia"
	case WesternAsia:
		return "WesternAsia"
	case Oceania:
		return "Oceania"
	default:
		return "INVALID BIN " + strconv.Itoa(int(b))
	}
}

// GetRegion converts the region to a numerical representation.
func GetRegion(region string) (GeoBin, error) {
	switch region {
	case "NorthAmerica":
		return NorthAmerica, nil
	case "SouthAndCentralAmerica":
		return SouthAndCentralAmerica, nil
	case "WesternEurope":
		return WesternEurope, nil
	case "CentralEurope":
		return CentralEurope, nil
	case "EasternEurope":
		return EasternEurope, nil
	case "MiddleEast":
		return MiddleEast, nil
	case "NorthernAfrica":
		return NorthernAfrica, nil
	case "SouthernAfrica":
		return SouthernAfrica, nil
	case "Russia":
		return Russia, nil
	case "EasternAsia":
		return EasternAsia, nil
	case "WesternAsia":
		return WesternAsia, nil
	case "Oceania":
		return Oceania, nil
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
