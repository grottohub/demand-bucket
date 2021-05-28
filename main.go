package main

import localserver "demand-bucket/local-server"

func main() {
	localserver.Start()
	localserver.Status()
}
