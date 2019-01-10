package main

import (
	"github.com/fatih/color"
)

var (
	WhiteColor    = color.FgHiWhite
	GreenColor    = color.FgHiGreen
	RedColor      = color.FgHiRed
	CyanColor     = color.FgHiCyan
	HiYellowColor = color.FgHiYellow
	YellowColor   = color.FgYellow
	MagentaColor  = color.FgMagenta
	HiWhiteColor  = color.FgHiWhite
	FaintColor    = color.Faint
)

func ColorOutput(msg string, gColor color.Attribute) {
	x := color.New(gColor)
	x.Fprintf(color.Output, "%s\n", msg)
}
