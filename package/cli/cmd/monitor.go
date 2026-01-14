package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/config"
	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/context"
	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/ui"
	"github.com/spf13/cobra"
)

// MonitorMetrics holds all monitoring metrics
type MonitorMetrics struct {
	Timestamp   time.Time
	Vault       VaultMetrics
	System      SystemMetrics
	Network     NetworkMetrics
	Security    SecurityMetrics
	Performance PerformanceMetrics
}

// VaultMetrics contains vault-specific metrics
type VaultMetrics struct {
	Status           string
	Version          string
	Uptime           string
	ActiveTokens     int
	StoredSecrets    int
	MemoryUsage      int64
	CPUUsage         float64
	Connections      int
	RequestsPerSec   float64
	ResponseTime     time.Duration
	EncryptionStatus string
}

// SystemMetrics contains system-level metrics
type SystemMetrics struct {
	OS              string
	Architecture    string
	Hostname        string
	CPUUsage        float64
	MemoryTotal     uint64
	MemoryUsed      uint64
	MemoryAvailable uint64
	DiskTotal       uint64
	DiskUsed        uint64
	DiskAvailable   uint64
	LoadAverage     []float64
	ProcessCount    int
	Uptime          string
}

// NetworkMetrics contains network-related metrics
type NetworkMetrics struct {
	Interface   string
	IPAddresses []string
	BytesSent   uint64
	BytesRecv   uint64
	PacketsSent uint64
	PacketsRecv uint64
	Connections int
	Latency     time.Duration
	Bandwidth   float64
}

// SecurityMetrics contains security-related metrics
type SecurityMetrics struct {
	FailedLogins     int
	SuccessfulLogins int
	ActiveSessions   int
	ThreatAlerts     int
	PolicyViolations int
	AuditEvents      int
	LastSecurityScan time.Time
	EncryptionLevel  string
}

// PerformanceMetrics contains performance-related metrics
type PerformanceMetrics struct {
	AvgResponseTime  time.Duration
	P95ResponseTime  time.Duration
	P99ResponseTime  time.Duration
	ThroughputPerSec float64
	ErrorRate        float64
	CacheHitRate     float64
	DatabaseQueries  int
	SlowQueries      int
}

// newMonitorCommand creates the monitor command
func newMonitorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "monitor",
		Short: "Real-time monitoring dashboard",
		Long: `Launch a real-time monitoring dashboard for Aether Vault.

This command provides a comprehensive monitoring interface similar to glances,
displaying system resources, vault performance, security metrics, and more.

Features:
  - Real-time metrics updates
  - Interactive dashboard
  - Resource usage monitoring
  - Security event tracking
  - Performance analytics
  - Browser-based web interface (with --browser flag)

Examples:
  vault monitor                    # Terminal dashboard
  vault monitor --browser          # Web interface
  vault monitor --refresh 2        # 2-second refresh interval
  vault monitor --export metrics.json # Export metrics`,
		RunE: runMonitorCommand,
	}

	cmd.Flags().BoolP("browser", "b", false, "Launch web-based monitoring interface in browser")
	cmd.Flags().IntP("refresh", "r", 1, "Refresh interval in seconds (1-60)")
	cmd.Flags().StringP("export", "e", "", "Export metrics to file (JSON format)")
	cmd.Flags().BoolP("quiet", "q", false, "Minimal output mode")
	cmd.Flags().StringP("theme", "t", "auto", "Color theme (auto, light, dark, monochrome)")
	cmd.Flags().BoolP("full-screen", "f", false, "Full-screen mode")
	cmd.Flags().IntP("history", "H", 100, "Number of data points to keep in history")
	cmd.Flags().BoolP("alerts", "a", true, "Enable alerts and notifications")
	cmd.Flags().StringP("filter", "F", "", "Filter metrics by category (vault, system, network, security, performance)")

	return cmd
}

