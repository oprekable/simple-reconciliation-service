package main

import (
	"embed"
	"simple-reconciliation-service/cmd"
)

//go:embed all:embeds
var embedFS embed.FS

func main() {
	cmd.Execute(&embedFS)
}
