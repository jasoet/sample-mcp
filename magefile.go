//go:build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/mg"
	"os"
	"os/exec"
	"runtime"
	"time"
)

var Default = Build

func Build() error {
	fmt.Println("Building with GoReleaser...")

	if err := ensureToolInstalled("goreleaser", "github.com/goreleaser/goreleaser@latest"); err != nil {
		return err
	}

	cmd := exec.Command("goreleaser", "release", "--snapshot", "--clean")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func CompileBuild() error {
	fmt.Println("Building with standard Go build...")

	// Create dist directory if it doesn't exist
	if err := os.MkdirAll("dist", 0755); err != nil {
		return fmt.Errorf("failed to create dist directory: %w", err)
	}

	// Determine target OS and architecture
	goos := getEnvOrDefault("GOOS", runtime.GOOS)
	goarch := getEnvOrDefault("GOARCH", runtime.GOARCH)

	// Build for the target platform
	outputPath := fmt.Sprintf("mcp-server")
	if goos == "windows" {
		outputPath += ".exe"
	}

	fmt.Printf("Building for %s/%s to %s\n", goos, goarch, outputPath)

	cmd := exec.Command("go", "build", "-o", outputPath, "./main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Test() error {
	fmt.Println("Running tests...")
	cmd := exec.Command("go", "test", "-count=1", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func IntegrationTest() error {
	fmt.Println("Running integration tests...")

	docker := Docker{}
	if err := docker.Up(); err != nil {
		return fmt.Errorf("failed to start docker services: %w", err)
	}

	// Give PostgreSQL time to initialize
	fmt.Println("Waiting for PostgreSQL to initialize...")
	time.Sleep(2 * time.Second)

	cmd := exec.Command("go", "test", "-count=1", "-tags=integration", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Lint() error {
	fmt.Println("Running linter...")

	if err := ensureToolInstalled("golangci-lint", "github.com/golangci/golangci-lint/cmd/golangci-lint@latest"); err != nil {
		return err
	}

	cmd := exec.Command("golangci-lint", "run", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Dev() error {
	if err := Build(); err != nil {
		return err
	}

	docker := Docker{}
	return docker.Up()
}

type Docker mg.Namespace

func (d Docker) Up() error {
	fmt.Println("Starting Docker Compose services...")
	cmd := exec.Command("docker", "compose", "up", "-d")
	cmd.Dir = "compose"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (d Docker) Down() error {
	fmt.Println("Stopping Docker Compose services...")
	cmd := exec.Command("docker", "compose", "down")
	cmd.Dir = "compose"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (d Docker) Logs() error {
	fmt.Println("Showing Docker Compose logs...")
	cmd := exec.Command("docker", "compose", "logs", "-f")
	cmd.Dir = "compose"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (d Docker) Restart() error {
	fmt.Println("Restarting Docker Compose services...")
	if err := d.Down(); err != nil {
		return err
	}
	return d.Up()
}

func Clean() error {
	fmt.Println("Cleaning...")
	listDeletedFiles := []string{"dist", "mcp-server"}
	for _, file := range listDeletedFiles {
		if err := os.RemoveAll(file); err != nil {
			return fmt.Errorf("failed to remove %s directory: %w", file, err)
		}
	}
	return nil
}

func ensureToolInstalled(toolName, installPackage string) error {
	if _, err := exec.LookPath(toolName); err != nil {
		fmt.Printf("Installing %s...\n", toolName)
		installCmd := exec.Command("go", "install", installPackage)
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr
		if err := installCmd.Run(); err != nil {
			return fmt.Errorf("failed to install %s: %w", toolName, err)
		}
	}
	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
