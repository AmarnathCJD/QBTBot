package main


import (

"github.com/cenkalti/rain/torrent"
"log"
)

var (
	client *torrent.Session
)

func init() {
	config := torrent.DefaultConfig
	client, err := torrent.NewSession(config)
	if err != nil {
		log.Fatal(err)
	}
}

