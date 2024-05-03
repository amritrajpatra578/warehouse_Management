package main

import (
	"fmt"
	"log"
	"time"
)

type ProductService interface {
	Create(Product) error
	Update(Product) error
	GetById(id int) (Product, error)
	GetAll() ([]Product, error)
	Delete(id int) error
	subscribe(Subscriber) error
	unsubscribe(Subscriber) error
	notify()
}

type ProductServiceImpl struct {
	repo        Repo
	subscribers []Subscriber
}

func NewProductServiceImpl(repo Repo) *ProductServiceImpl {
	return &ProductServiceImpl{
		repo: repo,
	}
}

func (s *ProductServiceImpl) Create(product Product) error {

	if err := validateProduct(product); err != nil {
		return fmt.Errorf("create product: %w", err)
	}

	timeNow := time.Now()
	product.CreatedAt = timeNow
	product.UpdatedAt = timeNow

	if err := s.repo.Create(product); err != nil {
		return err
	}
	s.notify()
	return nil
}

func (s *ProductServiceImpl) Update(product Product) error {
	if err := validateProduct(product); err != nil {
		return fmt.Errorf("update product: %w", err)
	}

	product.UpdatedAt = time.Now()

	if err := s.repo.Update(product); err != nil {
		return err
	}
	s.notify()
	return nil

}

func (s *ProductServiceImpl) GetAll() ([]Product, error) {
	return s.repo.GetAll()
}

func (s *ProductServiceImpl) GetById(id int) (Product, error) {
	return s.repo.GetById(id)
}

func (s *ProductServiceImpl) Delete(id int) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	s.notify()
	return nil
}

func (s *ProductServiceImpl) subscribe(subscriber Subscriber) error {
	if subscriber.Id() == "" {
		return errEmptyId
	}
	s.subscribers = append(s.subscribers, subscriber)
	return nil
}

func (s *ProductServiceImpl) unsubscribe(subscriber Subscriber) error {
	if subscriber.Id() == "" {
		return errEmptyId
	}

	subscriberId := subscriber.Id()

	for index, sub := range s.subscribers {
		if subscriberId == sub.Id() {
			s.subscribers = append(s.subscribers[:index], s.subscribers[index+1:]...)
			return nil
		}
	}

	return fmt.Errorf("subscriber with %s id not found", subscriberId)
}

func (s *ProductServiceImpl) notify() {
	products, err := s.GetAll()
	if err != nil {
		log.Println("error while getting list of products:", err)
		return
	}
	for _, sub := range s.subscribers {
		sub.Update(products)
	}
}
