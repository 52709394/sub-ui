package main

import (
	"sub-ui/serve"
	"sub-ui/setup"
)

func main() {
	serve.ToggleContent = ""
	setup.GetData()
	s := serve.Server{}
	s.Run()
}
