#!/bin/sh
# Helper script for db-operation workflow
set -e

case "$DBMATE_OPERATION_MIGRATION_TYPE" in
   "default")
   ;;
   "postdeployment") 
    ;;
   *) 
    echo "Unknown migration type $DBMATE_OPERATION_MIGRATION_TYPE !!" 
    exit 1
    ;;
esac

export DBMATE_MIGRATION_TYPE=$DBMATE_OPERATION_MIGRATION_TYPE
source ./files/deployment/dbmate.sh

dbmate --migrations-dir $DBMATE_MIGRATIONS_DIR --migrations-table $DBMATE_MIGRATIONS_TABLE $DBMATE_OPERATION_COMMAND
