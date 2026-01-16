package cli

import (
	"fmt"
	"strings"
	"time"
)

// ANSI color codes
const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Italic    = "\033[3m"
	Underline = "\033[4m"

	// Colors
	Black   = "\033[30m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
	Gray    = "\033[90m"

	// Bright colors
	BrightGreen  = "\033[92m"
	BrightYellow = "\033[93m"
	BrightCyan   = "\033[96m"
	BrightWhite  = "\033[97m"

	// Background colors
	BgGreen = "\033[42m"
	BgRed   = "\033[41m"
)

// Styled output functions
func Success(msg string) string {
	return fmt.Sprintf("%s%s%s%s", Bold, Green, msg, Reset)
}

func Error(msg string) string {
	return fmt.Sprintf("%s%s%s%s", Bold, Red, msg, Reset)
}

func Warning(msg string) string {
	return fmt.Sprintf("%s%s%s%s", Bold, Yellow, msg, Reset)
}

func Info(msg string) string {
	return fmt.Sprintf("%s%s%s", Cyan, msg, Reset)
}

func Muted(msg string) string {
	return fmt.Sprintf("%s%s%s", Gray, msg, Reset)
}

func Highlight(msg string) string {
	return fmt.Sprintf("%s%s%s%s", Bold, BrightGreen, msg, Reset)
}

func Label(label, value string) string {
	return fmt.Sprintf("%s%s:%s %s", Gray, label, Reset, value)
}

// Icons
const (
	IconCheck    = "✓"
	IconCross    = "✗"
	IconArrow    = "→"
	IconDot      = "●"
	IconCircle   = "○"
	IconStar     = "★"
	IconSparkle  = "✦"
	IconBolt     = "⚡"
	IconBrain    = "🧠"
	IconRobot    = "🤖"
	IconLink     = "🔗"
	IconLock     = "🔒"
	IconKey      = "🔑"
	IconGear     = "⚙"
	IconChart    = "📊"
	IconTarget   = "🎯"
	IconRocket   = "🚀"
)

// Box drawing characters
const (
	BoxTopLeft     = "┌"
	BoxTopRight    = "┐"
	BoxBottomLeft  = "└"
	BoxBottomRight = "┘"
	BoxHorizontal  = "─"
	BoxVertical    = "│"
	BoxTLeft       = "├"
	BoxTRight      = "┤"
)

// Banner prints the Squaremind ASCII banner
func Banner() string {
	banner := `
    ███████╗ ██████╗ ██╗   ██╗ █████╗ ██████╗ ███████╗███╗   ███╗██╗███╗   ██╗██████╗
    ██╔════╝██╔═══██╗██║   ██║██╔══██╗██╔══██╗██╔════╝████╗ ████║██║████╗  ██║██╔══██╗
    ███████╗██║   ██║██║   ██║███████║██████╔╝█████╗  ██╔████╔██║██║██╔██╗ ██║██║  ██║
    ╚════██║██║▄▄ ██║██║   ██║██╔══██║██╔══██╗██╔══╝  ██║╚██╔╝██║██║██║╚██╗██║██║  ██║
    ███████║╚██████╔╝╚██████╔╝██║  ██║██║  ██║███████╗██║ ╚═╝ ██║██║██║ ╚████║██████╔╝
    ╚══════╝ ╚══▀▀═╝  ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝╚═╝     ╚═╝╚═╝╚═╝  ╚═══╝╚═════╝ `

	return fmt.Sprintf("%s%s%s", Green, banner, Reset)
}

// SmallBanner prints a smaller banner
func SmallBanner() string {
	return fmt.Sprintf(`
  %s%s┌─────────────────────────────────────┐%s
  %s│%s  %s◆ SQUAREMIND%s  %s│%s
  %s│%s     %sMany Agents. One Mind.%s          %s│%s
  %s└─────────────────────────────────────┘%s
`, Green, Bold, Reset,
		Green, Reset, Bold+BrightGreen, Reset, Green, Reset,
		Green, Reset, Dim, Reset, Green, Reset,
		Green, Reset)
}

// Divider creates a horizontal divider
func Divider(width int) string {
	return fmt.Sprintf("%s%s%s", Gray, strings.Repeat(BoxHorizontal, width), Reset)
}

// Section prints a section header
func Section(title string) string {
	return fmt.Sprintf("\n%s%s %s %s%s\n", Bold, Green, IconSparkle, title, Reset)
}

// StatusLine prints a status with icon
func StatusLine(status, message string) string {
	var icon, color string
	switch status {
	case "success":
		icon, color = IconCheck, Green
	case "error":
		icon, color = IconCross, Red
	case "warning":
		icon, color = "!", Yellow
	case "info":
		icon, color = IconArrow, Cyan
	case "pending":
		icon, color = IconCircle, Gray
	default:
		icon, color = IconDot, White
	}
	return fmt.Sprintf("  %s%s%s %s", color, icon, Reset, message)
}

