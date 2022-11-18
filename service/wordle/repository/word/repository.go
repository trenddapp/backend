package word

import (
	"context"
	"embed"
	"encoding/json"
	"math/rand"
	"strings"

	"github.com/trenddapp/backend/service/wordle/model"
)

//go:embed data/*
var fsys embed.FS

type repository struct {
	words map[string][]string
}

func NewRepository() Repository {
	return &repository{
		words: map[string][]string{
			"en_US": readWordsByLocale("en_US"),
		},
	}
}

func (r repository) GetRandomWord(ctx context.Context, locale string) (model.Word, error) {
	return model.Word{
		Locale: locale,
		Value:  r.words[locale][rand.Intn(len(r.words[locale]))],
	}, nil
}

func (r repository) IsValidWord(ctx context.Context, in model.Word) (bool, error) {
	in.Value = strings.ToLower(in.Value)

	for _, word := range r.words[in.Locale] {
		if word == in.Value {
			return true, nil
		}
	}

	return false, nil
}

func readWordsByLocale(locale string) []string {
	var data struct {
		Words []string `json:"words"`
	}

	content, err := fsys.ReadFile("data/" + strings.ToLower(locale) + ".json")
	if err != nil {
		return nil
	}

	if err := json.Unmarshal(content, &data); err != nil {
		return nil
	}

	return data.Words
}
