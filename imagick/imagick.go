////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// the imagick package contains implementations and an interface for the
// ImageMagick library, used for image transformations and reductions.  There
// are two implementations at the moment, one using the golang library which
// targets the CGO bindings, the other using the javascript/wasm library.

package imagick

var m Magick = nil

// GetImageMagick returns a Magick interface
func GetImageMagick() Magick {
	if m == nil {
		getImageMagick(func(im *im) {
			m = im
		})
	}
	return m
}

// Reduce uses the underlying imageMagick implementation to reduce the size of an image
func (m *im) Reduce(image []byte) ([]byte, error) {
	return m.reduce(image)
}
