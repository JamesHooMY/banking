#!/bin/bash

# Wait for the MySQL master to be fully started
until mysqladmin ping -h mysql-master -uroot -proot_password --silent; do
  >&2 echo "Master is unavailable - sleeping"
  sleep 5
done

# Get the log file and position from the master
MASTER_STATUS=$(mysql -h mysql-master -uroot -proot_password -e "SHOW MASTER STATUS\G")
MASTER_LOG_FILE=$(echo "$MASTER_STATUS" | grep 'File:' | awk '{print $2}')
MASTER_LOG_POS=$(echo "$MASTER_STATUS" | grep 'Position:' | awk '{print $2}')

# Debugging information
echo "Master log file: $MASTER_LOG_FILE"
echo "Master log position: $MASTER_LOG_POS"

# Temporarily disable super_read_only mode
mysql -uroot -proot_password -e "SET GLOBAL super_read_only = OFF;" || { echo "Failed to disable super_read_only"; exit 1; }

# Configure the slave to replicate from the master
mysql -uroot -proot_password -e "CHANGE MASTER TO MASTER_HOST='mysql-master', MASTER_USER='replica', MASTER_PASSWORD='replica_password', MASTER_LOG_FILE='$MASTER_LOG_FILE', MASTER_LOG_POS=$MASTER_LOG_POS;" || { echo "Failed to change master"; exit 1; }
mysql -uroot -proot_password -e "START SLAVE;" || { echo "Failed to start slave"; exit 1; }

# Re-enable super_read_only mode
mysql -uroot -proot_password -e "SET GLOBAL super_read_only = ON;" || { echo "Failed to enable super_read_only"; exit 1; }

echo "Slave initialization completed successfully."
