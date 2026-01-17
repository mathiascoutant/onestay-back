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
		id   string
		name string
		slug string
	}{
		{"1", "Client", models.RoleClient},
		{"2", "Loueur", models.RoleLoueur},
		{"3", "Admin", models.RoleAdmin},
		{"4", "Super Admin", models.RoleSuperAdmin},
	}

	for _, roleData := range roles {
		exists, err := roleRepo.ExistsBySlug(ctx, roleData.slug)
		if err != nil {
			return err
		}

		if !exists {
			role := &models.Role{
				ID:   roleData.id,
				Name: roleData.name,
				Slug: roleData.slug,
			}

			if err := roleRepo.Create(ctx, role); err != nil {
				return err
			}
			log.Printf("Rôle créé: %s (%s) avec ID: %s", roleData.name, roleData.slug, roleData.id)
		} else {
			log.Printf("Rôle déjà existant: %s (%s)", roleData.name, roleData.slug)
		}
	}

	return nil
}
