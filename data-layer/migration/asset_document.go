package migration

import (
	"be-lecsens/asset_management/data-layer/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// CreateAssetDocumentTable creates the asset_documents table with proper foreign key constraints
func CreateAssetDocumentTable(cfg *config.Config) error {
	log.Println("Creating asset_documents table...")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// SQL untuk membuat tabel asset_documents
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS asset_documents (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		tenant_id UUID NULL,
		asset_id UUID NULL,
		document_type VARCHAR(100) NOT NULL,
		file_url VARCHAR(1000) NOT NULL,
		cloudinary_id VARCHAR(500) NOT NULL,
		original_filename VARCHAR(500) NOT NULL,
		file_size BIGINT NOT NULL DEFAULT 0,
		mime_type VARCHAR(100) NOT NULL,
		uploaded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		
		-- Foreign key constraints (nullable)
		CONSTRAINT fk_asset_documents_asset_id 
			FOREIGN KEY (asset_id) REFERENCES assets(id) 
			ON DELETE SET NULL ON UPDATE CASCADE,
			
		-- Unique constraint on cloudinary_id
		CONSTRAINT unique_cloudinary_id UNIQUE (cloudinary_id)
	);

	-- Create indexes for better query performance
	CREATE INDEX IF NOT EXISTS idx_asset_documents_tenant_id ON asset_documents(tenant_id);
	CREATE INDEX IF NOT EXISTS idx_asset_documents_asset_id ON asset_documents(asset_id);
	CREATE INDEX IF NOT EXISTS idx_asset_documents_document_type ON asset_documents(document_type);
	CREATE INDEX IF NOT EXISTS idx_asset_documents_uploaded_at ON asset_documents(uploaded_at);
	CREATE INDEX IF NOT EXISTS idx_asset_documents_cloudinary_id ON asset_documents(cloudinary_id);
	
	-- Create composite index for tenant-scoped queries
	CREATE INDEX IF NOT EXISTS idx_asset_documents_tenant_asset ON asset_documents(tenant_id, asset_id) WHERE tenant_id IS NOT NULL AND asset_id IS NOT NULL;
	`

	// Execute the SQL
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create asset_documents table: %v", err)
	}

	log.Println("Asset documents table created successfully")
	return nil
}

// CreateAssetDocumentTableIfNotExists creates the asset_documents table if it doesn't exist
func CreateAssetDocumentTableIfNotExists(db *sql.DB) error {
	// Check if table exists
	var exists bool
	query := `SELECT EXISTS (
		SELECT 1 FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_name = 'asset_documents'
	)`

	err := db.QueryRow(query).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if asset_documents table exists: %v", err)
	}

	if exists {
		log.Println("Asset documents table already exists")
		return nil
	}

	log.Println("Creating asset_documents table...")

	// SQL untuk membuat tabel asset_documents
	createTableSQL := `
	CREATE TABLE asset_documents (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		tenant_id UUID NULL,
		asset_id UUID NULL,
		document_type VARCHAR(100) NOT NULL,
		file_url VARCHAR(1000) NOT NULL,
		cloudinary_id VARCHAR(500) NOT NULL,
		original_filename VARCHAR(500) NOT NULL,
		file_size BIGINT NOT NULL DEFAULT 0,
		mime_type VARCHAR(100) NOT NULL,
		uploaded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		
		-- Foreign key constraints (nullable)
		CONSTRAINT fk_asset_documents_asset_id 
			FOREIGN KEY (asset_id) REFERENCES assets(id) 
			ON DELETE SET NULL ON UPDATE CASCADE,
			
		-- Unique constraint on cloudinary_id
		CONSTRAINT unique_cloudinary_id UNIQUE (cloudinary_id)
	);

	-- Create indexes for better query performance
	CREATE INDEX idx_asset_documents_tenant_id ON asset_documents(tenant_id);
	CREATE INDEX idx_asset_documents_asset_id ON asset_documents(asset_id);
	CREATE INDEX idx_asset_documents_document_type ON asset_documents(document_type);
	CREATE INDEX idx_asset_documents_uploaded_at ON asset_documents(uploaded_at);
	CREATE INDEX idx_asset_documents_cloudinary_id ON asset_documents(cloudinary_id);
	
	-- Create composite index for tenant-scoped queries
	CREATE INDEX idx_asset_documents_tenant_asset ON asset_documents(tenant_id, asset_id) WHERE tenant_id IS NOT NULL AND asset_id IS NOT NULL;
	`

	// Execute the SQL
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create asset_documents table: %v", err)
	}

	log.Println("Asset documents table created successfully")
	return nil
}
