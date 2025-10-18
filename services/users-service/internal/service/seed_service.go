package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/youngermaster/bookstore/services/users-service/internal/domain"
	"gorm.io/gorm"
)

// SeedRoles creates default roles if they don't exist
func SeedRoles(db *gorm.DB) error {
	ctx := context.Background()

	// Define default roles
	customerPermissions, _ := json.Marshal([]string{
		"books:read",
		"wishlist:read",
		"wishlist:write",
		"orders:read",
		"orders:write",
		"profile:read",
		"profile:write",
	})

	adminPermissions, _ := json.Marshal([]string{
		"books:read",
		"books:write",
		"books:delete",
		"users:read",
		"users:write",
		"orders:read",
		"orders:write",
		"orders:manage",
		"logs:read",
	})

	roles := []domain.Role{
		{
			Name:        "customer",
			Permissions: string(customerPermissions),
		},
		{
			Name:        "admin",
			Permissions: string(adminPermissions),
		},
	}

	for _, role := range roles {
		var existing domain.Role
		err := db.WithContext(ctx).Where("name = ?", role.Name).First(&existing).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Role doesn't exist, create it
				if err := db.WithContext(ctx).Create(&role).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}
