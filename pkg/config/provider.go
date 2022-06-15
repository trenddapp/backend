package config

import (
	"os"
	"strings"

	"go.uber.org/config"
)

func NewYAML() (*config.YAML, error) {
	options := []config.YAMLOption{
		config.Expand(os.LookupEnv),
	}

	for _, path := range generatePaths("base") {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		options = append(options, config.File(path))
	}

	options = append(
		options,
		config.Source(strings.NewReader(os.Getenv("CONFIG"))),
	)

	return config.NewYAML(options...)
}

func generatePaths(filenames ...string) []string {
	paths := []string{}

	for _, filename := range filenames {
		paths = append(
			paths,
			"/config/"+filename+".yaml",
			"./config/"+filename+".yaml",
		)
	}

	return paths
}
