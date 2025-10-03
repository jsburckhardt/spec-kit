// Package ui provides terminal user interface components
package ui

import (
	"strings"

	"github.com/jsburckhardt/spec-kit/gospecify/internal/config"
)

// ShowBanner displays the application banner
func ShowBanner() {
	banner := BannerStyle.Render(config.Banner)
	tagline := TaglineStyle.Render(config.Tagline)

	// Print with proper spacing
	println(banner)
	println()
	println(tagline)
	println()
}

// GetBanner returns the formatted banner as a string
func GetBanner() string {
	var output strings.Builder
	output.WriteString(BannerStyle.Render(config.Banner))
	output.WriteString("\n\n")
	output.WriteString(TaglineStyle.Render(config.Tagline))
	output.WriteString("\n\n")
	return output.String()
}

// ShowMiniBanner displays a compact version of the banner
func ShowMiniBanner() {
	miniBanner := `
╔═╗╔═╗╔═╗╔═╗╦╔═╗╦ ╦
╚═╗╠═╝║╣ ║  ║╠╣ ╚╦╝
╚═╝╩  ╚═╝╚═╝╩╚   ╩ `

	println(CyanStyle.Render(miniBanner))
}
