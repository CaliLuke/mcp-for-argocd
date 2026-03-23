package main

import (
	"log"
	"os"

	"github.com/argoproj-labs/mcp-for-argocd/internal/mcpserver"
)

func main() {
	if err := mcpserver.Run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
