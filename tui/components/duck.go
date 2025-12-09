package components

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hassek/bc-cli/tui/styles"
)

type DuckState int

const (
	DuckStateIdle DuckState = iota
	DuckStateAction
)

type tickMsg time.Time

type DuckComponent struct {
	state       DuckState
	frame       int
	actionFrame int
	tickRate    time.Duration
}

func NewDuckComponent() *DuckComponent {
	return &DuckComponent{
		state:    DuckStateIdle,
		frame:    0,
		tickRate: 500 * time.Millisecond,
	}
}

func (d *DuckComponent) Init() tea.Cmd {
	return d.tick()
}

func (d *DuckComponent) tick() tea.Cmd {
	return tea.Tick(d.tickRate, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (d *DuckComponent) Update(msg tea.Msg) (*DuckComponent, tea.Cmd) {
	switch msg.(type) {
	case tickMsg:
		switch d.state {
		case DuckStateIdle:
			d.frame = (d.frame + 1) % len(idleFrames)
		case DuckStateAction:
			d.actionFrame++
			if d.actionFrame >= len(actionFrames) {
				// Return to idle state after action completes
				d.state = DuckStateIdle
				d.actionFrame = 0
				d.frame = 0
			}
		}
		return d, d.tick()
	}
	return d, nil
}

func (d *DuckComponent) TriggerAction() {
	d.state = DuckStateAction
	d.actionFrame = 0
}

func (d *DuckComponent) View() string {
	var art string
	if d.state == DuckStateIdle {
		art = idleFrames[d.frame]
	} else {
		art = actionFrames[d.actionFrame]
	}

	// Apply styling
	lines := strings.Split(art, "\n")
	styledLines := make([]string, len(lines))
	for i, line := range lines {
		// Duck body in cyan, coffee/steam in yellow
		if strings.Contains(line, "☕") || strings.Contains(line, "~") {
			styledLines[i] = styles.DuckAccentStyle.Render(line)
		} else {
			styledLines[i] = styles.DuckStyle.Render(line)
		}
	}

	return strings.Join(styledLines, "\n") + "\n\n"
}

// Idle animation frames - subtle blink (fixed height)
var idleFrames = []string{
	// Frame 0 - eyes open
	`       ,~&.
(_   (  e )>    ☕
 ) ` + "`~~'   (      ~" + `
(   ` + "`-._)  )" + `
 ` + "`-._____,'" + ``,

	// Frame 1 - eyes open
	`       ,~&.
(_   (  e )>    ☕
 ) ` + "`~~'   (      ~" + `
(   ` + "`-._)  )" + `
 ` + "`-._____,'" + ``,

	// Frame 2 - eyes open
	`       ,~&.
(_   (  e )>    ☕
 ) ` + "`~~'   (      ~" + `
(   ` + "`-._)  )" + `
 ` + "`-._____,'" + ``,

	// Frame 3 - subtle blink
	`       ,~&.
(_   (  - )>    ☕
 ) ` + "`~~'   (      ~" + `
(   ` + "`-._)  )" + `
 ` + "`-._____,'" + ``,
}

// Action animation frames - celebration/coffee sip (fixed height)
var actionFrames = []string{
	// Frame 0 - prepare
	`       ,~&.
(_   (  e )>    ☕
 ) ` + "`~~'   (      ~" + `
(   ` + "`-._)  )" + `
 ` + "`-._____,'" + ``,

	// Frame 1 - excited!
	`       ,~&.
(_   (  O )>  ☕
 ) ` + "`~~'   (    ~~" + `
(   ` + "`-._)  )" + `
 ` + "`-._____,'" + ``,

	// Frame 2 - very happy
	`       ,~&.
(_   (  ^ )> ☕
 ) ` + "`~~'   (   ~~~" + `
(   ` + "`-._)  )" + `
 ` + "`-._____,'" + ``,

	// Frame 3 - drinking coffee
	`       ,~&.☕
(_   (  ◠ )>
 ) ` + "`~~'   (      ~" + `
(   ` + "`-._)  )" + `
 ` + "`-._____,'" + ``,

	// Frame 4 - satisfied
	`       ,~&.
(_   (  ◡ )>    ☕
 ) ` + "`~~'   (      ~" + `
(   ` + "`-._)  )" + `
 ` + "`-._____,'" + ``,

	// Frame 5 - return to normal
	`       ,~&.
(_   (  e )>    ☕
 ) ` + "`~~'   (      ~" + `
(   ` + "`-._)  )" + `
 ` + "`-._____,'" + ``,
}
