package styles

import "github.com/charmbracelet/lipgloss"

var LabelStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF")).Underline(true)
var LabelStyleRender = LabelStyle.Render
var SubLabelStyle = LabelStyle.Copy().Bold(false).Italic(true).Underline(false).Render
var ErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ab2727")).Render
var TitleStyle = lipgloss.NewStyle().Margin(1).Padding(0, 2).Align(lipgloss.Center).Bold(true).Foreground(lipgloss.Color("228")).BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("63")).Render
var TabStyle = lipgloss.NewStyle().Padding(0, 1).Margin(0, 0, 1).Align(lipgloss.Center).BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("63")).BorderBottom(false).BorderTop(true).BorderLeft(true).BorderRight(true)
var ActivatedTabStyle = TabStyle.Copy().Bold(true).Foreground(lipgloss.Color("#FFFFFF")).Render
var DeactivatedTabStyle = TabStyle.Copy().Italic(true).Foreground(lipgloss.Color("#a4a6a5")).Render

const (
	TabSep   = "    "
	SpaceSep = " "
	Cr       = "\n"
)
