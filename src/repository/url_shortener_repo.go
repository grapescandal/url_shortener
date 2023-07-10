package repository

import (
	"context"
	"url_shortener/src/shared/model"
)

type URLShortenerRepository interface {
	CreateShortener(ctx context.Context, inputData model.ShortURL) error
	GetFullURL(ctx context.Context, code string) *model.ShortURL
	IncreaseAccessCount(ctx context.Context, code string) error
}

type URLShortenerRepositoryImpl struct {
	shortURLs []model.ShortURL
}

func NewURLShortenerRepositoryImpl() *URLShortenerRepositoryImpl {
	return &URLShortenerRepositoryImpl{}
}

func (repo *URLShortenerRepositoryImpl) CreateShortener(ctx context.Context, inputData model.ShortURL) error {
	repo.shortURLs = append(repo.shortURLs, inputData)
	return nil
}

func (repo *URLShortenerRepositoryImpl) GetFullURL(ctx context.Context, code string) *model.ShortURL {
	for i := 0; i < len(repo.shortURLs); i++ {
		if repo.shortURLs[i].ShortCode == code {
			return &repo.shortURLs[i]
		}
	}
	return nil
}

func (repo *URLShortenerRepositoryImpl) IncreaseAccessCount(ctx context.Context, code string) error {
	for i := 0; i < len(repo.shortURLs); i++ {
		if repo.shortURLs[i].ShortCode == code {
			repo.shortURLs[i].AccessCount++
			break
		}
	}

	return nil
}
