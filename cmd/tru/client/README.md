# Teonet truproxy client

This project contains the client and server of Teonet truproxy.

## Client

Client is a static web page created with golang template engine. It displays some servers static data.

It connects to server with teonet webrtc protocol using teoweb script. And displays connection info in "This page webrtc connection" web page section.

It connects also to server from it webasm part with teoproxy client package. And displays connection info in "Webasm webrtc connection" web page section.

So this sample application connect to tha same server from web page and from webasm.

## How to use

Copy the JavaScript support file:

```bash
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ./wasm
```

To run this server use next command:

```bash
# Build wasm
cd wasm && ./build.sh && cd ..
PORT=8085 go run .
```

Then open <http://localhost:8085> in browser.

## Licence

[BSD](LICENSE)
