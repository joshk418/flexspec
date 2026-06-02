package migrate

import "io/fs"

// Registry returns ordered migrations. templatesFS is the embedded templates
// subtree (rooted at "templates"). force allows overwriting differing template files.
func Registry(templatesFS fs.FS, force bool) []Migration {
	return []Migration{
		&statusRenameMigration{},
		&templatesResyncMigration{templates: templatesFS, force: force},
		&configKeysMigration{},
		&charterCheckMigration{templates: templatesFS},
	}
}
