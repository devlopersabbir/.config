package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"xauusd/internal/analysis"
	"xauusd/internal/market"
	"xauusd/internal/streamer"
	"xauusd/internal/ui"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	log.Println("======================================")
	log.Println(" XAUUSD / Gold Live Market Daemon")
	log.Println("======================================")

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	// Initialize state and modular components
	state := market.NewState()
	analyzer := analysis.NewAnalyzer(state)
	streamerClient := streamer.NewStreamer(state)
	renderer := ui.NewRenderer(state)

	// Launch background workers
	go analyzer.Start(ctx)
	go renderer.Start(ctx)
	go streamerClient.Start(ctx)

	// Await graceful termination
	<-ctx.Done()

	log.Println("Shutting down Gold market daemon...")
}
