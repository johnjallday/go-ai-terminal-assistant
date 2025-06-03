package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"go-ai-terminal-assistant/server"
)

const (
	defaultPort    = "8080"
	defaultLogFile = "/usr/local/var/log/ai-terminal-assistant.log"
	pidFile        = "/usr/local/var/run/ai-terminal-assistant.pid"
)

func main() {
	var (
		port      = flag.String("port", defaultPort, "Port to run the server on")
		logFile   = flag.String("log", defaultLogFile, "Log file path")
		daemon    = flag.Bool("daemon", false, "Run as daemon")
		install   = flag.Bool("install", false, "Install as macOS service")
		uninstall = flag.Bool("uninstall", false, "Uninstall macOS service")
		start     = flag.Bool("start", false, "Start the service")
		stop      = flag.Bool("stop", false, "Stop the service")
		status    = flag.Bool("status", false, "Check service status")
	)
	flag.Parse()

	// Handle service management commands
	if *install {
		if err := installService(*port, *logFile); err != nil {
			log.Fatalf("Failed to install service: %v", err)
		}
		fmt.Println("✅ AI Terminal Assistant service installed successfully!")
		fmt.Printf("💡 Start with: ai-terminal-assistant -start\n")
		return
	}

	if *uninstall {
		if err := uninstallService(); err != nil {
			log.Fatalf("Failed to uninstall service: %v", err)
		}
		fmt.Println("✅ AI Terminal Assistant service uninstalled successfully!")
		return
	}

	if *start {
		if err := startService(); err != nil {
			log.Fatalf("Failed to start service: %v", err)
		}
		fmt.Println("✅ AI Terminal Assistant service started!")
		fmt.Printf("📡 API available at: http://localhost:%s\n", *port)
		return
	}

	if *stop {
		if err := stopService(); err != nil {
			log.Fatalf("Failed to stop service: %v", err)
		}
		fmt.Println("✅ AI Terminal Assistant service stopped!")
		return
	}

	if *status {
		running, err := isServiceRunning()
		if err != nil {
			log.Fatalf("Failed to check service status: %v", err)
		}
		if running {
			fmt.Println("✅ AI Terminal Assistant service is running")
			if pid := getPID(); pid != "" {
				fmt.Printf("📊 PID: %s\n", pid)
			}
			fmt.Printf("📡 API available at: http://localhost:%s\n", *port)
		} else {
			fmt.Println("❌ AI Terminal Assistant service is not running")
		}
		return
	}

	// Setup logging for daemon mode
	if *daemon {
		if err := setupLogging(*logFile); err != nil {
			log.Fatalf("Failed to setup logging: %v", err)
		}

		// Write PID file
		if err := writePIDFile(); err != nil {
			log.Printf("Warning: Failed to write PID file: %v", err)
		}
		defer removePIDFile()
	}

	// Create and start the server
	srv, err := server.NewDaemonServer(*port)
	if err != nil {
		log.Fatalf("Failed to create daemon server: %v", err)
	}

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal %v, shutting down gracefully...", sig)
		cancel()
	}()

	// Start server in a goroutine
	go func() {
		log.Printf("Starting AI Terminal Assistant daemon on port %s", *port)
		if err := srv.Start(ctx); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	log.Println("AI Terminal Assistant daemon stopped")
}

func setupLogging(logFile string) error {
	// Create log directory if it doesn't exist
	logDir := filepath.Dir(logFile)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open log file
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	// Set log output to file
	log.SetOutput(file)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	return nil
}

func writePIDFile() error {
	pidDir := filepath.Dir(pidFile)
	if err := os.MkdirAll(pidDir, 0755); err != nil {
		return fmt.Errorf("failed to create PID directory: %w", err)
	}

	file, err := os.Create(pidFile)
	if err != nil {
		return fmt.Errorf("failed to create PID file: %w", err)
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%d\n", os.Getpid())
	return err
}

func removePIDFile() {
	os.Remove(pidFile)
}

func getPID() string {
	data, err := os.ReadFile(pidFile)
	if err != nil {
		return ""
	}
	return string(data)
}
