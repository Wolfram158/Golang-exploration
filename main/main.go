package main

import (
	"exploration/server"
	"flag"
	"log"
)

var (
	addr           = flag.String("addr", server.SettingsInstance.Addr, "http service address")
	wsEndpoint     = flag.String("ws-endpoint", server.SettingsInstance.WsEndpoint, "websocket endpoint")
	maxConnCount   = flag.Int("max-conn", server.SettingsInstance.MaxConnCount, "max count of simultaneous connections")
	secondsToSolve = flag.Int("sts", server.SettingsInstance.SecondsToSolve, "seconds given to solve one task")
)

func main() {
	flag.Parse()
	server.SettingsInstance.
		ChangeAddr(*addr).
		ChangeWsEndpoint(*wsEndpoint).
		ChangeMaxConnCount(*maxConnCount).
		ChangeSecondsToSolve(*secondsToSolve)
	log.Fatalf("%v", <-server.LaunchServer())
}
