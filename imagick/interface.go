////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package imagick

// Magick interface contains functions which transform images using ImageMagick
type Magick interface {
	// Reduce accepts an image in bytes and reduces its size, returning the
	// bytes of the new image.
	Reduce(image []byte) ([]byte, error)
}
