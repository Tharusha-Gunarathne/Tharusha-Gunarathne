// Command gen fetches live GitHub stats, renders the SVG card + buttons and
// writes README.md from templates/README.tmpl. Needs GH_ACCESS_TOKEN.
package main

import (
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/Tharusha-Gunarathne/Tharusha-Gunarathne/internal/card"
	"github.com/Tharusha-Gunarathne/Tharusha-Gunarathne/internal/github"
	"github.com/Tharusha-Gunarathne/Tharusha-Gunarathne/internal/posts"
	"github.com/Tharusha-Gunarathne/Tharusha-Gunarathne/internal/profile"
	"github.com/joho/godotenv"
)

type TemplateData struct {
	Posts     []posts.Post
	ShowPosts bool
	Github    *github.Stats
	Profile   profile.Profile
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("note: no .env file loaded (%v)", err)
	}

	prof, err := profile.GetProfile()
	if err != nil {
		log.Fatal(err)
	}
	portrait, err := profile.GetPortrait()
	if err != nil {
		log.Fatal(err)
	}

	showPosts := strings.EqualFold(os.Getenv("SHOW_LATEST_POSTS"), "true")
	var blogPosts []posts.Post
	if showPosts {
		if blogPosts, err = posts.GetPosts(); err != nil {
			log.Fatalf("Error getting posts: %v", err)
		}
	}

	stats, err := github.GetGithubStats(prof.Links.GitHub.Label)
	if err != nil {
		log.Fatalf("Error getting GitHub stats: %v", err)
	}

	cardData := card.Assemble(prof, portrait, card.Stats{
		Commits:      stats.Commits,
		PullRequests: stats.PullRequests,
		Stars:        stats.Stars,
		Issues:       stats.Issues,
		Contributed:  stats.Contributed,
	})
	if err := card.WriteAll("assets", cardData); err != nil {
		log.Fatalf("Error writing profile card: %v", err)
	}
	if err := card.WriteButtons("assets"); err != nil {
		log.Fatalf("Error writing connect buttons: %v", err)
	}

	tmpl, err := template.ParseFiles("templates/README.tmpl")
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}
	out, err := os.Create("README.md")
	if err != nil {
		log.Fatalf("Error creating README.md: %v", err)
	}
	defer out.Close()

	data := TemplateData{
		Posts: blogPosts, ShowPosts: showPosts, Github: stats,
		Profile: prof,
	}
	if err := tmpl.Execute(out, data); err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
}
