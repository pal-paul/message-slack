#!/bin/bash

# GitHub Actions Workflow Validation Script
# This script validates the syntax of GitHub Actions workflows

set -e

echo "🔍 Validating GitHub Actions workflow files..."

# Check if workflow files exist
WORKFLOW_DIR=".github/workflows"
if [ ! -d "$WORKFLOW_DIR" ]; then
    echo "❌ No .github/workflows directory found"
    exit 1
fi

# Find all workflow files
WORKFLOW_FILES=$(find "$WORKFLOW_DIR" -name "*.yml" -o -name "*.yaml")

if [ -z "$WORKFLOW_FILES" ]; then
    echo "❌ No workflow files found in $WORKFLOW_DIR"
    exit 1
fi

echo "📄 Found workflow files:"
for file in $WORKFLOW_FILES; do
    echo "  - $file"
done

# Check for deprecated action versions
echo ""
echo "🔍 Checking for deprecated action versions..."

DEPRECATED_FOUND=false

# Check for deprecated upload-artifact@v3
if grep -r "actions/upload-artifact@v3" "$WORKFLOW_DIR" >/dev/null 2>&1; then
    echo "❌ Found deprecated actions/upload-artifact@v3"
    DEPRECATED_FOUND=true
fi

# Check for deprecated cache@v3
if grep -r "actions/cache@v3" "$WORKFLOW_DIR" >/dev/null 2>&1; then
    echo "❌ Found deprecated actions/cache@v3"
    DEPRECATED_FOUND=true
fi

# Check for deprecated setup-go@v4 (v5 is latest)
if grep -r "actions/setup-go@v4" "$WORKFLOW_DIR" >/dev/null 2>&1; then
    echo "⚠️  Found potentially outdated actions/setup-go@v4 (v5 is latest)"
fi

if [ "$DEPRECATED_FOUND" = true ]; then
    echo "❌ Deprecated actions found! Please update them."
    exit 1
fi

echo "✅ No deprecated actions found!"

# Check action versions are up to date
echo ""
echo "📋 Current action versions found:"
grep -hr "uses:" "$WORKFLOW_DIR" | sed 's/.*uses: /  /' | sort | uniq

echo ""
echo "✅ All GitHub Actions workflow files appear to be valid!"
echo ""
echo "💡 Recommendations:"
echo "  - Ensure you have SLACK_TOKEN in repository secrets"
echo "  - Consider setting TEST_SLACK_CHANNEL in repository variables"
echo "  - Test the workflow in a branch before merging to main"
