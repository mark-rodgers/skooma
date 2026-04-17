package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/tabwriter"
	"time"

	"charm.land/huh/v2"
	"github.com/briandowns/spinner"
	"github.com/skooma-cli/skooma/internal/sanitize"
	"github.com/skooma-cli/skooma/internal/templates"
	"github.com/skooma-cli/skooma/internal/types"
	"github.com/skooma-cli/skooma/internal/utils"
	"github.com/skooma-cli/skooma/internal/validators"
	"github.com/spf13/cobra"
)

var brewProjectNameArg string
var brewTemplateFlag string
var brewRepoUrlFlag string
var brewAuthorFlag string

var brewCmd = &cobra.Command{
	Use:   "brew PROJECT_NAME",
	Short: "Brew a new project",
	Long: `Brew a new project with the given name.
This command will create a new directory with the project name and generate
the necessary files for a basic project structure.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n\n", utils.GetRandomKhajiitPhrase())

		if len(args) > 0 {
			brewProjectNameArg = args[0]
		}

		groups := []*huh.Group{}

		// Validators for the project name input
		projectNameValidators := []func(string) error{
			validators.NotEmpty("Project name"), // only meaningful in the TUI, redundant if a flag is provided 🤷
			validators.NoSpaces("Project name"),
			validators.NoUnderscores("Project name"),
		}
		// If no project name was provided, prompt the user; otherwise validate the provided value
		if brewProjectNameArg == "" {
			groups = append(groups, huh.NewGroup(
				huh.NewInput().
					Title("Project name:").
					Value(&brewProjectNameArg).
					Validate(validators.All(projectNameValidators...)),
			))
		} else {
			if err := validators.All(projectNameValidators...)(brewProjectNameArg); err != nil {
				log.Fatalf("❌ Invalid project name: %v\n", err)
			}
		}

		// Load templates to build options for the template selection prompt
		templates, err := templates.GetTemplates()
		if err != nil {
			log.Fatalf("❌ Error loading templates: %v\n", err)
		}

		// If no template was provided, prompt the user; otherwise validate the provided template name exists
		if brewTemplateFlag == "" {
			templateOptions := make([]huh.Option[string], 0, len(templates))
			for name, tmpl := range templates {
				templateOptions = append(templateOptions, huh.NewOption(name+" - "+tmpl.Description, name))
			}
			groups = append(groups, huh.NewGroup(
				huh.NewSelect[string]().
					Title("Template").
					Options(templateOptions...).
					Value(&brewTemplateFlag),
			))
		} else if _, ok := templates[brewTemplateFlag]; !ok {
			log.Fatalf("❌ Invalid template name: '%s'. Use 'skooma template ls' to see available templates.\n", brewTemplateFlag)
		}

		// Validators for the repository URL input
		repoUrlValidators := []func(string) error{
			validators.NoSpaces("Repository URL"),
			validators.ValidURL("Repository URL"),
		}
		// If no repository URL was provided, prompt the user; otherwise validate the provided value
		if brewRepoUrlFlag == "" {
			groups = append(groups, huh.NewGroup(
				huh.NewInput().
					Title("Repository URL (e.g., github.com/user/repo):").
					Value(&brewRepoUrlFlag).
					Validate(validators.AllowEmpty(repoUrlValidators...)),
			))
		} else {
			if err := validators.All(repoUrlValidators...)(brewRepoUrlFlag); err != nil {
				log.Fatalf("❌ Invalid repository URL: %v\n", err)
			}
		}

		// Validators for the author name input
		authorValidators := []func(string) error{
			validators.RFC5322Address("Author"),
		}
		// If no author was provided, prompt the user; otherwise validate the provided value
		if brewAuthorFlag == "" {
			groups = append(groups, huh.NewGroup(
				huh.NewInput().
					Title("Author name (e.g., Name <email@example.com>):").
					Value(&brewAuthorFlag).
					Validate(validators.AllowEmpty(authorValidators...)),
			))
		} else {
			if err := validators.All(authorValidators...)(brewAuthorFlag); err != nil {
				log.Fatalf("❌ Invalid author name: %v\n", err)
			}
		}

		form := huh.NewForm(groups...)

		// Run the form to collect user input
		err = form.Run()
		if err != nil {
			log.Fatalf("❌ Failed to run form: %v\n", err)
		}

		// Get current working directory to build absolute path for project root directory
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("❌ Failed to get current working directory: %v\n", err)
		}

		// Build project data struct to pass to the brewing process
		project := types.ProjectData{
			Name:     brewProjectNameArg,
			RootDir:  filepath.Join(cwd, brewProjectNameArg),
			Template: templates[brewTemplateFlag],
			RepoURL:  sanitize.StripHTTPPrefix(brewRepoUrlFlag),
			Author:   brewAuthorFlag,
		}

		// Check if project root directory already exists before starting the brewing process
		if _, err := os.Stat(project.RootDir); !os.IsNotExist(err) {
			log.Fatalf("❌ Directory '%s' already exists: %v\n", project.RootDir, err)
		}

		// Start brewing spinner
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = " Brewing..."
		s.Start()

		// TODO: refactor brewing logic to pull templates for git repos instead of embedded templates

		// Prepare project data for templating
		// err = scaffoldProject(project)
		// if err != nil {
		// 	log.Fatalf("❌ Failed to brew project\n\n%v\n", err)
		// }

		// Simulate scaffolding work
		time.Sleep(2 * time.Second)

		s.Stop()
		fmt.Printf("\n✅ '%s' has finished brewing!\n\n", project.Name)

		// Print project details
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "Template\t%s - %s\n", project.Template.Name, project.Template.Description)
		fmt.Fprintf(w, "Repository\t%s\n", project.RepoURL)
		fmt.Fprintf(w, "Author\t%s\n", project.Author)
		fmt.Fprintf(w, "Directory\t%s\n", project.RootDir)
		w.Flush()
	},
}

// init registers the brew command and its flags with the root command.
func init() {
	rootCmd.AddCommand(brewCmd)
	brewCmd.Flags().StringVarP(&brewTemplateFlag, "template", "t", "", "Template name")
	brewCmd.Flags().StringVarP(&brewRepoUrlFlag, "repo", "r", "", "Repository URL (e.g., github.com/user/repo)")
	brewCmd.Flags().StringVarP(&brewAuthorFlag, "author", "a", "", "Author name")
}

// TODO: all this shit should go inside ./internal and out of the cmd package

// // scaffoldProject creates the project directory structure and generates files based on templates.
// func scaffoldProject(project types.ProjectData) error {
// 	err := createProjectRoot(project)
// 	if err != nil {
// 		return fmt.Errorf("failed to brew project: %w", err)
// 	}

// 	err = createBackend(project)
// 	if err != nil {
// 		return fmt.Errorf("failed to brew project: %w", err)
// 	}

// 	err = createFrontend(project)
// 	if err != nil {
// 		return fmt.Errorf("failed to brew project: %w", err)
// 	}
// 	return nil
// }

// // createProjectRoot creates the root project directory and processes root-level templates.
// func createProjectRoot(project types.ProjectData) error {
// 	err := os.Mkdir(project.RootDir, 0755)
// 	if err != nil {
// 		return fmt.Errorf("failed to create project root directory: %w", err)
// 	}

// 	// TODO: move template processing logic into a separate function that can be reused for all templates
// 	// instead of duplicating the code for each template

// 	// Process root-level templates
// 	err = processTemplate("templates/docker-compose.yml.tmpl", filepath.Join(project.RootDir, "docker-compose.yml"))
// 	if err != nil {
// 		return fmt.Errorf("failed to process root-level templates: %w", err)
// 	}
// 	return nil
// }

// // createBackend creates the backend directory and generates files based on templates.
// func createBackend(project types.ProjectData) error {
// 	backendPath := filepath.Join(project.RootDir, "backend")
// 	err := os.Mkdir(backendPath, 0755)
// 	if err != nil {
// 		return fmt.Errorf("failed to create backend directory: %w", err)
// 	}

// 	// Process backend templates
// 	templates := []struct {
// 		src, dst string
// 	}{
// 		{"templates/backend/go.mod.tmpl", filepath.Join(backendPath, "go.mod")},
// 		{"templates/backend/main.go.tmpl", filepath.Join(backendPath, "main.go")},
// 		{"templates/backend/Makefile.tmpl", filepath.Join(backendPath, "Makefile")},
// 	}

// 	for _, tmpl := range templates {
// 		if err := processTemplate(tmpl.src, tmpl.dst); err != nil {
// 			return fmt.Errorf("failed to process template %s: %w", tmpl.src, err)
// 		}
// 	}
// 	return nil
// }

// // createFrontend creates the frontend directory, subdirectories, and generates files based on templates.
// func createFrontend(project types.ProjectData) error {
// 	frontendPath := filepath.Join(project.RootDir, "frontend")
// 	err := os.Mkdir(frontendPath, 0755)
// 	if err != nil {
// 		return fmt.Errorf("failed to create frontend directory: %w", err)
// 	}

// 	subdirs := []string{"src", "src/assets", "public"}
// 	for _, subdir := range subdirs {
// 		err := os.Mkdir(filepath.Join(frontendPath, subdir), 0755)
// 		if err != nil {
// 			return fmt.Errorf("failed to create frontend subdirectory %s: %w", subdir, err)
// 		}
// 	}

// 	// Copy static asset files that don't require templating
// 	staticFiles := []struct {
// 		src, dst string
// 	}{
// 		// Public directory static files
// 		{"templates/frontend/public/favicon.svg", filepath.Join(frontendPath, "public", "favicon.svg")},
// 		{"templates/frontend/public/khajiit.webp", filepath.Join(frontendPath, "public", "khajiit.webp")},
// 	}
// 	for _, file := range staticFiles {
// 		if err := copyFile(file.src, file.dst); err != nil {
// 			return fmt.Errorf("failed to copy static file %s: %w", file.src, err)
// 		}
// 	}

// 	// Process frontend templates
// 	templates := []struct {
// 		src, dst string
// 	}{
// 		{"templates/frontend/gitignore.tmpl", filepath.Join(frontendPath, ".gitignore")},
// 		{"templates/frontend/eslint.config.js.tmpl", filepath.Join(frontendPath, "eslint.config.js")},
// 		{"templates/frontend/index.html.tmpl", filepath.Join(frontendPath, "index.html")},
// 		{"templates/frontend/package.json.tmpl", filepath.Join(frontendPath, "package.json")},
// 		{"templates/frontend/README.md.tmpl", filepath.Join(frontendPath, "README.md")},
// 		{"templates/frontend/tsconfig.json.tmpl", filepath.Join(frontendPath, "tsconfig.json")},
// 		{"templates/frontend/tsconfig.app.json.tmpl", filepath.Join(frontendPath, "tsconfig.app.json")},
// 		{"templates/frontend/tsconfig.node.json.tmpl", filepath.Join(frontendPath, "tsconfig.node.json")},
// 		{"templates/frontend/vite.config.ts.tmpl", filepath.Join(frontendPath, "vite.config.ts")},
// 		{"templates/frontend/src/App.css.tmpl", filepath.Join(frontendPath, "src", "App.css")},
// 		{"templates/frontend/src/App.tsx.tmpl", filepath.Join(frontendPath, "src", "App.tsx")},
// 		{"templates/frontend/src/index.css.tmpl", filepath.Join(frontendPath, "src", "index.css")},
// 		{"templates/frontend/src/main.tsx.tmpl", filepath.Join(frontendPath, "src", "main.tsx")},
// 	}

// 	for _, tmpl := range templates {
// 		if err := processTemplate(tmpl.src, tmpl.dst); err != nil {
// 			return fmt.Errorf("failed to process template %s: %w", tmpl.src, err)
// 		}
// 	}
// 	return nil
// }

// // copyFile reads a file from the embedded filesystem and writes it to the specified destination path.
// func copyFile(src, dst string) error {
// 	// Read file content from embedded filesystem
// 	content, err := templateFS.ReadFile(src)
// 	if err != nil {
// 		return fmt.Errorf("failed to read file %s: %w", src, err)
// 	}

// 	// Write content to destination path
// 	err = os.WriteFile(dst, content, 0644)
// 	if err != nil {
// 		return fmt.Errorf("failed to write file %s: %w", dst, err)
// 	}
// 	return nil
// }

// // processTemplate reads a template from the embedded filesystem, executes it with the project data, and writes the output to the specified path.
// func processTemplate(templatePath, outputPath string) error {
// 	// Read template from embedded filesystem
// 	content, err := templateFS.ReadFile(templatePath)
// 	if err != nil {
// 		return fmt.Errorf("failed to read template %s: %w", templatePath, err)
// 	}

// 	// Parse and execute template
// 	tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(content))
// 	if err != nil {
// 		return fmt.Errorf("failed to parse template %s: %w", templatePath, err)
// 	}

// 	// Create output file
// 	outputFile, err := os.Create(outputPath)
// 	if err != nil {
// 		return fmt.Errorf("failed to create output file %s: %w", outputPath, err)
// 	}
// 	defer outputFile.Close()

// 	// Execute template with data
// 	if err := tmpl.Execute(outputFile, project); err != nil {
// 		return fmt.Errorf("failed to execute template %s: %w", templatePath, err)
// 	}
// 	return nil
// }
