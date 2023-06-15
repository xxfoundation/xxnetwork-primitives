////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

//go:build !js || !wasm

package imagick

import (
	"gopkg.in/gographics/imagick.v3/imagick"
)

// reduce an image using the golang CGO bindings of imagemagick
func (m *im) reduce(image []byte) ([]byte, error) {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()

	err := mw.ReadImageBlob(image)
	if err != nil {
		return nil, err
	}

	err = mw.ResizeImage(40, 40, imagick.FILTER_POINT)
	if err != nil {
		return nil, err
	}

	err = mw.SetImageDepth(5)
	if err != nil {
		return nil, err
	}

	return mw.GetImageBlob(), nil
}
