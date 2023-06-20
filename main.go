package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/prashantv/tmux-editor-popup/prompt"
)

const name = "tmux-editor-popup"

func main() {
	var (
		flagSkipPrompt = flag.Bool("skip-prompt", false, "Whether to skip prompt detection, and only use the most recent file for the command")
		flagDir        = flag.String("dir", "/tmp/snippets", "Directory to use for storing entered snippets")
	)
	flag.Parse()

	opts := Opts{
		Dir:       *flagDir,
		UsePrompt: !*flagSkipPrompt,
	}

	if err := opts.Run(); err != nil {
		// tmux shows stdout but not stderr.
		fmt.Printf("%v failed: %v\n", name, err)
		os.Exit(1)
	}
}

type Opts struct {
	Dir       string
	UsePrompt bool
}

func (o Opts) Run() error {
	// sendBackspaces is used when we're replacing an existing prompt
	var sendBackspaces int

	cmd := getCommand()
	file := o.getFile(cmd)

	if err := os.MkdirAll(filepath.Dir(file), 0o777); err != nil {
		return fmt.Errorf("create dir for snippet: %v", err)
	}

	if o.UsePrompt {
		promptParser, err := prompt.NewAll()
		if err != nil {
			return fmt.Errorf("prompt parse config: %v", err)
		}

		p, err := o.findPrompt(promptParser)
		if err == nil {
			preEditContent := strings.Join(p.Parsed, "\n")
			if p.IsLast {
				sendBackspaces = len(preEditContent)
			}

			if err := os.WriteFile(file, []byte(preEditContent), 0o666); err != nil {
				tmuxShowMessage("Failed to write prompt contents to file")
			}

		} else {
			go tmuxShowMessage("No prompt found in scrollback, reusing last snippet for command %q", cmd)
		}
	}

	preEditorModTime := modTimeOr0(file)

	if err := tmuxShowPopup(o.Dir, fmt.Sprintf("$EDITOR %q", file)); err != nil {
		return fmt.Errorf("tmux display-popup with editor: %v", err)
	}

	postEditorModTime := modTimeOr0(file)

	if preEditorModTime == postEditorModTime {
		// File was not written, so drop.
		_ = tmuxShowMessage("editor contents not modified, not sending any input")
		return nil
	}

	// Otherwise, send the keys.
	contents, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("read temporary file: %v", err)
	}

	sendKeys := string(contents)
	if sendBackspaces > 0 {
		sendKeys = strings.Repeat("\b", sendBackspaces) + sendKeys
	}

	if _, err := runCmd("tmux", "send-keys",
		"-l", // read values as literal values, without key recognition like Enter => \n
		sendKeys,
	); err != nil {
		return fmt.Errorf("send back-space to remove current prompt: %v", err)
	}

	return nil
}

func (o Opts) findPrompt(promptParser *prompt.Parser) (*prompt.Prompt, error) {
	var lastPrompt *prompt.Prompt
	var lastPane []byte

	for start := 0; start > -10000; start -= 1000 {
		pane, err := tmuxCapturePane(start)
		if err != nil {
			return nil, err
		}

		if bytes.Equal(pane, lastPane) {
			// Going back in history didn't give us anything additional
			// so let's return the last match (if any).
			if lastPrompt == nil {
				break
			}

			return lastPrompt, nil
		}

		prompt, err := promptParser.Find(bytes.NewReader(pane))
		if err != nil {
			// We shouldn't get any errors, so retry with more screen.
			log.Printf("prompt parser got err: %v", err)
			continue
		}

		if prompt.IsFirst {
			// Retry with more history
			lastPrompt = &prompt
			continue
		}

		return &prompt, nil
	}

	return nil, errors.New("not found")
}

func getCommand() string {
	currentCommand, err := runCmd("tmux", "display",
		"-p",
		"#{pane_current_command}",
	)
	if err != nil {
		return "unknown"
	}

	return string(bytes.TrimSpace(currentCommand))
}

func (o Opts) getFile(cmd string) string {
	return filepath.Join(o.Dir, cmd+".last")
}

func tmuxShowMessage(msg string, args ...any) error {
	_, err := runCmd("tmux", "display-message", fmt.Sprintf(name+": "+msg, args...))
	return err
}

// start should be negative to go backwards.
func tmuxCapturePane(start int) ([]byte, error) {
	return runCmd("tmux", "capture-pane",
		"-p",                    // capture to stdout
		"-J",                    // joins wrapped lines
		"-S", fmt.Sprint(start), // Where to start capturing. 0 = first line, negative is history.
	)
}

func tmuxShowPopup(startDir, shellCmd string) error {
	if _, err := runCmd("tmux", "display-popup",
		"-d", startDir, // start-directory
		"-E", // closes popup automatically when the command exits.
		shellCmd,
	); err != nil {
		return fmt.Errorf("tmux display-popup: %v", err)
	}

	return nil
}

func modTimeOr0(path string) int64 {
	s, err := os.Stat(path)
	if err != nil {
		return 0
	}

	return s.ModTime().UnixNano()
}

func runCmd(command string, args ...string) (out []byte, _ error) {
	var stdout bytes.Buffer
	cmd := exec.Command(command, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return stdout.Bytes(), nil
}
