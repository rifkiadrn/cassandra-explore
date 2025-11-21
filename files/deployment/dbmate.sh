#!/bin/sh
# dbmate (DB Migration) Config
if [[ -z "${DATABASE_URL}" ]]; then
  export DATABASE_URL=postgres://$DB_USERNAME:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable
fi
export DBMATE_DEFAULT_MIGRATIONS_DIR='./db/postgres/dbmate/migrations'
export DBMATE_DEFAULT_MIGRATIONS_TABLE="dbmate_schema_migrations"
export DBMATE_POSTDEPLOYMENT_MIGRATIONS_DIR='./db/postgres/dbmate/postdeployment_migrations'
export DBMATE_POSTDEPLOYMENT_MIGRATIONS_TABLE="dbmate_schema_postdeployment_migrations"
export DBMATE_SCHEMA_FILE='./db/postgres/dbmate/schema'
export DBMATE_NO_DUMP_SCHEMA=false
export DBMATE_WAIT=false
export DBMATE_WAIT_TIMEOUT='2m0s'
export DBMATE_VERBOSE=true

case "$DBMATE_MIGRATION_TYPE" in
   "default")
    export DBMATE_MIGRATIONS_DIR=$DBMATE_DEFAULT_MIGRATIONS_DIR
    export DBMATE_MIGRATIONS_TABLE=$DBMATE_DEFAULT_MIGRATIONS_TABLE
   ;;
   "postdeployment") 
    export DBMATE_MIGRATIONS_DIR=$DBMATE_POSTDEPLOYMENT_MIGRATIONS_DIR
    export DBMATE_MIGRATIONS_TABLE=$DBMATE_POSTDEPLOYMENT_MIGRATIONS_TABLE
    ;;
   *)
    export DBMATE_MIGRATIONS_DIR=$DBMATE_DEFAULT_MIGRATIONS_DIR
    export DBMATE_MIGRATIONS_TABLE=$DBMATE_DEFAULT_MIGRATIONS_TABLE
    ;;
esac
