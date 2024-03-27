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
	const url = "wss://signal.teonet.dev/signal"
	const server = "server-1"
	var login string = uuid.String()
	global.Call("setIdText", "wa_login", login)

	teo := global.Get("teo")
	teo.Call("connect", url, login, server)

	const cmdClients = "clients"

	teo.Call("addReader", js.FuncOf(func(this js.Value, args []js.Value) any {
		cmd := args[0].Get("command").String()
		switch cmd {
		case cmdClients:
			global.Call("setIdText", "wa_clients", args[1])
		}
		return nil
	}))

	teo.Call("onOpen", js.FuncOf(func(this js.Value, args []js.Value) any {
		fmt.Println("webasm connected to webrtc")
		global.Call("setIdText", "wa_online", true)
		teo.Call("sendCmd", cmdClients)
		teo.Call("subscribeCmd", cmdClients)
		return nil
	}))

	teo.Call("onClose", js.FuncOf(func(this js.Value, args []js.Value) any {
		fmt.Println("webasm disconnected from webrtc")
		global.Call("setIdText", "wa_online", false)
		return nil
	}))

	// Wait forever
	select {}
}