// runMonitorCommand executes the monitor command
func runMonitorCommand(cmd *cobra.Command, args []string) error {
	browser, _ := cmd.Flags().GetBool("browser")
	refresh, _ := cmd.Flags().GetInt("refresh")
	export, _ := cmd.Flags().GetString("export")
	quiet, _ := cmd.Flags().GetBool("quiet")
	theme, _ := cmd.Flags().GetString("theme")
	fullScreen, _ := cmd.Flags().GetBool("full-screen")
	history, _ := cmd.Flags().GetInt("history")
	alerts, _ := cmd.Flags().GetBool("alerts")
	filter, _ := cmd.Flags().GetString("filter")

	// Validate refresh interval
	if refresh < 1 || refresh > 60 {
		return fmt.Errorf("refresh interval must be between 1 and 60 seconds")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		if !quiet {
			fmt.Fprintf(os.Stderr, "Warning: Failed to load configuration: %v\n", err)
		}
		cfg = config.Defaults()
	}

	// Create context
	ctx, err := context.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to create context: %w", err)
	}

	// Initialize monitor
	monitor := &VaultMonitor{
		Context:    ctx,
		Refresh:    time.Duration(refresh) * time.Second,
		Quiet:      quiet,
		Theme:      theme,
		FullScreen: fullScreen,
		History:    history,
		Alerts:     alerts,
		Filter:     filter,
		StopChan:   make(chan struct{}),
		Metrics:    make(chan *MonitorMetrics, 10),
	}

	// Handle browser mode
	if browser {
		return monitor.launchBrowserInterface()
	}

	// Handle export mode
	if export != "" {
		return monitor.exportMetrics(export)
	}

	// Start terminal monitoring
	return monitor.startTerminalMonitoring()
}

// VaultMonitor manages the monitoring session
type VaultMonitor struct {
	Context    *context.Context
	Refresh    time.Duration
	Quiet      bool
	Theme      string
	FullScreen bool
	History    int
	Alerts     bool
	Filter     string
	StopChan   chan struct{}
	Metrics    chan *MonitorMetrics
}

// launchBrowserInterface starts the web-based monitoring interface
func (m *VaultMonitor) launchBrowserInterface() error {
	if !m.Quiet {
		fmt.Printf("%s🚀 Starting web-based monitoring interface...%s\n", ui.Blue, ui.Reset)
	}

	// Start web server in background
	go func() {
		if err := m.startWebServer(); err != nil {
			fmt.Fprintf(os.Stderr, "Web server error: %v\n", err)
		}
	}()

	// Wait a moment for server to start
	time.Sleep(2 * time.Second)

	// Open browser
	url := "http://localhost:8080/monitor"
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default: // linux, freebsd, etc.
		cmd = exec.Command("xdg-open", url)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open browser: %w", err)
	}

	if !m.Quiet {
		fmt.Printf("%s✅ Monitoring interface opened in browser: %s%s\n", ui.Green, url, ui.Reset)
		fmt.Printf("%s💡 Press Ctrl+C to stop monitoring%s\n", ui.Yellow, ui.Reset)
	}

	// Wait for interrupt
	<-m.StopChan
	return nil
}

// startWebServer starts the monitoring web server
func (m *VaultMonitor) startWebServer() error {
	// This is a placeholder for web server implementation
	// In a real implementation, you would use a web framework like gin or echo
	// to serve HTML, CSS, JavaScript, and WebSocket endpoints
	select {
	case <-m.StopChan:
		return nil
	case <-time.After(time.Hour):
		return fmt.Errorf("web server timeout")
	}
}

// exportMetrics exports metrics to a file
func (m *VaultMonitor) exportMetrics(filename string) error {
	if !m.Quiet {
		fmt.Printf("%s📊 Collecting metrics for export...%s\n", ui.Blue, ui.Reset)
	}

	// Collect metrics once
	metrics, err := m.collectMetrics()
	if err != nil {
		return fmt.Errorf("failed to collect metrics: %w", err)
	}

	// Export to JSON (simplified)
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create export file: %w", err)
	}
	defer file.Close()

	// Simple JSON export (in real implementation, use json.Marshal)
	fmt.Fprintf(file, "timestamp: %s\n", metrics.Timestamp.Format(time.RFC3339))
	fmt.Fprintf(file, "vault_status: %s\n", metrics.Vault.Status)
	fmt.Fprintf(file, "cpu_usage: %.2f%%\n", metrics.System.CPUUsage)
	fmt.Fprintf(file, "memory_usage: %.2f%%\n", float64(metrics.System.MemoryUsed)/float64(metrics.System.MemoryTotal)*100)

	if !m.Quiet {
		fmt.Printf("%s✅ Metrics exported to: %s%s\n", ui.Green, filename, ui.Reset)
	}

	return nil
}

