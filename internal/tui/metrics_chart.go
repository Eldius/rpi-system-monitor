package tui

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/NimbleMarkets/ntcharts/sparkline"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eldius/rpi-system-monitor/internal/adapter"
)

// --- Constantes e Estilos ---

const (
	width  = 60
	height = 10
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

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			PaddingLeft(1)

	chartStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("63"))
	/*
		blockStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("63")) // purple

		blockStyle2 = lipgloss.NewStyle().
				Foreground(lipgloss.Color("9")). // red
				Background(lipgloss.Color("2"))  // green

		blockStyle3 = lipgloss.NewStyle().
				Foreground(lipgloss.Color("6")). // cyan
				Background(lipgloss.Color("3"))  // yellow

		blockStyle4 = lipgloss.NewStyle().
				Foreground(lipgloss.Color("3")) // yellow
	*/
)

// --- Tipos de Mensagem ---

// tickMsg é enviado a cada intervalo para atualizar os dados
type tickMsg time.Time

// --- Modelo da Aplicação ---

type hostMetricsDisplayModel struct {
	hostname string
	ip       string

	// Gráficos do Ntcharts
	cpuChart  *sparkline.Model
	memChart  *sparkline.Model
	tempChart *sparkline.Model

	// Dados brutos (simulados para o exemplo)
	tickCount int

	ctx context.Context
}

// initialModel configura o estado inicial
func initialModel(ctx context.Context) hostMetricsDisplayModel {
	// Obter Hostname e IP reais
	host, _ := os.Hostname()
	ip := getOutboundIP()

	// Configuração base dos gráficos
	//lcConfig := sparkline.Config{
	//	Width:  width,
	//	Height: height,
	//	MinY:   0,
	//	MaxY:   100,
	//}
	/*
		var cpuUsg, memUsg, temp []float64
		for _, v := range data {
			cpuUsg = append(cpuUsg, v.CPU.CPUUsage)
			memUsg = append(memUsg, v.Memory.MemoryUsagePercentage)
			temp = append(temp, v.Temp.Temperature)
		}

		cpuChart := sparkline.New(width, height, sparkline.WithData(cpuUsg))
		memChart := sparkline.New(width, height, sparkline.WithData(memUsg))
		tempChart := sparkline.New(width, height, sparkline.WithData(temp))
	*/

	cpuChart := sparkline.New(width, height, sparkline.WithStyle(chartStyle))
	memChart := sparkline.New(width, height, sparkline.WithStyle(chartStyle))
	tempChart := sparkline.New(width, height, sparkline.WithStyle(chartStyle))

	m := hostMetricsDisplayModel{
		hostname:  host,
		ip:        ip,
		cpuChart:  &cpuChart,
		memChart:  &memChart,
		tempChart: &tempChart,
		ctx:       ctx,
	}

	return m
}

// --- Métodos do Bubble Tea ---

func (m hostMetricsDisplayModel) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(), // Inicia o loop de atualização
	)
}

func (m hostMetricsDisplayModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
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
		cpuVal := measures.CPU.CPUUsage
		memVal := measures.Memory.MemoryUsagePercentage
		tempVal := measures.Temp.Temperature

		// Atualiza os gráficos
		// Ntcharts usa coordenadas normalizadas ou raw. Aqui empurramos valores na timeline.
		m.cpuChart.Push(cpuVal)
		m.memChart.Push(memVal)
		m.tempChart.Push(tempVal)

		_, _ = fmt.Fprint(os.Stdout, "\007")

		return m, tickCmd()
	}

	return m, nil
}

func (m hostMetricsDisplayModel) View() string {
	// Renderiza os gráficos para strings
	cpuView := m.cpuChart.View()
	memView := m.memChart.View()
	tempView := m.tempChart.View()

	// Constrói o cabeçalho
	header := headerStyle.Render(fmt.Sprintf("🖥️  HOST: %s  |  🌐 IP: %s", m.hostname, m.ip))

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
		lipgloss.JoinVertical(lipgloss.Left, cpuBox, memBox, tempBox),
		labelStyle.Render("\nPressione 'q' para sair."),
	)

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
	defer conn.Close()
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
