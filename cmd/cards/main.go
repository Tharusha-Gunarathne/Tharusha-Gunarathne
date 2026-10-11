// Command cards regenerates the profile SVGs offline (no GitHub token needed),
// using the real profile data and placeholder stats. Use it to preview design
// changes locally: `go run ./cmd/cards`.
package main

import (
	"log"

	"github.com/Tharusha-Gunarathne/Tharusha-Gunarathne/internal/card"
	"github.com/Tharusha-Gunarathne/Tharusha-Gunarathne/internal/profile"
)

func main() {
	p, err := profile.GetProfile()
	if err != nil {
		log.Fatal(err)
	}
	portrait, err := profile.GetPortrait()
	if err != nil {
		log.Fatal(err)
	}
	// placeholder stats (the weekly `cmd/gen` run overwrites with live values)
	stats := card.Stats{Commits: 420, PullRequests: 37, Stars: 12}

	d := card.Assemble(p, portrait, stats)

	if err := card.WriteAll("assets", d); err != nil {
		log.Fatalf("write cards: %v", err)
	}
	if err := card.WriteButtons("assets"); err != nil {
		log.Fatalf("write buttons: %v", err)
	}
	log.Println("wrote assets/{dark,light}_mode.svg and connect buttons")
}