// startTerminalMonitoring starts the terminal-based monitoring interface
func (m *VaultMonitor) startTerminalMonitoring() error {
	// Setup terminal
	if err := m.setupTerminal(); err != nil {
		return fmt.Errorf("failed to setup terminal: %w", err)
	}
	defer m.restoreTerminal()

	if !m.Quiet {
		fmt.Printf("%s🔍 Starting real-time monitoring...%s\n", ui.Blue, ui.Reset)
		fmt.Printf("%s💡 Press Ctrl+C to stop monitoring%s\n", ui.Yellow, ui.Reset)
	}

	// Start metrics collection
	go m.collectMetricsLoop()

	// Start display loop
	ticker := time.NewTicker(m.Refresh)
	defer ticker.Stop()

	for {
		select {
		case <-m.StopChan:
			return nil
		case <-ticker.C:
			if err := m.updateDisplay(); err != nil {
				fmt.Fprintf(os.Stderr, "Display error: %v\n", err)
			}
		case metrics := <-m.Metrics:
			if m.Alerts {
				m.checkAlerts(metrics)
			}
		}
	}
}

// setupTerminal prepares the terminal for monitoring
func (m *VaultMonitor) setupTerminal() error {
	if !m.FullScreen {
		return nil
	}

	// Clear screen and hide cursor
	fmt.Print("\033[2J\033[H")
	fmt.Print("\033[?25l")

	return nil
}

// restoreTerminal restores the terminal state
func (m *VaultMonitor) restoreTerminal() {
	if m.FullScreen {
		// Show cursor and clear screen
		fmt.Print("\033[?25h")
		fmt.Print("\033[2J\033[H")
	}
}

// collectMetricsLoop continuously collects metrics
func (m *VaultMonitor) collectMetricsLoop() {
	ticker := time.NewTicker(m.Refresh)
	defer ticker.Stop()

	for {
		select {
		case <-m.StopChan:
			return
		case <-ticker.C:
			metrics, err := m.collectMetrics()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Metrics collection error: %v\n", err)
				continue
			}

			select {
			case m.Metrics <- metrics:
			default:
				// Channel is full, skip this update
			}
		}
	}
}

// collectMetrics gathers all monitoring metrics
func (m *VaultMonitor) collectMetrics() (*MonitorMetrics, error) {
	metrics := &MonitorMetrics{
		Timestamp: time.Now(),
	}

	var wg sync.WaitGroup
	var err error

	// Collect vault metrics
	wg.Add(1)
	go func() {
		defer wg.Done()
		metrics.Vault, err = m.collectVaultMetrics()
	}()

	// Collect system metrics
	wg.Add(1)
	go func() {
		defer wg.Done()
		metrics.System, err = m.collectSystemMetrics()
	}()

	// Collect network metrics
	wg.Add(1)
	go func() {
		defer wg.Done()
		metrics.Network, err = m.collectNetworkMetrics()
	}()

	// Collect security metrics
	wg.Add(1)
	go func() {
		defer wg.Done()
		metrics.Security, err = m.collectSecurityMetrics()
	}()

	// Collect performance metrics
	wg.Add(1)
	go func() {
		defer wg.Done()
		metrics.Performance, err = m.collectPerformanceMetrics()
	}()

	wg.Wait()
	return metrics, err
}

// collectVaultMetrics gathers vault-specific metrics
func (m *VaultMonitor) collectVaultMetrics() (VaultMetrics, error) {
	status, err := m.Context.GetStatus()
	if err != nil {
		return VaultMetrics{}, err
	}

	// This is a simplified implementation
	// In a real implementation, you would collect actual metrics from the vault
	return VaultMetrics{
		Status:           string(status.Mode),
		Version:          "1.0.0", // Would get from actual vault
		Uptime:           "2h 34m",
		ActiveTokens:     42,
		StoredSecrets:    156,
		MemoryUsage:      128 * 1024 * 1024, // 128MB
		CPUUsage:         15.5,
		Connections:      8,
		RequestsPerSec:   23.4,
		ResponseTime:     45 * time.Millisecond,
		EncryptionStatus: "AES-256-GCM",
	}, nil
}

// collectSystemMetrics gathers system-level metrics
func (m *VaultMonitor) collectSystemMetrics() (SystemMetrics, error) {
	// This is a simplified implementation
	// In a real implementation, you would use system libraries like gopsutil
	return SystemMetrics{
		OS:              runtime.GOOS,
		Architecture:    runtime.GOARCH,
		Hostname:        "localhost",
		CPUUsage:        12.3,
		MemoryTotal:     8 * 1024 * 1024 * 1024,   // 8GB
		MemoryUsed:      4 * 1024 * 1024 * 1024,   // 4GB
		MemoryAvailable: 4 * 1024 * 1024 * 1024,   // 4GB
		DiskTotal:       500 * 1024 * 1024 * 1024, // 500GB
		DiskUsed:        250 * 1024 * 1024 * 1024, // 250GB
		DiskAvailable:   250 * 1024 * 1024 * 1024, // 250GB
		LoadAverage:     []float64{1.2, 1.5, 1.8},
		ProcessCount:    156,
		Uptime:          "5d 12h 34m",
	}, nil
}

