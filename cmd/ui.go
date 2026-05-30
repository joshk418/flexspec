/*
Copyright © 2026 Josh Kyte
*/
package cmd

import (
	"context"
	"fmt"
	"io/fs"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/joshk418/flexspec/internal/ui"
)

// UIFS holds the embedded web UI assets (web/dist).
var UIFS fs.FS

var (
	uiPort   int
	uiHost   string
	uiOpen   bool
	uiNoOpen bool
)

var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Start the local management UI in your browser",
	Long: `Start a local HTTP server that serves the FlexSpec management UI and API.

The UI provides a kanban/table board, spec browser, and settings editor.
Changes to spec files on disk refresh automatically via server-sent events.`,
	RunE: runUI,
}

func runUI(cmd *cobra.Command, _ []string) error {
	root, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("resolve working directory: %w", err)
	}

	static := UIFS
	if static == nil {
		static = ui.StubStaticFS()
	}

	host := uiHost
	port := uiPort
	var srv *ui.Server
	for i := 0; i < 21; i++ {
		tryPort := port + i
		srv, err = ui.NewServer(root, host, tryPort, static)
		if err != nil {
			return err
		}
		ln, listenErr := net.Listen("tcp", srv.Addr())
		if listenErr == nil {
			_ = ln.Close()
			if i > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "Port %d in use, using %d\n", port, tryPort)
			}
			break
		}
		if i == 20 {
			return fmt.Errorf("no available port in range %d-%d", port, port+20)
		}
	}

	url := fmt.Sprintf("http://%s", srv.Addr())
	fmt.Fprintf(cmd.OutOrStdout(), "FlexSpec UI at %s\n", url)
	fmt.Fprintf(cmd.OutOrStdout(), "Press Ctrl+C to stop\n")

	shouldOpen := uiOpen && !uiNoOpen
	if shouldOpen {
		if err := openBrowser(url); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "could not open browser: %v\n", err)
		}
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return srv.Run(ctx)
}

func openBrowser(url string) error {
	var c *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		c = exec.Command("open", url)
	case "windows":
		c = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		c = exec.Command("xdg-open", url)
	}
	return c.Start()
}

func init() {
	rootCmd.AddCommand(uiCmd)
	uiCmd.Flags().IntVarP(&uiPort, "port", "p", 3000, "Preferred HTTP port")
	uiCmd.Flags().StringVar(&uiHost, "host", "127.0.0.1", "HTTP listen host")
	uiCmd.Flags().BoolVar(&uiOpen, "open", true, "Open the UI in the default browser")
	uiCmd.Flags().BoolVar(&uiNoOpen, "no-open", false, "Do not open the browser")
}
