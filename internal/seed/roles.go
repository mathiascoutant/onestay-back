package seed

import (
	"context"
	"log"

	"onestay-back/internal/models"
	"onestay-back/internal/repository"
)

func SeedRoles() error {
	ctx := context.Background()
	roleRepo := repository.NewRoleRepository()

	roles := []struct {
		name string
		slug string
	}{
		{"Client", models.RoleClient},
		{"Loueur", models.RoleLoueur},
		{"Admin", models.RoleAdmin},
		{"Super Admin", models.RoleSuperAdmin},
	}

	for _, roleData := range roles {
		exists, err := roleRepo.ExistsBySlug(ctx, roleData.slug)
		if err != nil {
			return err
		}

		if !exists {
			role := &models.Role{
				Name: roleData.name,
				Slug: roleData.slug,
			}

			if err := roleRepo.Create(ctx, role); err != nil {
				return err
			}
			log.Printf("Rôle créé: %s (%s)", roleData.name, roleData.slug)
		} else {
			log.Printf("Rôle déjà existant: %s (%s)", roleData.name, roleData.slug)
		}
	}

	return nil
}
