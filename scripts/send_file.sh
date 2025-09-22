#!/bin/bash

SFTP_USER=$1
SFTP_HOST=$2
LOCAL_FILE=$3
REMOTE_PATH=$4

sftp $SFTP_USER@$SFTP_HOST <<EOF
  put "$LOCAL_FILE" "$REMOTE_PATH"
  exit
EOF

if [ $? -eq 0 ]; then
  echo "File '$LOCAL_FILE' successfully sent to $SFTP_HOST:$REMOTE_PATH"
else
  echo "Error sending file '$LOCAL_FILE' to $SFTP_HOST:$REMOTE_PATH"
fi
