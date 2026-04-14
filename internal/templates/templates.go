package templates

import (
	"github.com/mark-rodgers/skooma/internal/config"
	"github.com/mark-rodgers/skooma/internal/types"
)

func GetTemplates() (map[string]types.Template, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	return cfg.Templates, nil
}

func GetTemplateByName(name string) (*types.Template, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	tmpl, exists := cfg.Templates[name]
	if !exists {
		return nil, nil
	}

	return &tmpl, nil
}

func AddTemplate(name string, template types.Template) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	cfg.Templates[name] = template
	return config.SaveConfig(cfg)
}

func RemoveTemplate(name string) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	delete(cfg.Templates, name)
	return config.SaveConfig(cfg)
}
