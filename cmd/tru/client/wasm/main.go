package main

import (
	"fmt"
	"syscall/js"
)

func main() {
	// Logo prompt to console when wasm is ready
	fmt.Println("hello from wasm")

	// Define global js variable
	global := js.Global()

	// Get uuid from js module
	uuid := global.Call("uuidv4")
	fmt.Println(uuid)

	// Connect to Teonet WebRTC server
	// var teoweb = global.Get("teoweb")
	const url = "wss://signal.teonet.dev/signal"
	const server = "server-1"
	var login string = uuid.String()
	global.Call("setIdText", "wa_login", login) // setIdText("login", login);
	// let teo = new teoweb();
	// teo.connect(url, login, server);
	teo := global.Get("teo")
	teo.Call("connect", url, login, server)

	teo.Call("onOpen", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("webasm connected to webrtc")
		global.Call("setIdText", "wa_online", true)
		return nil
	}))

	const cmdClients = "clients"

	teo.Call("onOpen", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("webasm connected to webrtc")
		global.Call("setIdText", "wa_online", true)
		teo.Call("sendCmd", cmdClients)
		return nil
	}))

	teo.Call("onClose", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("webasm disconnected from webrtc")
		global.Call("setIdText", "wa_online", false)
		return nil
	}))

	// Set functions
	// global.Set("teoHashKey", js.FuncOf(hashKey))
	// global.Set("teoEncrypt", js.FuncOf(encrypt))
	// global.Set("teoDecrypt", js.FuncOf(decrypt))

	// Wait forever
	select {}
}