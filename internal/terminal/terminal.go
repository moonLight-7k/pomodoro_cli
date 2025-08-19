package terminal

import (
	"fmt"
	"os"
	"strings"
	"time"

	"pomodoro_cli/internal/errors"
)

type TerminalCapabilities struct {
	SupportsColor  bool
	SupportsANSI   bool
	SupportsClear  bool
	TerminalWidth  int
	TerminalHeight int
}

type Terminal struct {
	capabilities *TerminalCapabilities
	logger       *errors.Logger
}

func NewTerminal(logger *errors.Logger) *Terminal {
	caps := detectTerminalCapabilities()
	return &Terminal{
		capabilities: caps,
		logger:       logger,
	}
}

func detectTerminalCapabilities() *TerminalCapabilities {
	caps := &TerminalCapabilities{
		SupportsColor:  false,
		SupportsANSI:   false,
		SupportsClear:  false,
		TerminalWidth:  80, // Safe default
		TerminalHeight: 24, // Safe default
	}

	term := os.Getenv("TERM")
	colorTerm := os.Getenv("COLORTERM")

	if strings.Contains(term, "color") ||
		strings.Contains(term, "xterm") ||
		strings.Contains(term, "screen") ||
		colorTerm != "" {
		caps.SupportsColor = true
		caps.SupportsANSI = true
	}

	if term != "" && term != "dumb" {
		caps.SupportsClear = true
	}

	return caps
}

func (t *Terminal) ClearScreen() error {
	if !t.capabilities.SupportsClear {
		for i := 0; i < 10; i++ {
			fmt.Println()
		}
		return nil
	}

	fmt.Print("\033[H\033[2J")
	return nil
}

type Colors struct {
	Purple   string
	White    string
	DarkGray string
	Reset    string
}

func (t *Terminal) GetColors() Colors {
	if !t.capabilities.SupportsColor {
		return Colors{
			Purple:   "",
			White:    "",
			DarkGray: "",
			Reset:    "",
		}
	}

	return Colors{
		Purple:   "\033[38;2;138;43;226m",
		White:    "\033[37m",
		DarkGray: "\033[48;2;64;64;64m",
		Reset:    "\033[0m",
	}
}

func formatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%dm%02ds", minutes, seconds)
}

func (t *Terminal) DrawProgressBar(progress float64, width int) string {
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}

	if width > t.capabilities.TerminalWidth-10 {
		width = t.capabilities.TerminalWidth - 10
	}
	if width < 10 {
		width = 10
	}

	filledWidth := int(float64(width) * progress)
	emptyWidth := width - filledWidth

	if !t.capabilities.SupportsColor {
		filled := strings.Repeat("█", filledWidth)
		empty := strings.Repeat("░", emptyWidth)
		return fmt.Sprintf("[%s%s]", filled, empty)
	}

	colors := t.GetColors()
	filled := colors.Purple + "\033[48;2;138;43;226m" + strings.Repeat(" ", filledWidth) + colors.Reset
	empty := colors.DarkGray + strings.Repeat(" ", emptyWidth) + colors.Reset

	return fmt.Sprintf("%s%s", filled, empty)
}

type SessionInfo struct {
	Label            string
	Elapsed          time.Duration
	Progress         float64
	ProgressBarWidth int
}

func (t *Terminal) DisplaySession(info SessionInfo) error {
	if err := t.ClearScreen(); err != nil {
		return fmt.Errorf("failed to clear screen: %w", err)
	}

	colors := t.GetColors()

	fmt.Printf("%s%s%s\n", colors.Purple, strings.ToLower(info.Label), colors.Reset)

	now := time.Now()
	elapsedFormatted := formatDuration(info.Elapsed)
	fmt.Printf("\033[1m%s%s\033[0m - \033[1m%s%s\033[0m\n",
		colors.White, now.Format("3:04 PM"), elapsedFormatted, colors.White)

	progressBar := t.DrawProgressBar(info.Progress, info.ProgressBarWidth)
	fmt.Printf("%s", progressBar)

	fmt.Printf("  %s%d%%%s\n", colors.White, int(info.Progress*100), colors.Reset)

	if t.capabilities.SupportsANSI {
		fmt.Printf("\n%sPress Ctrl+C to exit%s\n", colors.DarkGray, colors.Reset)
	} else {
		fmt.Printf("\nPress Ctrl+C to exit\n")
	}

	return nil
}

func (t *Terminal) DisplayCompletion(label string) error {
	if err := t.ClearScreen(); err != nil {
		return err
	}

	colors := t.GetColors()
	fmt.Printf("%s%s complete!%s\n", colors.Purple, label, colors.Reset)

	if t.capabilities.SupportsANSI {
		fmt.Printf("%s✓ Session finished successfully%s\n", colors.White, colors.Reset)
	} else {
		fmt.Printf("* Session finished successfully\n")
	}

	return nil
}
