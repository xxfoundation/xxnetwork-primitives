////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package imagick

import (
	"encoding/base64"
	"fmt"
	"testing"
)

// TestGetImageMagick can be targeted at either js or native to test the
// initialization and reduction functions.
func TestGetImageMagick(t *testing.T) {
	var imv *im
	wait := make(chan interface{})
	getImageMagick(func(i *im) {
		imv = i
		wait <- nil
	})
	<-wait
	t.Log(imv)

	imgBytes, err := base64.StdEncoding.DecodeString(imgB64)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(len(imgBytes))

	reduced, err := imv.Reduce(imgBytes)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(len(reduced))
	t.Log(base64.StdEncoding.EncodeToString(reduced))
}
