package simple_charts

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/NimbleMarkets/ntcharts/canvas/runes"
	"github.com/NimbleMarkets/ntcharts/linechart/timeserieslinechart"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

var defaultStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("63")) // purple

var graphLineStyle1 = lipgloss.NewStyle().
	Foreground(lipgloss.Color("4")) // blue

var graphLineStyle2 = lipgloss.NewStyle().
	Foreground(lipgloss.Color("10")) // green

var axisStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("3")) // yellow

var labelStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("6")) // cyan

var timePoint1 timeserieslinechart.TimePoint
var randf1 float64

type simpleChartsModel struct {
	c1 timeserieslinechart.Model
	zM *zone.Manager
}

func (m simpleChartsModel) Init() tea.Cmd {
	m.c1.DrawXYAxisAndLabel()
	return nil
}

func (m simpleChartsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	addPoint := false
	forwardMsg := false
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			m.c1.ClearAllData()
			m.c1.Clear()
			m.c1.DrawXYAxisAndLabel()

			return m, nil
		case "up", "down", "left", "right", "pgup", "pgdown":
			forwardMsg = true
		case "q", "ctrl+c":
			return m, tea.Quit
		default:
			addPoint = true
		}
	case tea.MouseMsg:
		if msg.Action == tea.MouseActionPress {
			m.c1.Blur()

			// switch to whichever canvas was clicked on
			switch {
			case m.zM.Get(m.c1.ZoneID()).InBounds(msg):
				m.c1.Focus()
			}
		}
		forwardMsg = true
	}
	if addPoint {
		// generate random numbers within the given Y value range
		//rangeNumbers := m.tslc2.MaxY() - m.tslc2.MinY()
		rangeNumbers := float64(100)
		randf1 = rand.Float64()*rangeNumbers + m.c1.MinY()

		now := time.Now()
		timePoint1 = timeserieslinechart.TimePoint{Time: now, Value: randf1}

		// timeserieslinechart 1 and 3 pushes random value randf1 to default data set
		m.c1.Push(timePoint1)

		m.c1.Draw()
	}
	// timeserieslinechart handles mouse events
	if forwardMsg {
		switch {
		case m.c1.Focused():
			m.c1, _ = m.c1.Update(msg)
			m.c1.Draw()
		}
	}
	return m, nil
}

func (m simpleChartsModel) View() string {
	s := "any key to push randomized data value,`r` to clear data, `q/ctrl+c` to quit\n"
	s += "pgup/pdown/mouse wheel scroll to zoom in and out along X axis\n"
	s += "mouse click+drag or arrow keys to move view along X axis while zoomed in\n"
	s += lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.JoinVertical(lipgloss.Left,
			defaultStyle.Render(fmt.Sprintf("ts:%s, f1:(%.02f)\n", timePoint1.Time.UTC().Format("15:04:05"), randf1)+m.c1.View()),
		),
	) + "\n"
	return m.zM.Scan(s) // call zone Manager.Scan() at root model
}

func Start(ctx context.Context) error {
	width := 36
	height := 8
	minYValue := 0.0
	maxYValue := 100.0

	// timeserieslinecharts creates line charts starting with time as time.Now()
	// There are two sets of charts, one show regular lines and one showing braille lines
	// Pressing keys will insert random Y value data into chart with time.Now() (when key was pressed)

	// create new bubblezone Manager to enable mouse support to zoom in and out of chart
	zoneManager := zone.New()

	// timeserieslinechart 1 created with New() and setting options afterwards
	tslc1 := timeserieslinechart.New(width, height)
	tslc1.AxisStyle = axisStyle
	tslc1.LabelStyle = labelStyle
	tslc1.XLabelFormatter = timeserieslinechart.HourTimeLabelFormatter()
	tslc1.UpdateHandler = timeserieslinechart.SecondUpdateHandler(1)
	tslc1.SetYRange(minYValue, maxYValue)     // set expected Y values (values can be less or greater than what is displayed)
	tslc1.SetViewYRange(minYValue, maxYValue) // setting display Y values will fail unless set expected Y values first
	tslc1.SetStyle(graphLineStyle1)
	tslc1.SetLineStyle(runes.ThinLineStyle) // ThinLineStyle replaces default linechart arcline rune style
	tslc1.SetZoneManager(zoneManager)
	tslc1.Focus()

	m := simpleChartsModel{c1: tslc1,
		zM: zoneManager,
	}
	if _, err := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion()).Run(); err != nil {
		err = fmt.Errorf("executing screen: %w", err)
		fmt.Println("Error running program:", err)
		return err
	}
	return nil
}
