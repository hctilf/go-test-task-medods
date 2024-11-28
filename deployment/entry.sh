#!/bin/sh

envsubst < ${CONFIG_NAME}.${CONFIG_EXT}.template > ${CONFIG_NAME}.${CONFIG_EXT}

echo "Environment variables substituted and saved to config.yaml"

exec "$@"