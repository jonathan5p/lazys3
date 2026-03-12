package main

import (
	"context"
	"log"

	"github.com/jonathan5p/lazys3/cmd/lazys3-tui/logger"
	"github.com/jonathan5p/lazys3/cmd/lazys3-tui/s3"
	"github.com/jonathan5p/lazys3/cmd/lazys3-tui/ui"
)

func main() {
	ctx := context.Background()

	appLog, err := logger.New("lazys3.log")
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer appLog.Close()

	s3Client, err := s3.NewClient(ctx, appLog)
	if err != nil {
		log.Fatalf("Failed to create S3 client: %v", err)
	}

	app := ui.NewApp(s3Client, appLog)
	if err := app.Run(ctx); err != nil {
		log.Fatalf("TUI error: %v", err)
	}
}
