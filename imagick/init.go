////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

//go:build !js || !wasm

package imagick

// im when compiled normally has no internals
type im struct {
}

// getImageMagick calls the given callback with a new *im struct.
// note that if more functions are added, we may want to move the init/destroy logic to this function
func getImageMagick(setFunc func(*im)) {
	go setFunc(&im{})
}
