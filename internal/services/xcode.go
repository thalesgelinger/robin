package services

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/thalesgelinger/robin/internal/entities"
)

type Xcode struct {
	ProjectPath string
	AppName     string
	Scheme      string
	ArchivePath string
	ExportPath  string
	ExportPlist string
	LogPath     string
	Verbose     bool
}

func NewXcode(ios entities.IOS, env entities.Environment, credentials *Credentials, verbose bool) (*Xcode, error) {
	archivePath := filepath.Join(ios.OutputDir, "Archives", fmt.Sprintf("%s.xcarchive", ios.AppName))
	exportPath := filepath.Join(ios.OutputDir, "IPA")
	logPath := filepath.Join(ios.OutputDir, "Logs")

	// Ensure directories exist
	os.MkdirAll(filepath.Dir(archivePath), os.ModePerm)
	os.MkdirAll(exportPath, os.ModePerm)
	os.MkdirAll(logPath, os.ModePerm)

	// Generate ExportOptions.plist dynamically
	exportPlistPath := filepath.Join(ios.ProjectPath, "ExportOptions.plist")
	if err := generateExportOptionsPlist(exportPlistPath, env.ExportMethod, credentials.TeamID); err != nil {
		return nil, err
	}

	return &Xcode{
		ProjectPath: ios.ProjectPath,
		AppName:     ios.AppName,
		Scheme:      env.Scheme,
		ArchivePath: archivePath,
		ExportPath:  exportPath,
		ExportPlist: "./ExportOptions.plist",
		LogPath:     generateTempLogPath(),
		Verbose:     verbose,
	}, nil
}

// Generate a log file path inside /tmp
func generateTempLogPath() string {
	timestamp := time.Now().Format("20060102-150405")
	return filepath.Join(os.TempDir(), fmt.Sprintf("xcodebuild-%s.log", timestamp))
}

func (x *Xcode) Archive() error {
	cmd := exec.Command("xcodebuild",
		"-workspace", fmt.Sprintf("%s.xcworkspace", x.AppName),
		"-scheme", x.Scheme,
		"-destination", "generic/platform=iOS",
		"-archivePath", x.ArchivePath,
		"archive",
	)

	cmd.Dir = x.ProjectPath

	fmt.Println("📦 Archiving the project...")
	return x.runAndLog(cmd, "❌ Failed to archive")
}

func (x *Xcode) ExportIPA() error {
	cmd := exec.Command("xcodebuild",
		"-exportArchive",
		"-archivePath", x.ArchivePath,
		"-exportPath", x.ExportPath,
		"-exportOptionsPlist", x.ExportPlist,
	)

	cmd.Dir = x.ProjectPath // Change directory before running the command

	fmt.Println("📦 Exporting IPA...")
	return x.runAndLog(cmd, "❌ Failed to export IPA")
}

// Run a command and log errors
func (x *Xcode) runAndLog(cmd *exec.Cmd, errMsg string) error {
	logFile, err := os.Create(x.LogPath)
	if err != nil {
		return fmt.Errorf("❌ Could not create log file: %w", err)
	}
	defer logFile.Close()

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("❌ Failed to get stdout pipe: %w", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("❌ Failed to get stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("❌ Failed to start command: %w", err)
	}

	// Capture logs and write to file
	go func() {
		scanner := bufio.NewScanner(io.MultiReader(stdoutPipe, stderrPipe))
		for scanner.Scan() {
			line := scanner.Text()
			logFile.WriteString(line + "\n")

			if x.Verbose {
				fmt.Println(line) // Print logs in real-time if verbose mode is on
			}
		}
	}()

	if err := cmd.Wait(); err != nil {
		fmt.Printf("%s. Check logs: %s\n", errMsg, x.LogPath)
		return fmt.Errorf("%s: %w", errMsg, err)
	}

	fmt.Println("✅ Process completed successfully!")
	return nil
}

// Generate ExportOptions.plist dynamically
func generateExportOptionsPlist(path, method, teamID string) error {
	content := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>method</key>
    <string>%s</string>
    <key>teamID</key>
    <string>%s</string>
    <key>compileBitcode</key>
    <true/>
    <key>destination</key>
    <string>export</string>
    <key>signingStyle</key>
    <string>automatic</string>
    <key>stripSwiftSymbols</key>
    <true/>
</dict>
</plist>`, method, teamID)

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("❌ Failed to create ExportOptions.plist: %w", err)
	}

	fmt.Println("✅ ExportOptions.plist created successfully!")
	return nil
}
