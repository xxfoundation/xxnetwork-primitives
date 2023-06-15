////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package imagick

import (
	"github.com/pkg/errors"
	"syscall/js"
	"time"
)

// reduce an image using the wasm/js imagemagick bindings
func (m *im) reduce(image []byte) ([]byte, error) {
	var fileInputs []interface{}
	fileInputs = append(fileInputs, map[string]interface{}{"name": "img.png", "content": copyBytesToJS(image)})
	// Define the command which will run
	// Note that the filter argument must precede the resize argument for it to apply
	command := []interface{}{"convert", "img.png", "-filter", "Point", "-resize", "40x40", "-depth", "5", "out.png"}

	i1 := js.ValueOf(fileInputs)
	i2 := js.ValueOf(command)

	reducePromise := m.jsMagick.Call("Call", i1, i2)

	reduced, err := resolve(reducePromise)
	if err != nil {
		return nil, err
	}
	blob := reduced.Index(0).Get("blob")
	size := blob.Get("size")
	consoleLog(size)

	arrayBufferPromise := blob.Call("arrayBuffer")

	arrayBuffer, err := resolve(arrayBufferPromise)
	if err != nil {
		return nil, err
	}
	uint8array := Uint8Array.New(arrayBuffer)
	return copyBytesToGo(uint8array), nil
}

func resolve(promise js.Value) (js.Value, error) {
	wait := make(chan js.Value)
	errCh := make(chan error)
	promise.Call("then", js.FuncOf(func(_ js.Value, args []js.Value) any {
		wait <- args[0]
		return nil
	}), js.FuncOf(func(_ js.Value, args []js.Value) any {
		errCh <- errors.New("Failed to resolve promise")
		return nil
	}))
	select {
	case v := <-wait:
		return v, nil
	case err := <-errCh:
		return js.Null(), err
	case <-time.Tick(time.Second * 5):
		return js.Null(), errors.New("Timed out waiting to resolve promise")
	}
}

func consoleLog(v js.Value) {
	js.Global().Get("console").Call("log", v)
}

// TODO: this is from wasmutils, should this package go elsewhere instead?

var Uint8Array = js.Global().Get("Uint8Array")

// CopyBytesToGo copies the [Uint8Array] stored in the [js.Value] to []byte.
// This is a wrapper for [js.CopyBytesToGo] to make it more convenient.
func copyBytesToGo(src js.Value) []byte {
	b := make([]byte, src.Length())
	js.CopyBytesToGo(b, src)
	return b
}

// CopyBytesToJS copies the []byte to a [Uint8Array] stored in a [js.Value].
// This is a wrapper for [js.CopyBytesToJS] to make it more convenient.
func copyBytesToJS(src []byte) js.Value {
	dst := Uint8Array.New(len(src))
	js.CopyBytesToJS(dst, src)
	return dst
}
