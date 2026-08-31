#!/bin/bash
ENV_VARS=$(cat .env)
sh -c "$ENV_VARS go run cmd/main.go"