package main

import (
	"go-clean-arch/cmd"
)

// @title Motion CRM RESTful APIs
// @description This is a documentation for Motion CRM RESTful APIs. <br>

// @securityDefinitions.basic BasicAuth
// @in header
// @name Authorization

// @securityDefinitions.apikey BearerToken
// @in header
// @name Authorization

func main() {
	cmd.Run()
}
