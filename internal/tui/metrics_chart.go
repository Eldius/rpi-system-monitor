package tui

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/NimbleMarkets/ntcharts/canvas/runes"
	"github.com/NimbleMarkets/ntcharts/linechart/timeserieslinechart"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eldius/rpi-system-monitor/internal/adapter"
	"github.com/eldius/rpi-system-monitor/internal/tui/helper"
	zone "github.com/lrstanley/bubblezone"
)

var (
	_ tea.Model = &hostMetricsDisplayModel{}
)

var (
	// Estilos Lipgloss para as caixas

	borderStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(0, 1)

	headerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Padding(0, 1)

	graphLineStyle1 = lipgloss.NewStyle().
		Foreground(lipgloss.Color("4")) // blue

	axisStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("3")) // yellow

	labelStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("6")) // cyan
)

// --- Tipos de Mensagem ---

// tickMsg é enviado a cada intervalo para atualizar os dados
type tickMsg time.Time

// --- Modelo da Aplicação ---

type hostMetricsDisplayModel struct {
	hostname string
	ip       string

	// Gráficos do Ntcharts
	cpuChart  *timeserieslinechart.Model
	memChart  *timeserieslinechart.Model
	tempChart *timeserieslinechart.Model

	zm *zone.Manager

	// Dados brutos (simulados para o exemplo)
	tickCount int

	ctx context.Context
}

// initialModel configura o estado inicial
func initialModel(ctx context.Context) *hostMetricsDisplayModel {
	// Obter Hostname e IP reais
	host, _ := os.Hostname()
	ip := getOutboundIP()

	windowSize, _ := helper.GetTerminalSize()

	w := windowSize.Width - 10
	h := (windowSize.Height - 5) / 4

	cpuChart := timeserieslinechart.New(w, h)
	memChart := timeserieslinechart.New(w, h)
	tempChart := timeserieslinechart.New(w, h)

	zm := zone.New()

	setupMetricChart(&cpuChart, zm)
	setupMetricChart(&memChart, zm)
	setupMetricChart(&tempChart, zm)

	m := hostMetricsDisplayModel{
		hostname:  host,
		ip:        ip,
		cpuChart:  &cpuChart,
		memChart:  &memChart,
		tempChart: &tempChart,
		ctx:       ctx,
		zm:        zm,
	}

	return &m
}

func setupMetricChart(c *timeserieslinechart.Model, zm *zone.Manager) {
	c.AxisStyle = axisStyle
	c.LabelStyle = labelStyle
	c.XLabelFormatter = timeserieslinechart.HourTimeLabelFormatter()
	c.UpdateHandler = timeserieslinechart.SecondUpdateHandler(1)
	c.SetYRange(0, 100)     // set expected Y values (values can be less or greater than what is displayed)
	c.SetViewYRange(0, 100) // setting display Y values will fail unless set expected Y values first
	c.SetStyle(graphLineStyle1)
	c.SetLineStyle(runes.ThinLineStyle) // ThinLineStyle replaces default linechart arcline rune style
	c.SetZoneManager(zm)

	c.DrawXYAxisAndLabel()
}

// --- Métodos do Bubble Tea ---

func (m *hostMetricsDisplayModel) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(), // Inicia o loop de atualização
	)
}

func (m *hostMetricsDisplayModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", tea.KeyCtrlC.String(), tea.KeyEscape.String():
			return m, tea.Quit
		}

	case tickMsg:
		m.tickCount++

		// Simulação de dados (Substitua por leitura real de sensores aqui)
		// Usamos funções seno/cos com ruído para parecer dados reais
		measures, err := adapter.Measure(m.ctx)
		if err != nil {
			fmt.Println(err)
			return m, tea.Quit
		}
		ts := measures.Timestamp
		cpuVal := timeserieslinechart.TimePoint{Time: ts, Value: measures.CPU.CPUUsage}
		memVal := timeserieslinechart.TimePoint{Time: ts, Value: measures.Memory.MemoryUsagePercentage}
		tempVal := timeserieslinechart.TimePoint{Time: ts, Value: measures.Temp.Temperature}

		m.cpuChart.Push(cpuVal)
		m.memChart.Push(memVal)
		m.tempChart.Push(tempVal)

		m.cpuChart.Draw()
		m.memChart.Draw()
		m.tempChart.Draw()

		return m, tickCmd()
	}

	return m, nil
}

func (m *hostMetricsDisplayModel) View() string {
	// Renderiza os gráficos para strings
	cpuView := m.cpuChart.View()
	memView := m.memChart.View()
	tempView := m.tempChart.View()

	// Constrói o cabeçalho
	header := headerStyle.Render(
		lipgloss.JoinVertical(lipgloss.Top, fmt.Sprintf("🖥️  HOST: %s  |  🌐 IP: %s", m.hostname, m.ip)),
	) + "\n"

	// Cria caixas com títulos para cada métrica
	cpuBox := borderStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, labelStyle.Render("CPU Usage (%)"), cpuView),
	)
	memBox := borderStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, labelStyle.Render("Memory Usage (%)"), memView),
	)
	tempBox := borderStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, labelStyle.Render("Temperature (°C)"), tempView),
	)

	// Layout final: Cabeçalho em cima, gráficos lado a lado (se couber) ou vertical
	// Aqui usaremos vertical para garantir visualização simples
	body := lipgloss.JoinVertical(lipgloss.Left,
		header,
		cpuBox,
		memBox,
		tempBox,
		labelStyle.Render("\nPressione 'q' para sair."),
	)

	m.zm.Scan(body)

	return body
}

// --- Utilitários ---

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*5, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// getOutboundIP pega o IP preferido para saída dessa máquina
func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "Desconhecido"
	}
	defer func() { _ = conn.Close() }()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func MetricsChart(ctx context.Context) {
	p := tea.NewProgram(initialModel(ctx), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Erro ao executar TUI: %v", err)
		os.Exit(1)
	}
}