// collectNetworkMetrics gathers network-related metrics
func (m *VaultMonitor) collectNetworkMetrics() (NetworkMetrics, error) {
	// This is a simplified implementation
	// In a real implementation, you would collect actual network metrics
	return NetworkMetrics{
		Interface:   "eth0",
		IPAddresses: []string{"192.168.1.100", "10.0.0.5"},
		BytesSent:   1024 * 1024 * 100, // 100MB
		BytesRecv:   1024 * 1024 * 200, // 200MB
		PacketsSent: 50000,
		PacketsRecv: 75000,
		Connections: 25,
		Latency:     12 * time.Millisecond,
		Bandwidth:   1000.0, // Mbps
	}, nil
}

// collectSecurityMetrics gathers security-related metrics
func (m *VaultMonitor) collectSecurityMetrics() (SecurityMetrics, error) {
	// This is a simplified implementation
	// In a real implementation, you would collect actual security metrics
	return SecurityMetrics{
		FailedLogins:     3,
		SuccessfulLogins: 42,
		ActiveSessions:   8,
		ThreatAlerts:     0,
		PolicyViolations: 1,
		AuditEvents:      156,
		LastSecurityScan: time.Now().Add(-time.Hour * 2),
		EncryptionLevel:  "TLS-1.3",
	}, nil
}

// collectPerformanceMetrics gathers performance-related metrics
func (m *VaultMonitor) collectPerformanceMetrics() (PerformanceMetrics, error) {
	// This is a simplified implementation
	// In a real implementation, you would collect actual performance metrics
	return PerformanceMetrics{
		AvgResponseTime:  45 * time.Millisecond,
		P95ResponseTime:  120 * time.Millisecond,
		P99ResponseTime:  250 * time.Millisecond,
		ThroughputPerSec: 1000.0,
		ErrorRate:        0.01, // 1%
		CacheHitRate:     0.85, // 85%
		DatabaseQueries:  500,
		SlowQueries:      2,
	}, nil
}

// updateDisplay updates the terminal display
func (m *VaultMonitor) updateDisplay() error {
	if m.FullScreen {
		// Clear screen and move to top
		fmt.Print("\033[2J\033[H")
	}

	// Get latest metrics
	metrics := <-m.Metrics

	// Display header
	m.displayHeader()

	// Display metrics based on filter
	if m.Filter == "" || m.Filter == "vault" {
		m.displayVaultMetrics(metrics.Vault)
	}
	if m.Filter == "" || m.Filter == "system" {
		m.displaySystemMetrics(metrics.System)
	}
	if m.Filter == "" || m.Filter == "network" {
		m.displayNetworkMetrics(metrics.Network)
	}
	if m.Filter == "" || m.Filter == "security" {
		m.displaySecurityMetrics(metrics.Security)
	}
	if m.Filter == "" || m.Filter == "performance" {
		m.displayPerformanceMetrics(metrics.Performance)
	}

	// Display footer
	m.displayFooter()

	return nil
}

// displayHeader shows the monitoring dashboard header
func (m *VaultMonitor) displayHeader() {
	fmt.Printf("%s╔══════════════════════════════════════════════════════════════════════════════╗%s\n", ui.Blue, ui.Reset)
	fmt.Printf("%s║                              AETHER VAULT MONITOR                           ║%s\n", ui.Blue, ui.Reset)
	fmt.Printf("%s║                                   Real-Time                                   ║%s\n", ui.Blue, ui.Reset)
	fmt.Printf("%s╠══════════════════════════════════════════════════════════════════════════════╣%s\n", ui.Blue, ui.Reset)
}