// ProgressBar creates a simple progress bar
func ProgressBar(current, total int, width int) string {
	if total == 0 {
		total = 1
	}
	percent := float64(current) / float64(total)
	filled := int(percent * float64(width))
	empty := width - filled

	bar := fmt.Sprintf("%s%s%s%s",
		Green, strings.Repeat("█", filled),
		Gray, strings.Repeat("░", empty))

	return fmt.Sprintf("[%s%s] %s%.0f%%%s", bar, Reset, Dim, percent*100, Reset)
}

// Spinner animation frames
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// Spinner provides an animated spinner
type Spinner struct {
	message string
	done    chan bool
	running bool
}

// NewSpinner creates a new spinner
func NewSpinner(message string) *Spinner {
	return &Spinner{
		message: message,
		done:    make(chan bool),
	}
}

// Start starts the spinner
func (s *Spinner) Start() {
	if s.running {
		return
	}
	s.running = true

	go func() {
		i := 0
		for {
			select {
			case <-s.done:
				return
			default:
				frame := spinnerFrames[i%len(spinnerFrames)]
				fmt.Printf("\r  %s%s%s %s", Green, frame, Reset, s.message)
				i++
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()
}

// Stop stops the spinner
func (s *Spinner) Stop(success bool) {
	if !s.running {
		return
	}
	s.running = false
	s.done <- true

	icon := IconCheck
	color := Green
	if !success {
		icon = IconCross
		color = Red
	}

	fmt.Printf("\r  %s%s%s %s\n", color, icon, Reset, s.message)
}

// StopWithMessage stops with a custom message
func (s *Spinner) StopWithMessage(success bool, message string) {
	if !s.running {
		return
	}
	s.running = false
	s.done <- true

	icon := IconCheck
	color := Green
	if !success {
		icon = IconCross
		color = Red
	}

	fmt.Printf("\r  %s%s%s %s\n", color, icon, Reset, message)
}

// Table helpers
func TableHeader(headers ...string) string {
	var parts []string
	for _, h := range headers {
		parts = append(parts, fmt.Sprintf("%s%s%s", Bold, h, Reset))
	}
	return "  " + strings.Join(parts, "  │  ")
}

// AgentCard prints a formatted agent card
func AgentCard(name, sid, state string, reputation float64, capabilities []string) string {
	stateColor := Gray
	switch state {
	case "idle":
		stateColor = Green
	case "working":
		stateColor = Yellow
	case "paused":
		stateColor = Gray
	}

	capsStr := strings.Join(capabilities, ", ")
	if len(capsStr) > 40 {
		capsStr = capsStr[:37] + "..."
	}

	return fmt.Sprintf(`  %s┌─────────────────────────────────────────────┐%s
  %s│%s %s%-20s%s %s%s%s %s│%s
  %s│%s   SID: %s%-36s%s %s│%s
  %s│%s   Rep: %s%-5.1f%s  State: %s%s%-10s%s %s│%s
  %s│%s   Caps: %s%-35s%s %s│%s
  %s└─────────────────────────────────────────────┘%s`,
		Gray, Reset,
		Gray, Reset, Bold+BrightGreen, name, Reset, Dim, IconRobot, Reset, Gray, Reset,
		Gray, Reset, Cyan, sid[:36], Reset, Gray, Reset,
		Gray, Reset, Yellow, reputation, Reset, stateColor, Bold, state, Reset, Gray, Reset,
		Gray, Reset, Dim, capsStr, Reset, Gray, Reset,
		Gray, Reset,
	)
}

// TaskCard prints a formatted task card
func TaskCard(id, description, status string, complexity string, assignedTo string) string {
	statusIcon := IconCircle
	statusColor := Gray
	switch status {
	case "completed":
		statusIcon = IconCheck
		statusColor = Green
	case "running", "assigned":
		statusIcon = IconBolt
		statusColor = Yellow
	case "failed":
		statusIcon = IconCross
		statusColor = Red
	case "pending":
		statusIcon = IconCircle
		statusColor = Cyan
	}

	desc := description
	if len(desc) > 45 {
		desc = desc[:42] + "..."
	}

	return fmt.Sprintf(`  %s%s%s %s%s%s
      ID: %s%s%s
      Complexity: %s  Assigned: %s`,
		statusColor, statusIcon, Reset, Bold, desc, Reset,
		Dim, id[:8], Reset,
		complexity, assignedTo,
	)
}

// PrintKeyValue prints a key-value pair with styling
func PrintKeyValue(key, value string) {
	fmt.Printf("  %s%s:%s %s\n", Gray, key, Reset, value)
}

// PrintStats prints statistics in a nice format
func PrintStats(stats map[string]interface{}) {
	fmt.Println()
	for k, v := range stats {
		fmt.Printf("  %s%-20s%s %s%v%s\n", Dim, k, Reset, Bold, v, Reset)
	}
	fmt.Println()
}
