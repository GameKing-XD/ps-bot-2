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


if [ ! -z "$POST_HOOK" ]; then
        if [ -z "$AMQP_SERVER" ]; then
                echo "AMQP_SERVER env var is empty"
                exit 1
        fi

        if [ -z "$AMQP_USER" ]; then
                echo "AMQP_USER env var is empty"
                exit 1
        fi

        if [ -z "$AMQP_PASS" ]; then
                echo "AMQP_PASS env var is empty"
                exit 1
        fi

        if [ -z "$AMQP_CHAN" ]; then
                echo "AMQP_CHAN env var is empty"
                exit 1
        fi
        if [ -z "$AMQP_BODY" ]; then
                echo "AMQP_BODY env var is empty"
                exit 1
        fi

        if [ -z "$AMQP_CONTENT_TYPE" ];then
                AMQP_CONTENT_TYPE="application/json"
        fi


        amqp-publish \
                --username "$AMQP_USER" \
                --password "$AMQP_PASS" \
                -s "$AMQP_SERVER" \
                -b "$AMQP_BODY" \
                -r "$AMQP_CHAN" \
                -e "$AMQP_EXCHANGE"

fi


