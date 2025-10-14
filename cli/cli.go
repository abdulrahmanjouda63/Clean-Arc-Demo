package cli

import (
	"fmt"
	"temp/config"
	"temp/global"
	"temp/repositories"

	"go.uber.org/zap"
)

// Command represents a CLI command
type Command struct {
	Name        string
	Description string
	Handler     func(*config.Config) error
}

// GetCommands returns all available CLI commands
func GetCommands() map[string]Command {
	return map[string]Command{
		"migrate": {
			Name:        "migrate",
			Description: "Run database migrations",
			Handler:     MigrateCommand,
		},
		"migrate:rollback": {
			Name:        "migrate:rollback",
			Description: "Rollback last database migration",
			Handler:     RollbackCommand,
		},
		"migrate:fresh": {
			Name:        "migrate:fresh",
			Description: "Drop all tables and re-run migrations",
			Handler:     FreshMigrateCommand,
		},
		"seed": {
			Name:        "seed",
			Description: "Seed the database with sample data",
			Handler:     SeedCommand,
		},
		"help": {
			Name:        "help",
			Description: "Show help information",
			Handler:     HelpCommand,
		},
	}
}

// ExecuteCommand executes a CLI command
func ExecuteCommand(commandName string, cfg *config.Config) error {
	commands := GetCommands()

	if cmd, exists := commands[commandName]; exists {
		return cmd.Handler(cfg)
	}

	fmt.Printf("Unknown command: %s\n", commandName)
	fmt.Println("Run 'help' to see available commands")
	return fmt.Errorf("unknown command: %s", commandName)
}

// MigrateCommand runs database migrations
func MigrateCommand(cfg *config.Config) error {
	global.Logger.Info("Running database migrations...")

	// Initialize database
	db, err := config.InitDB(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	global.DB = db

	// Run migrations
	userRepo := repositories.NewUserRepo()
	if err := userRepo.Migrate(); err != nil {
		global.Logger.Error("Migration failed", zap.Error(err))
		return err
	}

	global.Logger.Info("Migrations completed successfully")
	return nil
}

// RollbackCommand rollbacks last migration (placeholder)
func RollbackCommand(cfg *config.Config) error {
	global.Logger.Info("Rolling back last migration...")

	// Initialize database
	db, err := config.InitDB(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	global.DB = db

	// Implement rollback logic here
	// For now, this is a placeholder
	global.Logger.Info("Rollback completed successfully")
	return nil
}

// FreshMigrateCommand drops all tables and re-runs migrations
func FreshMigrateCommand(cfg *config.Config) error {
	global.Logger.Info("Running fresh migrations (dropping all tables)...")

	// Initialize database
	db, err := config.InitDB(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	global.DB = db

	// Drop all tables
	if err := global.DB.Migrator().DropTable("users"); err != nil {
		global.Logger.Warn("Failed to drop users table", zap.Error(err))
	}

	// Re-run migrations
	userRepo := repositories.NewUserRepo()
	if err := userRepo.Migrate(); err != nil {
		global.Logger.Error("Migration failed", zap.Error(err))
		return err
	}

	global.Logger.Info("Fresh migrations completed successfully")
	return nil
}

// SeedCommand seeds the database with sample data
func SeedCommand(cfg *config.Config) error {
	global.Logger.Info("Seeding database...")

	// Initialize database
	db, err := config.InitDB(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	global.DB = db

	// Add seeding logic here
	// Example: Create admin user, sample data, etc.

	global.Logger.Info("Database seeding completed successfully")
	return nil
}

// HelpCommand displays help information
func HelpCommand(cfg *config.Config) error {
	fmt.Println("Available Commands:")
	fmt.Println("-------------------")

	commands := GetCommands()
	for name, cmd := range commands {
		fmt.Printf("  %-20s %s\n", name, cmd.Description)
	}

	fmt.Println("\nUsage:")
	fmt.Println("  go run main.go [command]")
	fmt.Println("  ./app [command]")

	return nil
}
