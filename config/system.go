package config

import (
	"evedict/internal/promptgen"
	"fmt"
)

type SystemConfig struct {
	PromptTemplate *promptgen.PromptTemplate
}

func ConfigureSystem() (*SystemConfig, error) {
	promptTemplate, err := promptgen.LoadPromptTemplate("config/prompts.yaml")
	if err != nil {
		return nil, fmt.Errorf("error loading prompt template: %v", err)
	}

	return &SystemConfig{PromptTemplate: promptTemplate}, nil
}