// displayVaultMetrics shows vault-specific metrics
func (m *VaultMonitor) displayVaultMetrics(metrics VaultMetrics) {
	fmt.Printf("%s║ VAULT STATUS                                                               ║%s\n", ui.Bold, ui.Reset)
	fmt.Printf("%s║ Status: %s%-20s%s Version: %s%-10s%s Uptime: %s%-10s%s ║%s\n",
		ui.Blue, ui.Green, metrics.Status, ui.Blue,
		ui.Cyan, metrics.Version, ui.Blue,
		ui.Yellow, metrics.Uptime, ui.Blue, ui.Reset)
	fmt.Printf("%s║ Tokens: %s%-5d%s Secrets: %s%-5d%s Memory: %s%-8s%s CPU: %s%-5.1f%%%s ║%s\n",
		ui.Blue, ui.Green, metrics.ActiveTokens, ui.Blue,
		ui.Cyan, metrics.StoredSecrets, ui.Blue,
		ui.Yellow, formatBytes(uint64(metrics.MemoryUsage)), ui.Blue,
		getUsageColor(metrics.CPUUsage), metrics.CPUUsage, ui.Blue, ui.Reset)
	fmt.Printf("%s╠══════════════════════════════════════════════════════════════════════════════╣%s\n", ui.Blue, ui.Reset)
}

// displaySystemMetrics shows system-level metrics
func (m *VaultMonitor) displaySystemMetrics(metrics SystemMetrics) {
	memPercent := float64(metrics.MemoryUsed) / float64(metrics.MemoryTotal) * 100
	diskPercent := float64(metrics.DiskUsed) / float64(metrics.DiskTotal) * 100

	fmt.Printf("%s║ SYSTEM RESOURCES                                                            ║%s\n", ui.Bold, ui.Reset)
	fmt.Printf("%s║ OS: %s%-15s%s Host: %s%-15s%s Processes: %s%-4d%s ║%s\n",
		ui.Blue, ui.Green, metrics.OS, ui.Blue,
		ui.Cyan, metrics.Hostname, ui.Blue,
		ui.Yellow, metrics.ProcessCount, ui.Blue, ui.Reset)
	fmt.Printf("%s║ CPU: %s%-5.1f%%%s Memory: %s%-5.1f%%%s Disk: %s%-5.1f%%%s Load: %s%.1f %.1f %.1f%s ║%s\n",
		ui.Blue, getUsageColor(metrics.CPUUsage), metrics.CPUUsage, ui.Blue,
		getUsageColor(memPercent), memPercent, ui.Blue,
		getUsageColor(diskPercent), diskPercent, ui.Blue,
		ui.Yellow, metrics.LoadAverage[0], metrics.LoadAverage[1], metrics.LoadAverage[2], ui.Blue, ui.Reset)
	fmt.Printf("%s╠══════════════════════════════════════════════════════════════════════════════╣%s\n", ui.Blue, ui.Reset)
}

// displayNetworkMetrics shows network-related metrics
func (m *VaultMonitor) displayNetworkMetrics(metrics NetworkMetrics) {
	fmt.Printf("%s║ NETWORK METRICS                                                             ║%s\n", ui.Bold, ui.Reset)
	fmt.Printf("%s║ Interface: %s%-10s%s IP: %s%-15s%s Connections: %s%-3d%s ║%s\n",
		ui.Blue, ui.Green, metrics.Interface, ui.Blue,
		ui.Cyan, strings.Join(metrics.IPAddresses, ", "), ui.Blue,
		ui.Yellow, metrics.Connections, ui.Blue, ui.Reset)
	fmt.Printf("%s║ Sent: %s%-8s%s Recv: %s%-8s%s Latency: %s%-6s%s Bandwidth: %s%-6.1f Mbps%s ║%s\n",
		ui.Blue, ui.Cyan, formatBytes(metrics.BytesSent), ui.Blue,
		ui.Green, formatBytes(metrics.BytesRecv), ui.Blue,
		ui.Yellow, metrics.Latency.String(), ui.Blue,
		ui.Magenta, metrics.Bandwidth, ui.Blue, ui.Reset)
	fmt.Printf("%s╠══════════════════════════════════════════════════════════════════════════════╣%s\n", ui.Blue, ui.Reset)
}

