package shortener

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"url_shortener/src/repository"
	"url_shortener/src/shared/model"
)

type Service interface {
	CreateShortener(ctx context.Context, inputData *model.ShortURL) (string, error)
	RedirectURL(ctx context.Context, code string) (string, error)
}

type ServiceImpl struct {
	URLShortenerRepo repository.URLShortenerRepository
}

func NewServiceImpl(urlShortenerRepo repository.URLShortenerRepository) *ServiceImpl {
	return &ServiceImpl{
		URLShortenerRepo: urlShortenerRepo,
	}
}

func (svc *ServiceImpl) CreateShortener(ctx context.Context, inputData *model.ShortURL) (string, error) {
	path := "CreateShortener"

	shortener := model.ShortURL{}

	shortener.ID = string(rand.Intn(100000000))
	shortener.ShortCode = generateShortCode()

	err := svc.URLShortenerRepo.CreateShortener(ctx, shortener)
	if err != nil {
		fmt.Printf("%s: err: %s", path, err.Error())
		return "Failed", err
	}
	return "Success", nil
}

func generateShortCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 6)

	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}

func (svc *ServiceImpl) RedirectURL(ctx context.Context, code string) (string, error) {
	path := "RedirectURL"

	url := svc.URLShortenerRepo.GetFullURL(ctx, code)

	if url == nil {
		err := errors.New("Shortener not found")
		fmt.Printf("%s: ", path, err)
		return "", err
	}

	svc.URLShortenerRepo.IncreaseAccessCount(ctx, code)

	return url.FullURL, nil
}
