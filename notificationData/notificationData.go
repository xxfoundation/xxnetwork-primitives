package notificationData

import (
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"strings"
)

type NotificationData struct {
	EphemeralID int64
	RoundID     uint64
	IdentityFP  []byte
	MessageHash []byte
}

func BuildNotificationCSV(ndList []*NotificationData, maxSize int) ([]byte, []*NotificationData) {
	buf := &bytes.Buffer{}

	numWritten := 0

	for _, nd := range ndList {
		line := &bytes.Buffer{}
		w := csv.NewWriter(line)
		output := []string{base64.StdEncoding.EncodeToString(nd.MessageHash),
			base64.StdEncoding.EncodeToString(nd.IdentityFP)}

		if err := w.Write(output); err != nil {
			jww.FATAL.Printf("Failed to write notificationsCSV line: %+v", err)
		}
		w.Flush()

		if buf.Len()+line.Len() > maxSize {
			break
		}

		if _, err := buf.Write(line.Bytes()); err != nil {
			jww.FATAL.Printf("Failed to write to notificationsCSV: %+v", err)
		}

		numWritten++
	}

	return buf.Bytes(), ndList[numWritten:]
}

func DecodeNotificationsCSV(data string) ([]*NotificationData, error) {
	r := csv.NewReader(strings.NewReader(data))
	read, err := r.ReadAll()
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to decode notifications CSV")
	}

	l := make([]*NotificationData, len(read))
	for i, touple := range read {
		messageHash, err := base64.StdEncoding.DecodeString(touple[0])
		if err != nil {
			return nil, errors.WithMessage(err, "Failed decode an element")
		}
		identityFP, err := base64.StdEncoding.DecodeString(touple[1])
		if err != nil {
			return nil, errors.WithMessage(err, "Failed decode an element")
		}
		l[i] = &NotificationData{
			EphemeralID: 0,
			IdentityFP:  identityFP,
			MessageHash: messageHash,
		}
	}
	return l, nil
}
