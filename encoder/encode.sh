#!/bin/sh

set -o pipefail

if [ -z "$MINIO_HOST" ]; then
        echo "MINIO_HOST env var is empty"
        exit 1
fi

if [ -z "$MINIO_USER" ]; then
        echo "MINIO_USER env var is empty"
        exit 1
fi

if [ -z "$MINIO_PASS" ]; then
        echo "MINIO_PASS env var is empty"
        exit 1
fi

if [ -z "$SOURCE" ]; then
        echo "SOURCE env var is empty"
        exit 1
fi

if [ -z "$TARGET" ]; then
        echo "TARGET env var is empty"
        exit 1
fi

if [ -z "$FILENAME" ]; then
        echo "FILENAME env var is empty"
        exit 1
fi



mcli alias set target "http://$MINIO_HOST" "$MINIO_USER" "$MINIO_PASS"

curl "$SOURCE" > /tmp/source

dca /tmp/source /tmp/output.dca

mcli cp /tmp/output.dca target/$TARGET/$FILENAME.dca 

redis-cli \
        -h "$REDIS_HOSTNAME" \
        -p "$REDIS_PORT" \
        -a "$REDIS_PASSWORD" \
        -n "$REDIS_DATABASE_INDEX" \
        PUBLISH "$REDIS_CHAN" "$REDIS_PAYLOAD"
