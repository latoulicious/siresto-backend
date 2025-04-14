package main

import (
	"flag"
	"log"

	"github.com/joho/godotenv"
	"github.com/latoulicious/siresto-backend/internal/config"
	"github.com/latoulicious/siresto-backend/migrations"
)

func main() {
	// Parse command flags
	seedFlag := flag.Bool("seed", false, "Seed the database with initial data after migration")
	resetFlag := flag.Bool("reset", false, "Drop all tables before migration (DESTRUCTIVE)")
	flag.Parse()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found, using environment variables")
	}

	// Connect to database
	db, err := config.NewGormDB()
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ Failed to get sql.DB instance: %v", err)
	}
	defer sqlDB.Close()
	log.Println("✅ Connected to database")

	// Handle reset if requested
	if *resetFlag {
		log.Println("🔄 Dropping existing tables...")
		// Enable cascade to avoid foreign key issues
		db.Exec("SET session_replication_role = 'replica';")

		// Drop all tables - BE CAREFUL WITH THIS!
		err := db.Exec(`
			DO $$ DECLARE
				r RECORD;
			BEGIN
				FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = current_schema()) LOOP
					EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';
				END LOOP;
			END $$;
		`).Error

		// Reset to normal role
		db.Exec("SET session_replication_role = 'origin';")

		if err != nil {
			log.Fatalf("❌ Failed to drop tables: %v", err)
		}
		log.Println("✅ Tables dropped successfully")
	}

	// Run migrations
	log.Println("🔄 Running migrations...")
	if err := migrations.RunMigrations(db); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}
	log.Println("✅ Migrations completed successfully")

	// Seed data if requested
	if *seedFlag {
		log.Println("🌱 Seeding initial data...")
		if err := migrations.SeedData(db); err != nil {
			log.Fatalf("❌ Data seeding failed: %v", err)
		}
		log.Println("✅ Data seeding completed successfully")
	}

	log.Println("✨ Migration process completed")
}
