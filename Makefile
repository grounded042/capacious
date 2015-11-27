setup:
	psql -d capacious-dev -a -f bin/setup/sql/drop_public_schema.sql
	psql -d capacious-dev -a -f bin/setup/sql/base_db.sql
	psql -d capacious-dev -a -f bin/setup/sql/seed_dev_data.sql

test: setup
	cd e2e-tests/; \
	npm run test
	
