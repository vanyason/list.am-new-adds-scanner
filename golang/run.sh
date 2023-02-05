#!/bin/bash

clear && clear && go test -v ./... && go run cmd/listam_parser/*.go -p 200000