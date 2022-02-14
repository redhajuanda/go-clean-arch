package main

import (
	"go-clean-arch/cmd"
)

// @title Go Clean Arch RESTful APIs
// @description This is a documentation for Go Clean Arch RESTful APIs. <br>

// @securityDefinitions.basic BasicAuth
// @in header
// @name Authorization

// @securityDefinitions.apikey BearerToken
// @in header
// @name Authorization

func main() {
	cmd.Run()
}
