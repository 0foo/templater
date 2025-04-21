package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Config struct {
	Aliases map[string]string `json:"aliases"`
}

func getConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Unable to determine home directory: " + err.Error())
	}

	configDir := filepath.Join(homeDir, ".templater")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		panic("Unable to create config directory: " + err.Error())
	}

	return filepath.Join(configDir, "config.json")
}

func loadConfig() (*Config, error) {
	path := getConfigPath()
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return &Config{Aliases: make(map[string]string)}, nil
	}
	defer file.Close()
	var config Config
	err = json.NewDecoder(file).Decode(&config)
	return &config, err
}

func saveConfig(config *Config) error {
	path := getConfigPath()
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(config)
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()
		destFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer destFile.Close()
		_, err = io.Copy(destFile, srcFile)
		return err
	})
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mycli <save|build|list|delete> [alias]")
		return
	}

	command := os.Args[1]
	config, err := loadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	switch command {
	case "save":
		if len(os.Args) < 3 {
			fmt.Println("Usage: mycli save <alias>")
			return
		}
		alias := os.Args[2]
		cwd, _ := os.Getwd()
		config.Aliases[alias] = cwd
		if err := saveConfig(config); err != nil {
			fmt.Println("Failed to save config:", err)
		} else {
			fmt.Printf("Saved alias '%s' -> %s\n", alias, cwd)
		}

	case "build":
		if len(os.Args) < 3 {
			fmt.Println("Usage: mycli build <alias>")
			return
		}
		alias := os.Args[2]
		src, ok := config.Aliases[alias]
		if !ok {
			fmt.Println("Alias not found:", alias)
			return
		}

		cwd, _ := os.Getwd()
		destDir := filepath.Join(cwd, alias)

		if err := os.MkdirAll(destDir, 0755); err != nil {
			fmt.Println("Failed to create target directory:", err)
			return
		}

		if err := copyDir(src, destDir); err != nil {
			fmt.Println("Failed to copy:", err)
		} else {
			fmt.Printf("Copied from %s to %s\n", src, destDir)
		}

	case "dump":
		if len(os.Args) < 3 {
			fmt.Println("Usage: mycli dump <alias>")
			return
		}
		alias := os.Args[2]
		src, ok := config.Aliases[alias]
		if !ok {
			fmt.Println("Alias not found:", alias)
			return
		}

		dst, _ := os.Getwd()
		if err := copyDir(src, dst); err != nil {
			fmt.Println("Failed to copy:", err)
		} else {
			fmt.Printf("Dumped from %s to %s\n", src, dst)
		}

	case "list":
		if len(config.Aliases) == 0 {
			fmt.Println("No aliases found.")
		}
		for alias, path := range config.Aliases {
			fmt.Printf("%s -> %s\n", alias, path)
		}

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: mycli delete <alias>")
			return
		}
		alias := os.Args[2]
		if _, ok := config.Aliases[alias]; !ok {
			fmt.Println("Alias not found:", alias)
			return
		}
		delete(config.Aliases, alias)
		if err := saveConfig(config); err != nil {
			fmt.Println("Failed to update config:", err)
		} else {
			fmt.Printf("Deleted alias '%s'\n", alias)
		}

	default:
		fmt.Println("Unknown command:", command)
		fmt.Println("Usage: mycli <save|build|list|delete> [alias]")
	}
}
