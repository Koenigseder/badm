package packages

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/creack/pty"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

// ReadConfigFile reads and returns the config file .badm.yaml as ConfigFile
func ReadConfigFile(cfgFilePath string) (*ConfigFile, error) {
	data, err := os.ReadFile(cfgFilePath)
	if err != nil {
		return nil, err
	}

	cfgFile := new(ConfigFile)

	err = yaml.Unmarshal(data, cfgFile)
	if err != nil {
		return nil, err
	}

	return cfgFile, nil
}

func (c *ConfigFile) InstallPackages(installScopeName string) error {
	// Check if install scope name exists
	installScope, exists := c.Packages.InstallScopes[installScopeName]
	if !exists {
		return fmt.Errorf("install scope '%s' does not exist", installScopeName)
	}

	// First, install dependent packages
	dependsOn := installScope.DependsOn
	if dependsOn == installScopeName {
		return errors.New("cyclic dependencies are not allowed")
	}

	if dependsOn != "" {
		if err := c.InstallPackages(dependsOn); err != nil {
			return fmt.Errorf("failed installing depending scope '%s': %v", dependsOn, err)
		}
	}

	for packageManagerName, packageNames := range installScope.Packages {
		// Check if package manager name exists
		packageManagerCmd, exists := c.Packages.PackageManagers[packageManagerName]
		if !exists {
			return fmt.Errorf("package manager alias '%s' does not exist", packageManagerName)
		}

		cmdName, cmdArgs := parseCommandString(packageManagerCmd)

		// Construct command
		cmd := exec.Command(cmdName, append(cmdArgs, packageNames...)...)

		err := executeCommandInPTY(cmd)
		if err != nil {
			return fmt.Errorf("failed installing packages for '%s': %v", packageManagerName, err)
		}
	}

	// Execute scripts if present
	if installScope.Scripts == nil {
		return nil
	}

	// Create temp script directory
	tempDir := fmt.Sprintf("%s/.temp", c.Directory)

	err := os.MkdirAll(tempDir, 0777)
	if err != nil {
		return fmt.Errorf("error creating temporary directory at '%s': %v", tempDir, err)
	}

	// ... and defer delete it
	defer func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			fmt.Printf("error removing temporary directory at '%s': %v", tempDir, err)
		}
	}()

	// Execute each script
	for _, scriptName := range installScope.Scripts {
		script, exists := c.Scripts[scriptName]
		if !exists {
			return fmt.Errorf("script '%s' does not exist", scriptName)
		}

		scriptPath := fmt.Sprintf("%s/%s", tempDir, scriptName)

		err := os.WriteFile(scriptPath, []byte(script.Exec), 0777)
		if err != nil {
			return fmt.Errorf("error writing temporary script at '%s': %v", scriptPath, err)
		}

		// Execute in PTY
		cmd := exec.Command(script.Shell, scriptPath)

		err = executeCommandInPTY(cmd)
		if err != nil {
			return fmt.Errorf("error executing temporary script at '%s': %v", scriptName, err)
		}
	}

	return nil
}

func executeCommandInPTY(cmd *exec.Cmd) error {
	// Hand over to PTY
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return fmt.Errorf("failed starting pseudo terminal: %v", err)
	}

	defer ptmx.Close()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				fmt.Printf("error resizing pty: %v\n", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH
	defer func() { signal.Stop(ch); close(ch) }()

	// Set stdin in raw mode.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	// Copy stdin to the pty and the pty to stdout.
	// NOTE: The goroutine will keep reading until the next keystroke before returning.
	go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
	_, _ = io.Copy(os.Stdout, ptmx)

	return nil
}

func parseCommandString(cmdString string) (string, []string) {
	cmdParts := strings.Split(cmdString, " ")

	if len(cmdParts) == 1 {
		return cmdParts[0], nil
	}

	return cmdParts[0], cmdParts[1:]
}
