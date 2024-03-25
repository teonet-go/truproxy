package main

import (
	"fmt"
)

func main() {
	// Logo prompt to console when wasm is ready
	fmt.Println("hello from wasm")

	// Define global js variable
	// global := js.Global()

	// Basic teonet crypt functions
	// global.Set("teoHashKey", js.FuncOf(hashKey))
	// global.Set("teoEncrypt", js.FuncOf(encrypt))
	// global.Set("teoDecrypt", js.FuncOf(decrypt))

	// Wait forever
	select {}
}
