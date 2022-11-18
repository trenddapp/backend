package workflow

import (
	"context"
	"errors"
	"strings"

	"github.com/trenddapp/backend/service/wordle/model"
	"github.com/trenddapp/backend/service/wordle/repository/word"
	"github.com/trenddapp/backend/service/wordle/repository/wordle"
)

type Engine struct {
	wordRepository   word.Repository
	wordleRepository wordle.Repository
}

func NewEngine(
	wordRepository word.Repository,
	wordleRepository wordle.Repository,
) *Engine {
	return &Engine{
		wordRepository:   wordRepository,
		wordleRepository: wordleRepository,
	}
}

func (e Engine) GetWordle(ctx context.Context, id string, opts ...wordle.Option) (model.Wordle, error) {
	out, err := e.wordleRepository.GetWordle(ctx, id, opts...)
	if err != nil {
		return model.Wordle{}, err
	}

	return out, nil
}

func (e Engine) CreateWordle(ctx context.Context, in model.Wordle) (model.Wordle, error) {
	word, err := e.wordRepository.GetRandomWord(ctx, "en_US")
	if err != nil {
		return model.Wordle{}, err
	}

	in.Solution = word.Value

	out, err := e.wordleRepository.CreateWordle(ctx, in)
	if err != nil {
		return model.Wordle{}, err
	}

	return out, nil
}

func (e Engine) UpdateWordle(ctx context.Context, in model.Wordle, opts ...wordle.Option) (model.Wordle, error) {
	out, err := e.wordleRepository.GetWordle(ctx, in.ID, wordle.WithUserID(in.UserID))
	if err != nil {
		return model.Wordle{}, err
	}

	if out.Status == model.StatusComplete {
		return model.Wordle{}, errors.New("cannot update complete wordle")
	}

	if len(in.Guesses) > 6 || len(in.Guesses) < len(out.Guesses) {
		return model.Wordle{}, errors.New("invalid guesses length")
	}

	for i, guess := range out.Guesses {
		if in.Guesses[i] == guess {
			continue
		}

		return model.Wordle{}, errors.New("cannot change old guesses")
	}

	if len(in.Guesses) == 6 {
		in.Status = model.StatusComplete
	}

	for _, guess := range in.Guesses {
		word := model.Word{
			Locale: "en_US",
			Value:  strings.ToLower(guess),
		}

		isValidWord, err := e.wordRepository.IsValidWord(ctx, word)
		if err != nil {
			return model.Wordle{}, err
		}

		if !isValidWord {
			return model.Wordle{}, errors.New("invalid word")
		}

		if guess == out.Solution {
			in.Status = model.StatusComplete
		}
	}

	out, err = e.wordleRepository.UpdateWordle(ctx, in, opts...)
	if err != nil {
		return model.Wordle{}, err
	}

	return out, nil
}

func (e Engine) DeleteWordle(ctx context.Context, id string, opts ...wordle.Option) (model.Wordle, error) {
	out, err := e.wordleRepository.DeleteWordle(ctx, id, opts...)
	if err != nil {
		return model.Wordle{}, err
	}

	return out, err
}

func (e Engine) ListWordles(ctx context.Context, opts ...wordle.Option) ([]model.Wordle, string, error) {
	wordles, nextPageToken, err := e.wordleRepository.ListWordles(ctx, opts...)
	if err != nil {
		return nil, "", err
	}

	return wordles, nextPageToken, nil
}
