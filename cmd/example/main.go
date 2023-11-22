package main

import (
	"fmt"
	"github.com/sagernet/cronet-go"
	"io"
	"log"
	"net/http"
)

func main() {
	engineParams := cronet.NewEngineParams()
	engineParams.SetEnableHTTP2(true)
	engineParams.SetEnableBrotli(true)
	engineParams.SetEnableQuic(false)
	defer engineParams.Destroy()

	engine := cronet.NewEngine()
	result := engine.StartWithParams(engineParams)
	if result != cronet.ResultSuccess {
		log.Fatalln(fmt.Errorf("cronet: failed to start engine. result code: %d", result))
	}

	log.Printf("started engine version %s\n", engine.Version())

	rt := cronet.RoundTripper{
		CheckRedirect: func(newLocationUrl string) bool {
			return true
		},
		Engine: engine,
	}

	req, err := http.NewRequest("GET", "https://tls.peet.ws/api/all", nil)
	if err != nil {
		log.Fatalln(err)
	}

	res, err := rt.RoundTrip(req)
	if err != nil {
		log.Fatalln(err)
	}

	body, _ := io.ReadAll(res.Body)

	log.Println(res.StatusCode)
	log.Println(string(body))
}