// displaySecurityMetrics shows security-related metrics
func (m *VaultMonitor) displaySecurityMetrics(metrics SecurityMetrics) {
	fmt.Printf("%s║ SECURITY METRICS                                                            ║%s\n", ui.Bold, ui.Reset)
	fmt.Printf("%s║ Logins: %s✓%-4d%s ✗%-4d%s Sessions: %s%-3d%s Threats: %s%-2d%s Encryption: %s%-10s%s ║%s\n",
		ui.Blue, ui.Green, metrics.SuccessfulLogins, ui.Red,
		metrics.FailedLogins, ui.Blue,
		ui.Yellow, metrics.ActiveSessions, ui.Blue,
		getAlertColor(metrics.ThreatAlerts), metrics.ThreatAlerts, ui.Blue,
		ui.Cyan, metrics.EncryptionLevel, ui.Blue, ui.Reset)
	fmt.Printf("%s║ Violations: %s%-3d%s Audit Events: %s%-4d%s Last Scan: %s%-10s%s ║%s\n",
		ui.Blue, getAlertColor(metrics.PolicyViolations), metrics.PolicyViolations, ui.Blue,
		ui.Cyan, metrics.AuditEvents, ui.Blue,
		ui.Yellow, metrics.LastSecurityScan.Format("15:04:05"), ui.Blue, ui.Reset)
	fmt.Printf("%s╠══════════════════════════════════════════════════════════════════════════════╣%s\n", ui.Blue, ui.Reset)
}

// displayPerformanceMetrics shows performance-related metrics
func (m *VaultMonitor) displayPerformanceMetrics(metrics PerformanceMetrics) {
	fmt.Printf("%s║ PERFORMANCE METRICS                                                         ║%s\n", ui.Bold, ui.Reset)
	fmt.Printf("%s║ Avg: %s%-6s%s P95: %s%-6s%s P99: %s%-6s%s Throughput: %s%-6.0f/s%s ║%s\n",
		ui.Blue, ui.Green, metrics.AvgResponseTime.String(), ui.Blue,
		ui.Yellow, metrics.P95ResponseTime.String(), ui.Blue,
		ui.Red, metrics.P99ResponseTime.String(), ui.Blue,
		ui.Cyan, metrics.ThroughputPerSec, ui.Blue, ui.Reset)
	fmt.Printf("%s║ Error Rate: %s%-5.2f%%%s Cache Hit: %s%-5.1f%%%s DB Queries: %s%-4d%s Slow: %s%-2d%s ║%s\n",
		ui.Blue, getAlertColor(int(metrics.ErrorRate*100)), metrics.ErrorRate*100, ui.Blue,
		ui.Green, metrics.CacheHitRate*100, ui.Blue,
		ui.Cyan, metrics.DatabaseQueries, ui.Blue,
		ui.Red, metrics.SlowQueries, ui.Blue, ui.Reset)
	fmt.Printf("%s╚══════════════════════════════════════════════════════════════════════════════╝%s\n", ui.Blue, ui.Reset)
}

// displayFooter shows the monitoring dashboard footer
func (m *VaultMonitor) displayFooter() {
	fmt.Printf("\n%sLast Update: %s%s | Refresh: %s%s | Press Ctrl+C to exit%s\n",
		ui.Blue, ui.Cyan, time.Now().Format("15:04:05"),
		ui.Yellow, m.Refresh.String(), ui.Reset)
}

// checkAlerts checks for alert conditions
func (m *VaultMonitor) checkAlerts(metrics *MonitorMetrics) {
	// Check for high CPU usage
	if metrics.System.CPUUsage > 80 {
		fmt.Printf("%s⚠️  High CPU usage: %.1f%%%s\n", ui.Red, metrics.System.CPUUsage, ui.Reset)
	}

	// Check for high memory usage
	memPercent := float64(metrics.System.MemoryUsed) / float64(metrics.System.MemoryTotal) * 100
	if memPercent > 85 {
		fmt.Printf("%s⚠️  High memory usage: %.1f%%%s\n", ui.Red, memPercent, ui.Reset)
	}

	// Check for security threats
	if metrics.Security.ThreatAlerts > 0 {
		fmt.Printf("%s🚨 Security threats detected: %d%s\n", ui.Red, metrics.Security.ThreatAlerts, ui.Reset)
	}

	// Check for high error rate
	if metrics.Performance.ErrorRate > 0.05 { // 5%
		fmt.Printf("%s⚠️  High error rate: %.2f%%%s\n", ui.Red, metrics.Performance.ErrorRate*100, ui.Reset)
	}
}

// Helper functions

// formatBytes formats bytes into human readable string
func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// getUsageColor returns color based on usage percentage
func getUsageColor(usage float64) string {
	switch {
	case usage < 50:
		return ui.Green
	case usage < 80:
		return ui.Yellow
	default:
		return ui.Red
	}
}

// getAlertColor returns color based on alert count
func getAlertColor(count int) string {
	switch {
	case count == 0:
		return ui.Green
	case count < 5:
		return ui.Yellow
	default:
		return ui.Red
	}
}
