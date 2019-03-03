package main

import (
	"github.com/fatih/color"
)

var (
	// WhiteColor
	// WhiteColor = color.FgHiWhite
	// GreenColor
	GreenColor = color.FgHiGreen
	// RedColor
	RedColor = color.FgHiRed
	// CyanColor
	CyanColor = color.FgHiCyan
	// HiYellowColor
	// HiYellowColor = color.FgHiYellow
	// YellowColor
	YellowColor = color.FgYellow
	// MagentaColor
	MagentaColor = color.FgMagenta
	// HiWhiteColor
	// HiWhiteColor = color.FgHiWhite
	// FaintColor
	// FaintColor = color.Faint
)

// ColorOutput prints the passed string in the passed color
func ColorOutput(msg string, gColor color.Attribute) {
	x := color.New(gColor)
	x.Fprintf(color.Output, "%s\n", msg)
}
