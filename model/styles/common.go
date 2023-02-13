package styles

import "github.com/charmbracelet/lipgloss"

const (
	White     = lipgloss.Color("#FFFFFF")
	RetroBlue = lipgloss.Color("63")
	Yellow    = lipgloss.Color("228")
	ErrorRed  = lipgloss.Color("#ab2727")
	BasicGrey = lipgloss.Color("#a4a6a5")
)

var LabelStyle = lipgloss.NewStyle().Bold(true).Foreground(White).Underline(true)
var LabelStyleRender = LabelStyle.Render
var SubLabelStyle = LabelStyle.Copy().Bold(false).Italic(true).Underline(false).Render
var ErrorStyle = lipgloss.NewStyle().Foreground(ErrorRed).Render
var TitleStyle = lipgloss.NewStyle().Margin(1).Padding(0, 2).Align(lipgloss.Center).Bold(true).Foreground(Yellow).BorderStyle(lipgloss.RoundedBorder()).BorderForeground(RetroBlue).Render
var TabStyle = lipgloss.NewStyle().Padding(0, 1).Margin(0, 0, 1).Align(lipgloss.Center).BorderStyle(lipgloss.RoundedBorder()).BorderForeground(RetroBlue).BorderBottom(false).BorderTop(true).BorderLeft(true).BorderRight(true)
var ActivatedTabStyle = TabStyle.Copy().Bold(true).Foreground(White).Render
var DeactivatedTabStyle = TabStyle.Copy().Italic(true).Foreground(BasicGrey).Render
var ModalStyle = lipgloss.NewStyle().Padding(2, 1).BorderStyle(lipgloss.RoundedBorder()).BorderForeground(RetroBlue)

const (
	TabSep   = "    "
	SpaceSep = " "
	Cr       = "\n"
)
