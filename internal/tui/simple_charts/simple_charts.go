package simple_charts

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/NimbleMarkets/ntcharts/canvas/runes"
	"github.com/NimbleMarkets/ntcharts/linechart/timeserieslinechart"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eldius/rpi-system-monitor/internal/adapter"
	"github.com/eldius/rpi-system-monitor/internal/tui/helper"
	zone "github.com/lrstanley/bubblezone"
)

var borderStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("63")).
	Padding(0, 1)

var graphLineStyle1 = lipgloss.NewStyle().
	Foreground(lipgloss.Color("4")) // blue

var axisStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("3")) // yellow

var labelStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("6")) // cyan

type tickMsg time.Time

type simpleChartsModel struct {
	cpuChart timeserieslinechart.Model
	zM       *zone.Manager
	mf       adapter.MeasureFunc
	ctx      context.Context

	lastTimestamp time.Time
}

func (m simpleChartsModel) Init() tea.Cmd {
	m.cpuChart.DrawXYAxisAndLabel()
	return tea.Batch(
		tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
			return tickMsg(t)
		}),
	)
}

func (m simpleChartsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}

	case tickMsg:
		mf, err := m.mf(m.ctx)
		if err != nil {
			fmt.Println(err)
			return m, tea.Quit
		}
		slog.With("cpu_usage", mf.CPU.CPUUsage).Debug("pushing cpu usage data")
		m.lastTimestamp = mf.Timestamp
		m.cpuChart.Push(timeserieslinechart.TimePoint{
			Time:  mf.Timestamp,
			Value: mf.Memory.MemoryUsagePercentage,
		})

		m.cpuChart.Draw()

		return m, tickCmd()
	}
	return m, nil
}

func (m simpleChartsModel) View() string {
	s := "any key to push randomized data value,`r` to clear data, `q/ctrl+c` to quit\n"
	s += "pgup/pdown/mouse wheel scroll to zoom in and out along X axis\n"
	s += "mouse click+drag or arrow keys to move view along X axis while zoomed in\n"
	s += "Latest update: " + m.lastTimestamp.Format(time.DateTime) + "\n"
	s += borderStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, labelStyle.Render("CPU Usage (%)"), m.cpuChart.View()),
	) + "\n"
	return m.zM.Scan(s) // call zone Manager.Scan() at root model
}

func tickCmd() tea.Cmd {
	return tea.Tick(30*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
func Start(ctx context.Context) error {
	minYValue := 0.0
	maxYValue := 100.0

	windowSize, _ := helper.GetTerminalSize()

	zoneManager := zone.New()

	// timeserieslinechart 1 created with New() and setting options afterwards
	cpuChart := timeserieslinechart.New(windowSize.Width, windowSize.Height)
	cpuChart.AxisStyle = axisStyle
	cpuChart.LabelStyle = labelStyle
	cpuChart.XLabelFormatter = timeserieslinechart.HourTimeLabelFormatter()
	cpuChart.UpdateHandler = timeserieslinechart.SecondUpdateHandler(1)
	cpuChart.SetYRange(minYValue, maxYValue)     // set expected Y values (values can be less or greater than what is displayed)
	cpuChart.SetViewYRange(minYValue, maxYValue) // setting display Y values will fail unless set expected Y values first
	cpuChart.SetStyle(graphLineStyle1)
	cpuChart.SetLineStyle(runes.ThinLineStyle) // ThinLineStyle replaces default linechart arcline rune style
	cpuChart.SetZoneManager(zoneManager)
	cpuChart.Focus()

	m := simpleChartsModel{
		cpuChart: cpuChart,
		zM:       zoneManager,
		mf:       adapter.Measure,
		ctx:      ctx,
	}
	if _, err := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion()).Run(); err != nil {
		err = fmt.Errorf("executing screen: %w", err)
		fmt.Println("Error running program:", err)
		return err
	}
	return nil
}
