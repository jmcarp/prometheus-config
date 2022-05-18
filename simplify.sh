#!/usr/bin/env bash

set -euo pipefail

# Delete modules we don't need
find -E . -type d -mindepth 1 -maxdepth 1 -not -regex "./(config|discovery|model|prompb|promql|scrape|storage|template|tsdb|util)" -not -name ".*" -exec git rm -rf {} \;
find -E discovery -type d -mindepth 1 -maxdepth 1 -not -regex "discovery/(kubernetes|targetgroup)" -exec git rm -rf {} \;
find -E discovery/kubernetes -type f -not -name "kubernetes.go" -exec git rm {} \;
find . -type f -name "*_test.go" -exec git rm {} \;

# Drop dependencies from kubernetes discovery package to avoid conflicts
git apply kubernetes-discovery.patch

# Downgrade otel to avoid conflicts in github.com/go-logr/logr
go get "go.opentelemetry.io/otel@v1.2.0"
go get "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp@v0.27.0"

# Update go.mod
go mod tidy
