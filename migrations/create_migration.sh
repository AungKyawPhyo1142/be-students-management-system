#!/bin/bash

# Get the current timestamp in YYYYMMDDHHMMSS format
timestamp=$(date +"%Y%m%d%H%M%S")

# Ask for a migration name
echo "Enter migration name (e.g., add_phone_column_to_students):"
read name

# Create a new migration file with the timestamp and name
filename="migration_${timestamp}_${name}.go"
touch $filename

# Add a template for the migration
cat <<EOL > $filename
package migrations

import (
	"gorm.io/gorm"
)

func Migrate_${timestamp}_${name}(tx *gorm.DB) error {
	// TODO: Implement the migration logic
	return nil
}

func Rollback_${timestamp}_${name}(tx *gorm.DB) error {
	// TODO: Implement the rollback logic
	return nil
}
EOL

echo "Migration file created: $filename"
