package cmd

import (
	"bufio"
	"embed"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

type ProjectData struct {
	ProjectName   string
	ProjectRoot   string
	RepositoryURL string
	Author        string
}

var projectData = ProjectData{
	ProjectName:   "",
	ProjectRoot:   "",
	RepositoryURL: "",
	Author:        "",
}

//go:embed templates/*
var templateFS embed.FS

func getRandomBrewMessage() string {
	messages := []string{
		"🧪 This one is brewing a fresh batch of Skooma...",
		"🦁 Khajiit has wares, if you have coin...",
		"🌙 By Azura! This one crafts magical elixir...",
		"🏝️ May your roads lead you to warm sands...",
		"🧙 This one mixes moon sugar and nightshade...",
		"🏺 Psst! Khajiit knows you come for the good stuff...",
	}
	return messages[rand.Intn(len(messages))]
}

var brewCmd = &cobra.Command{
	Use:   "brew <project_name>",
	Short: "Brew a new project",
	Long:  `Brew a new project with the given name. This command will create a new directory with the project name and generate the necessary files for a basic project structure.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(getRandomBrewMessage())

		// Prompt for project name if not provided by argument
		if len(args) == 1 {
			projectData.ProjectName = sanitizeProjectName(args[0])
		}
		if projectData.ProjectName == "" {
			for {
				err := promptProjectName()
				if err != nil {
					fmt.Printf("❌ %v\n", err)
					continue
				}
				break
			}
		}

		// Get current working directory
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("❌ Failed to get current working directory: %v\n", err)
			os.Exit(1)
		}
		projectData.ProjectRoot = filepath.Join(cwd, projectData.ProjectName)

		// Early return if project directory already exists
		if _, err := os.Stat(projectData.ProjectRoot); !os.IsNotExist(err) {
			fmt.Printf("❌ Directory '%s' already exists\n", projectData.ProjectRoot)
			os.Exit(1)
		}

		// Prompt for repository URL if not provided by flag
		if projectData.RepositoryURL != "" {
			projectData.RepositoryURL = sanitizeRepositoryURL(projectData.RepositoryURL)
		}
		if projectData.RepositoryURL == "" {
			err := promptRepositoryURL()
			if err != nil {
				fmt.Printf("❌ Invalid repository URL: %v\n", err)
			}
		}

		// Use default repository URL if not provided by user
		if projectData.RepositoryURL == "" {
			projectData.RepositoryURL = fmt.Sprintf("github.com/username/%s", projectData.ProjectName)
		}

		// Prompt for author name if not provided by flag
		if projectData.Author != "" {
			projectData.Author = sanitizeAuthorName(projectData.Author)
		}
		if projectData.Author == "" {
			err := promptAuthorName()
			if err != nil {
				fmt.Printf("❌ Invalid author name: %v\n", err)
			}
		}

		// Start brewing spinner
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = " Brewing..."
		s.Start()

		// Prepare project data for templating
		err = scaffoldProject()
		if err != nil {
			s.Stop()
			fmt.Printf("❌ Failed to brew project\n\n%v\n", err)
			os.Exit(1)
		}

		// Simulate scaffolding work
		time.Sleep(2 * time.Second)

		s.Stop()
		fmt.Printf("\n✅ '%s' has finished brewing!\n\n", projectData.ProjectName)

		// Print project details
		fmt.Printf("Directory: %s\n", projectData.ProjectRoot)
		fmt.Printf("Repository: https://%s\n", projectData.RepositoryURL)
		if projectData.Author != "" {
			fmt.Printf("Author: %s\n", projectData.Author)
		}
	},
}

func init() {
	rootCmd.AddCommand(brewCmd)
	brewCmd.Flags().StringVarP(&projectData.RepositoryURL, "repo", "r", "", "Repository URL (e.g., github.com/username/repo)")
	brewCmd.Flags().StringVarP(&projectData.Author, "author", "a", "", "Author name")
}

func sanitizeProjectName(name string) string {
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.TrimSpace(name)
	return name
}

func sanitizeRepositoryURL(url string) string {
	url = strings.TrimSpace(url)
	// Remove protocol if user accidentally included it
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "https://")
	return url
}

func sanitizeAuthorName(name string) string {
	return strings.TrimSpace(name)
}

func promptProjectName() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter project name: ")
	projectName, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}
	projectData.ProjectName = sanitizeProjectName(projectName)
	if projectData.ProjectName == "" {
		return fmt.Errorf("project name cannot be empty")
	}
	return nil
}

func promptRepositoryURL() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter repository URL (e.g., github.com/username/repo): ")
	repositoryURL, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}
	projectData.RepositoryURL = sanitizeRepositoryURL(repositoryURL)
	return nil
}

func promptAuthorName() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter author name (e.g., John Doe <john.doe@example.com>): ")
	author, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}
	projectData.Author = sanitizeAuthorName(author)
	return nil
}

func scaffoldProject() error {
	err := createProjectRoot()
	if err != nil {
		return fmt.Errorf("failed to brew project: %w", err)
	}

	// Scaffold backend
	err = createBackend()
	if err != nil {
		return fmt.Errorf("failed to brew project: %w", err)
	}

	// Scaffold frontend
	err = createFrontend()
	if err != nil {
		return fmt.Errorf("failed to brew project: %w", err)
	}
	return nil
}

func createProjectRoot() error {
	projectRoot := projectData.ProjectRoot

	// Create project directory
	err := os.Mkdir(projectRoot, 0755)
	if err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Process root-level templates
	err = processTemplate("templates/docker-compose.yml.tmpl", filepath.Join(projectRoot, "docker-compose.yml"))
	if err != nil {
		return fmt.Errorf("failed to process root-level templates: %w", err)
	}
	return nil
}

func createBackend() error {
	backendPath := filepath.Join(projectData.ProjectRoot, "backend")
	err := os.Mkdir(backendPath, 0755)
	if err != nil {
		return err
	}

	// Process backend templates
	templates := []struct {
		src, dst string
	}{
		{"templates/backend/go.mod.tmpl", filepath.Join(backendPath, "go.mod")},
		{"templates/backend/main.go.tmpl", filepath.Join(backendPath, "main.go")},
	}

	for _, tmpl := range templates {
		if err := processTemplate(tmpl.src, tmpl.dst); err != nil {
			return err
		}
	}
	return nil
}

func createFrontend() error {
	frontendPath := filepath.Join(projectData.ProjectRoot, "frontend")
	err := os.Mkdir(frontendPath, 0755)
	if err != nil {
		return err
	}

	// Process frontend templates
	templates := []struct {
		src, dst string
	}{
		{"templates/frontend/package.json.tmpl", filepath.Join(frontendPath, "package.json")},
	}

	for _, tmpl := range templates {
		if err := processTemplate(tmpl.src, tmpl.dst); err != nil {
			return err
		}
	}
	return nil
}

func processTemplate(templatePath, outputPath string) error {
	// Read template from embedded filesystem
	content, err := templateFS.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	// Parse and execute template
	tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	// Create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", outputPath, err)
	}
	defer outputFile.Close()

	// Execute template with data
	return tmpl.Execute(outputFile, projectData)
}
