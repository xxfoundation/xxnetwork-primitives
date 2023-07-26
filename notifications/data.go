////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package notifications

import (
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
)

type Data struct {
	EphemeralID int64
	RoundID     uint64
	IdentityFP  []byte
	MessageHash []byte
}

func (d *Data) String() string {
	fields := []string{
		strconv.FormatInt(d.EphemeralID, 10),
		strconv.FormatUint(d.RoundID, 10),
		base64.StdEncoding.EncodeToString(d.IdentityFP),
		base64.StdEncoding.EncodeToString(d.MessageHash),
	}
	return "{" + strings.Join(fields, " ") + "}"
}

// BuildNotificationCSV converts the [Data] list into a CSV of the specified max
// size and return it along with the included [Data] entries. Any [Data] entries
// over that size are excluded.
//
// The CSV contains each [Data] entry on its own row with column one the
// [Data.MessageHash] and column two having the [Data.IdentityFP], but base 64
// encoded
func BuildNotificationCSV(ndList []*Data, maxSize int) ([]byte, []*Data) {
	var buf bytes.Buffer
	var numWritten int

	for i, nd := range ndList {
		var line bytes.Buffer
		w := csv.NewWriter(&line)
		output := []string{
			base64.StdEncoding.EncodeToString(nd.MessageHash),
			base64.StdEncoding.EncodeToString(nd.IdentityFP)}

		if err := w.Write(output); err != nil {
			jww.FATAL.Printf("Failed to write record %d of %d to "+
				"notifications CSV line buffer: %+v", i, len(ndList), err)
		}
		w.Flush()

		if buf.Len()+line.Len() > maxSize {
			break
		}

		if _, err := buf.Write(line.Bytes()); err != nil {
			jww.FATAL.Printf("Failed to write record %d of %d to "+
				"notifications CSV: %+v", i, len(ndList), err)
		}

		numWritten++
	}

	return buf.Bytes(), ndList[numWritten:]
}

// DecodeNotificationsCSV decodes the Data list CSV into a slice of Data.
func DecodeNotificationsCSV(data string) ([]*Data, error) {
	r := csv.NewReader(strings.NewReader(data))
	records, err := r.ReadAll()
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read notifications CSV records.")
	}

	list := make([]*Data, len(records))
	for i, tuple := range records {
		messageHash, err := base64.StdEncoding.DecodeString(tuple[0])
		if err != nil {
			return nil, errors.Wrapf(err,
				"Failed to decode MessageHash for record %d of %d",
				i, len(records))
		}

		identityFP, err := base64.StdEncoding.DecodeString(tuple[1])
		if err != nil {
			return nil, errors.Wrapf(err,
				"Failed to decode IdentityFP for record %d of %d",
				i, len(records))
		}
		list[i] = &Data{
			IdentityFP:  identityFP,
			MessageHash: messageHash,
		}
	}

	return list, nil
}
