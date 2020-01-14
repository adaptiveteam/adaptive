#!/bin/bash
rename_table_backup_idempotent() {
    local original_table_name=$1
    local new_table_name=$2
    if [ ! -d "dump/${ADAPTIVE_CLIENT_ID}${new_table_name}" ] ; then
      echo "Renaming ${ADAPTIVE_CLIENT_ID}${original_table_name} to ${ADAPTIVE_CLIENT_ID}${new_table_name}"
      mv dump/${ADAPTIVE_CLIENT_ID}${original_table_name} dump/${ADAPTIVE_CLIENT_ID}${new_table_name}
      python dynamodump/dynamodump.py -m restore -r ${AWS_REGION} --dataOnly -s ${ADAPTIVE_CLIENT_ID}${new_table_name}
    fi
}

rename_table_backup_idempotent "_user_objectives" "_user_objective"
