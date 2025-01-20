package main

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
	"github.com/dustin/go-humanize"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"
)

func logAnd[T any](fn func(line []byte) (T, error)) func(line []byte) {
	return func(line []byte) {
		t, e := fn(line)
		fmt.Println(t, e)
	}
}

func unmarshalThen[T any, R any](unmarshaller func(line []byte) (T, error), then func(val T, line []byte) R) func(line []byte) R {
	return func(line []byte) R {
		t, e := unmarshaller(line)
		if e != nil {
			log.Panic("when unmarshalling", e)
		}

		return then(t, line)
	}
}

func actionUntypedMessage(m UntypedMessage, rest []byte) *SummaryMessage {
	switch ty := m.MessageType; ty {
	case "status":
		status, err := UnmarshalStatusMessage(rest)
		if err != nil {
			log.Panic(err)
		}
		
		elapsed, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(status.SecondsElapsed)))
		remaining, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(status.SecondsRemaining)))
		
		slog.Info("Status update",
			"done", fmt.Sprintf("%.2f%%", status.PercentDone*100),
			"elapsed", elapsed,
			"remaining", remaining,
			"files-done", status.FilesDone,
			"files-total", status.TotalFiles,
			"current-files", status.CurrentFiles,
		)
	case "summary":
		summary, err := UnmarshalSummaryMessage(rest)
		if err != nil {
			log.Panic(err)
		}

		duration, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(summary.TotalDuration)))
		processed := humanize.Bytes(uint64(summary.TotalBytesProcessed))

		slog.Info("Completed backup",
			"processed", processed,
			"took", duration,
			"snapshot", summary.SnapshotID)

		slog.Info("File summary",
			"files-new", summary.FilesNew,
			"files-changed", summary.FilesChanged,
			"files-unchanged", summary.FilesUnmodified)

		slog.Info("Directory summary",
			"dirs-new", summary.DirsNew,
			"dirs-changed", summary.DirsChanged,
			"dirs-unchanged", summary.DirsUnmodified)

		return &summary

	default:
		slog.Warn("unsupported action type", "type", ty)
	}

	return nil
}

func consume(line []byte) *string {
	slog.Debug(string(line))
	s := "unused"
	return &s
}

func sendResult(logs string, success bool, summary *SummaryMessage) {
	client, err := webhook.NewWithURL(os.Getenv("AUTORESTIC_WEBHOOK"))
	if err != nil {
		fmt.Println("Tried to send a result with no destination, sending here!")
		fmt.Println("success=", success, "logs=", logs)
		os.Exit(1)
	}

	colorFail := 0xf38ba8
	colorSuccess := 0xa6e3a1
	color := colorSuccess

	if !success {
		color = colorFail
	}

	embed := discord.NewEmbedBuilder().
		SetTitle("Restic backup summary").
		SetColor(color)

	if summary != nil {
		duration, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(summary.TotalDuration)))
		processed := humanize.Bytes(uint64(summary.TotalBytesProcessed))
		embed.
			AddField("Processed", processed, true).
			AddField("Duration", fmt.Sprintf("%s", duration), true).
			AddField("ID", summary.SnapshotID, true)
	} else {
		embed.SetDescription("No summary available")
	}

	logFile := discord.NewFile("restic.log", "the restic log file", strings.NewReader(logs))
	message := discord.NewWebhookMessageCreateBuilder().
		AddEmbeds(embed.Build()).
		AddFiles(logFile)

	if !success {
		message.SetContent("@everyone restic backup failure")
	}

	_, err = client.CreateMessage(message.Build())

	if err != nil {
		fmt.Println("Failed to upload logs", err)
		fmt.Println(logs)
	}
}

func doBackup(location string) {
	capturedLogs := new(strings.Builder)

	defer func() {
		if r := recover(); r != nil {
			slog.Error("Encountered an error when performing backup", "err", r)
			sendResult(capturedLogs.String(), false, nil)
			os.Exit(1)
		}
	}()

	logger := slog.New(
		NewTeeLogger(
			slog.NewTextHandler(capturedLogs, &slog.HandlerOptions{Level: slog.LevelInfo}),
			slog.NewTextHandler(
				os.Stderr,
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				},
			),
		),
	)
	slog.SetDefault(logger)

	fps := os.Getenv("RESTIC_PROGRESS_FPS")
	if fps == "" {
		slog.Warn("No RESTIC_PROGRESS_FPS set, defaulting to 0.016666 (1 / minute)")
		os.Setenv("RESTIC_PROGRESS_FPS", "0.016666")
	}

	_, err := invokeCommand(consume, "restic", "cat", "config")
	if err != nil {
		slog.Error("restic repo invalid, could not cat config", "err", err)
		return
	}

	ignore := os.Getenv("AUTORESTIC_IGNORE")

	slog.Info("Backing up", "dir", location)
	args := []string{"backup", location}
	if ignore != "" {
		args = append(args, "--exclude")
		args = append(args, ignore)
	}
	
	summary, _ := invokeCommand(
		unmarshalThen(UnmarshalUntypedMessage, actionUntypedMessage),
		"restic", args...,
	)

	sendResult(capturedLogs.String(), true, summary)
}

func main() {
	if os.Getenv("AUTORESTIC_WEBHOOK") == "" {
		log.Panic("Cannot proceed without a webhook to send to, set AUTORESTIC_WEBHOOK")
	}

	locationsRaw := os.Getenv("AUTORESTIC_LOCATIONS")
	if locationsRaw == "" {
		log.Panic("Cannot proceed without locations to backup, set AUTORESTIC_LOCATIONS to a comma separated list of directories")
	}

	locations := strings.Split(locationsRaw, ",")
	if len(locations) == 0 {
		log.Panic("Cannot proceed with no valid locations set AUTORESTIC_LOCATIONS to a comma separated list of directories")
	}

	for _, location := range locations {
		doBackup(location)
	}
}
