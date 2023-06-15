////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package imagick

import (
	"fmt"
	"syscall/js"
)

// im when compiled for javascript contains a js.Value pointing to the imported
// Magick object in js
type im struct {
	jsMagick js.Value
}

// getImageMagick retrieves an imagick js object from an import call, then
// calls the given callback with a new *im struct.
func getImageMagick(setFunc func(*im)) {
	window := js.Global().Get("window")
	imagickPromise := window.Call("eval", "import('https://knicknic.github.io/wasm-imagemagick/magickApi.js')")
	imagickPromise.Call("then", js.FuncOf(func(_ js.Value, args []js.Value) any {
		fmt.Println("Successful resolution of imagemagick promise")
		consoleLog(args[0])
		setFunc(&im{jsMagick: args[0]})
		return nil
	}), js.FuncOf(func(_ js.Value, args []js.Value) any {
		fmt.Println("Unsuccessful resolution of imagemagick promise")
		return nil
	}))
}
