package cli

import (
	"fmt"
	"os"
	"temp/config"
	"temp/global"
	"temp/repositories"
)

// ParseAndExecute parses CLI arguments and executes commands. Returns (handled, error)
func ParseAndExecute(cfg *config.Config) (bool, error) {
	if len(os.Args) < 2 {
		return false, nil
	}
	switch os.Args[1] {
	case "migrate":
		userRepo := repositories.NewUserRepo()
		if err := userRepo.Migrate(); err != nil {
			return true, fmt.Errorf("migration failed: %w", err)
		}
		global.Logger.Info("Migration completed successfully.")
		return true, nil
	default:
		fmt.Println("Unknown command:", os.Args[1])
		return true, nil
	}
}
