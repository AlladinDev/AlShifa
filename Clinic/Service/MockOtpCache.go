package service

import (
	"AlShifa/Clinic/models"
	"context"
	"time"
)

// Cache module for storing otp
type Cache struct {
}

func (c *Cache) Set(ctx context.Context, key string, value models.AddDoctorOtpPayload, ttl time.Duration) error {
	return nil
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	return nil
}

func (c *Cache) Get(ctx context.Context, key string) (models.AddDoctorOtpPayload, bool, error) {
	return models.AddDoctorOtpPayload{}, true, nil
}

func (c *Cache) Update(ctx context.Context, key string, value models.AddDoctorOtpPayload, ttl time.Duration) error {
	return nil
}
func ReturnNewCache() *Cache {
	return &Cache{}
}
