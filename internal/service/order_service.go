package service

import (
	"context"

	"github.com/sparxfort1ano/order-service/internal/cache"
	"github.com/sparxfort1ano/order-service/internal/repository"
)

// Интерфейс для репозитория Postgres
type Repo interface {
	Init(context.Context) error
	Upsert(context.Context, repository.Order) error
	Get(context.Context, string) (repository.Order, error)
	All(context.Context, int) ([]repository.Order, error)
}

// Основной сервис с б/д и кэшем
type Service struct {
	repo  Repo
	cache *cache.OrderCache
}

// Создание нового сервиса
func New(repo Repo, c *cache.OrderCache) *Service { return &Service{repo: repo, cache: c} }

// Прогреваем кэш: загружаем последние n заказов из б/д
func (s *Service) Warm(ctx context.Context, n int) error {
	list, err := s.repo.All(ctx, n)
	if err != nil {
		return err
	}
	s.cache.Warm(list)
	return nil
}

// Получаем заказ по id: сначала ищем в кэше, если нет -- берем из б/д
func (s *Service) Get(ctx context.Context, id string) (repository.Order, error) {
	if o, ok := s.cache.Get(id); ok {
		return o, nil
	}
	o, err := s.repo.Get(ctx, id)
	if err != nil {
		return repository.Order{}, err
	}
	s.cache.Set(o)
	return o, nil
}

// Сохраняем заказ: пишем в б/д и кладем в кэш
func (s *Service) Save(ctx context.Context, o repository.Order) error {
	if err := s.repo.Upsert(ctx, o); err != nil {
		return err
	}
	s.cache.Set(o)
	return nil
}
