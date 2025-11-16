package main

import (
	"io"
	"log"
	"os"

	"github.com/personalconnect/dragpass-keeper/internal/keystore"
)

func main() {
	log.SetOutput(os.Stderr)
	log.Println("dragpass extension helper started")

	for {
		req, err := keystore.ReadMessage()
		if err != nil {
			if err == io.EOF {
				log.Println("lost chrome extension connection")
				break
			}
			keystore.SendResponse(keystore.ResponseMessage{Success: false, Error: "wrong message format: " + err.Error()})
			continue
		}
		keystore.LogRequest(req)
		resp := keystore.HandleRequest(req)
		keystore.SendResponse(resp)
	}
}
