package main

import (
	"be-lecsens/asset_management/data-layer/config"
	"be-lecsens/asset_management/data-layer/repository"
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// parseUUID parses a string into a UUID
func parseUUID(uuidStr string) (uuid.UUID, error) {
	return uuid.Parse(uuidStr)
}

// CleanupDuplicatesCommand provides functionality to clean up duplicate asset documents
type CleanupDuplicatesCommand struct {
	assetDocumentRepo repository.AssetDocumentRepository
	db                *sql.DB
}

// NewCleanupDuplicatesCommand creates a new cleanup command
func NewCleanupDuplicatesCommand() (*CleanupDuplicatesCommand, error) {
	// Load configuration
	cfg := config.Load()

	// Build database connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)

	// Initialize database connection
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Initialize repository
	assetDocumentRepo := repository.NewAssetDocumentRepository(db)

	return &CleanupDuplicatesCommand{
		assetDocumentRepo: assetDocumentRepo,
		db:                db,
	}, nil
}

// Close closes the database connection
func (c *CleanupDuplicatesCommand) Close() error {
	return c.db.Close()
}

// GetDuplicatesReport generates a report of all duplicate documents
func (c *CleanupDuplicatesCommand) GetDuplicatesReport(ctx context.Context) error {
	log.Println("🔍 Scanning for duplicate asset documents...")

	// Get all duplicate documents
	duplicates, err := c.assetDocumentRepo.GetDuplicateDocuments(ctx)
	if err != nil {
		return fmt.Errorf("failed to get duplicate documents: %w", err)
	}

	if len(duplicates) == 0 {
		log.Println("✅ No duplicate documents found!")
		return nil
	}

	totalDuplicates := 0
	log.Printf("📊 Found duplicates in %d assets:\n", len(duplicates))
	log.Println("=" + fmt.Sprintf("%80s", "="))

	for assetID, docs := range duplicates {
		if len(docs) > 1 {
			duplicateCount := len(docs) - 1 // Subtract 1 because we keep the latest
			totalDuplicates += duplicateCount

			log.Printf("Asset ID: %s", assetID)
			log.Printf("  📄 Total documents: %d", len(docs))
			log.Printf("  🗑️  Duplicates to remove: %d", duplicateCount)
			log.Printf("  🏷️  Document types: ")

			for i, doc := range docs {
				status := "KEEP"
				if i > 0 { // First document (newest) is kept
					status = "DELETE"
				}
				log.Printf("    - %s [%s] (%s) - %s",
					doc.OriginalFilename,
					doc.DocumentType,
					doc.CreatedAt.Format("2006-01-02 15:04:05"),
					status)
			}
			log.Println()
		}
	}

	log.Println("=" + fmt.Sprintf("%80s", "="))
	log.Printf("📈 SUMMARY:")
	log.Printf("  🏢 Assets with duplicates: %d", len(duplicates))
	log.Printf("  📄 Total duplicate documents to remove: %d", totalDuplicates)
	log.Println()

	return nil
}

// CleanupAllDuplicates removes all duplicate documents from the database
func (c *CleanupDuplicatesCommand) CleanupAllDuplicates(ctx context.Context, dryRun bool) error {
	if dryRun {
		log.Println("🧪 DRY RUN MODE - No actual deletions will be performed")
		return c.GetDuplicatesReport(ctx)
	}

	log.Println("🧹 Starting cleanup of all duplicate asset documents...")

	// First show the report
	if err := c.GetDuplicatesReport(ctx); err != nil {
		return err
	}

	// Ask for confirmation
	fmt.Print("❓ Do you want to proceed with deletion? (yes/no): ")
	var response string
	fmt.Scanln(&response)

	if response != "yes" && response != "y" {
		log.Println("❌ Cleanup cancelled by user")
		return nil
	}

	// Perform actual cleanup
	deletedCount, err := c.assetDocumentRepo.CleanupAllDuplicateDocuments(ctx)
	if err != nil {
		return fmt.Errorf("failed to cleanup duplicate documents: %w", err)
	}

	log.Printf("✅ Successfully cleaned up %d duplicate documents!", deletedCount)
	return nil
}

// CleanupAssetDuplicates removes duplicate documents for a specific asset
func (c *CleanupDuplicatesCommand) CleanupAssetDuplicates(ctx context.Context, assetIDStr string, dryRun bool) error {
	// Parse asset ID
	assetID, err := parseUUID(assetIDStr)
	if err != nil {
		return fmt.Errorf("invalid asset ID format: %w", err)
	}

	if dryRun {
		log.Printf("🧪 DRY RUN MODE - Showing duplicates for asset %s", assetID)
	} else {
		log.Printf("🧹 Starting cleanup of duplicate documents for asset %s...", assetID)
	}

	// Get duplicates for this asset
	duplicates, err := c.assetDocumentRepo.GetDuplicateDocuments(ctx)
	if err != nil {
		return fmt.Errorf("failed to get duplicate documents: %w", err)
	}

	docs, exists := duplicates[assetID]
	if !exists || len(docs) <= 1 {
		log.Printf("✅ No duplicate documents found for asset %s", assetID)
		return nil
	}

	// Show what will be cleaned
	duplicateCount := len(docs) - 1
	log.Printf("📄 Found %d documents, %d duplicates to remove:", len(docs), duplicateCount)

	for i, doc := range docs {
		status := "KEEP"
		if i > 0 {
			status = "DELETE"
		}
		log.Printf("  - %s [%s] (%s) - %s",
			doc.OriginalFilename,
			doc.DocumentType,
			doc.CreatedAt.Format("2006-01-02 15:04:05"),
			status)
	}

	if dryRun {
		log.Printf("🧪 DRY RUN: Would delete %d duplicate documents", duplicateCount)
		return nil
	}

	// Ask for confirmation
	fmt.Printf("❓ Do you want to delete %d duplicate documents for asset %s? (yes/no): ", duplicateCount, assetID)
	var response string
	fmt.Scanln(&response)

	if response != "yes" && response != "y" {
		log.Println("❌ Cleanup cancelled by user")
		return nil
	}

	// Perform actual cleanup
	deletedCount, err := c.assetDocumentRepo.CleanupDuplicateDocuments(ctx, assetID)
	if err != nil {
		return fmt.Errorf("failed to cleanup duplicate documents for asset %s: %w", assetID, err)
	}

	log.Printf("✅ Successfully cleaned up %d duplicate documents for asset %s!", deletedCount, assetID)
	return nil
}

// RunCleanupCommand runs the cleanup command based on provided arguments
func (c *CleanupDuplicatesCommand) RunCleanupCommand(args []string) error {
	ctx := context.Background()

	if len(args) == 0 {
		// Show help
		c.showHelp()
		return nil
	}

	command := args[0]

	switch command {
	case "report", "scan":
		return c.GetDuplicatesReport(ctx)

	case "cleanup-all":
		dryRun := false
		if len(args) > 1 && args[1] == "--dry-run" {
			dryRun = true
		}
		return c.CleanupAllDuplicates(ctx, dryRun)

	case "cleanup-asset":
		if len(args) < 2 {
			return fmt.Errorf("asset ID is required for cleanup-asset command")
		}
		assetID := args[1]
		dryRun := false
		if len(args) > 2 && args[2] == "--dry-run" {
			dryRun = true
		}
		return c.CleanupAssetDuplicates(ctx, assetID, dryRun)

	default:
		c.showHelp()
		return fmt.Errorf("unknown command: %s", command)
	}
}

// showHelp displays help information
func (c *CleanupDuplicatesCommand) showHelp() {
	fmt.Println(
		`Asset Document Duplicate Cleanup Tool

Usage:
  go run cmd/cleanup-duplicates/main.go <command> [options]

Commands:
  report                    Show report of all duplicate documents
  scan                      Alias for 'report'
  cleanup-all [--dry-run]   Remove all duplicate documents (optionally dry run)
  cleanup-asset <asset-id> [--dry-run]  Remove duplicates for specific asset

Options:
  --dry-run                 Show what would be deleted without actually deleting

Examples:
  go run cmd/cleanup-duplicates/main.go report
  go run cmd/cleanup-duplicates/main.go cleanup-all --dry-run
  go run cmd/cleanup-duplicates/main.go cleanup-all
  go run cmd/cleanup-duplicates/main.go cleanup-asset 123e4567-e89b-12d3-a456-426614174000`)
}
