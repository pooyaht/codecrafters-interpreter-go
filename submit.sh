#!/bin/bash
set -e 

if git diff --quiet gh_upstream/$(git branch --show-current) HEAD 2>/dev/null; then
    echo "ℹ️  No new commits to push. Nothing to submit."
    exit 0
fi

echo "🧪 Running CodeCrafters tests..."

if codecrafters test; then
    echo "✅ Tests passed!"
    
    echo "🚀 Pushing to GitHub..."
    git push gh_upstream
    echo "✅ Pushed to GitHub"
    
    echo "📤 Submitting to CodeCrafters..."
    codecrafters submit
    echo "🎉 Submitted to CodeCrafters!"
    
else
    echo "❌ Tests failed! Please fix the issues before submitting."
    exit 1
fi