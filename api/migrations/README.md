## How the migration works
On startup, the system will load all sql scripts from this folder.\
It will loop them (starting with 0) and checks if it has been applied to the database.
If this is not the case, it will apply it.

## How to add a migration
If you want to add a migration after 0_init, create one with the name ``1_something``.\
***The migration script must finish with the sql statement ``INSERT INTO db_state VALUES (1);``; where `1` is the identifier of your migration.***