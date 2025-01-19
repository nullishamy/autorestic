package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os/exec"
	"strings"
)

func invokeCommand[R any](transformer func(line []byte) *R, cmd string, args ...string) (*R, error) {
	slog.Info("Invoking command", "cmd", cmd, "args", args)
	command := exec.Command(cmd, append(args, "--json")...)
	stdout, err := command.StdoutPipe()

	if err != nil {
		log.Fatal(err)
	}

	stderr, err := command.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := command.Start(); err != nil {
		log.Fatal(err)
	}

	var lastTransformer *R
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		slog.Debug("Got line", "line", line)
		lastTransformer = transformer([]byte(line))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	buf := new(strings.Builder)
	_, e := io.Copy(buf, stderr)
	if e != nil {
		slog.Warn("could not capture stderr for command")
	}

	if err := command.Wait(); err != nil {
		return nil, errors.New(fmt.Sprintf("%s, %s", buf.String(), err))
	}

	return lastTransformer, nil
}
