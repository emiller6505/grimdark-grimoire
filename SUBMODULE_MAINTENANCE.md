# Git Submodule Maintenance Guide

## Overview

The `wh40k-10e/` directory is a git submodule pointing to the BSData repository. This guide explains how to keep it updated.

## Checking for Updates

```bash
# Check if the submodule has updates available
cd wh40k-10e
git fetch
git log HEAD..origin/main --oneline  # See what's new upstream
cd ..
```

Or use git submodule commands:
```bash
# Check submodule status
git submodule status

# See if submodule is behind upstream
git submodule foreach git fetch
git submodule foreach git log HEAD..origin/main --oneline
```

## Updating the Submodule

### Option 1: Update to Latest (Recommended)

```bash
# Update submodule to latest commit from upstream
git submodule update --remote wh40k-10e

# Review the changes
cd wh40k-10e
git log --oneline -10  # See recent commits
cd ..

# Commit the submodule update to your main repo
git add wh40k-10e
git commit -m "Update wh40k-10e submodule to latest"
```

### Option 2: Update to Specific Tag/Commit

```bash
# Update to a specific tag (e.g., latest release)
cd wh40k-10e
git checkout v10.6.0  # or whatever tag you want
cd ..

# Commit the update
git add wh40k-10e
git commit -m "Update wh40k-10e submodule to v10.6.0"
```

### Option 3: Manual Update

```bash
cd wh40k-10e
git pull origin main  # or whatever branch you're tracking
cd ..

# Commit the update
git add wh40k-10e
git commit -m "Update wh40k-10e submodule"
```

## Initial Clone (for new contributors)

When someone clones your repository, they need to initialize submodules:

```bash
# Clone with submodules
git clone --recurse-submodules <your-repo-url>

# Or if already cloned
git submodule update --init --recursive
```

## Checking Current Version

```bash
# See what commit/tag the submodule is on
cd wh40k-10e
git describe --tags
git log -1 --oneline
cd ..
```

## Best Practices

1. **Regular Updates**: Check for updates periodically (weekly/monthly)
2. **Test After Updates**: Run your tests after updating the submodule
3. **Commit Separately**: Commit submodule updates separately from code changes
4. **Document Versions**: Consider documenting which submodule version you're using in your README

## Automated Updates (Optional)

You can set up a GitHub Action or script to check for updates:

```bash
#!/bin/bash
# update-submodule.sh
cd wh40k-10e
git fetch
if [ $(git rev-list HEAD..origin/main --count) != 0 ]; then
    echo "Updates available!"
    git submodule update --remote wh40k-10e
    cd ..
    git add wh40k-10e
    git commit -m "Update wh40k-10e submodule"
    echo "Submodule updated and committed"
else
    echo "Submodule is up to date"
fi
```

## Troubleshooting

### Submodule shows as modified but you haven't changed anything

```bash
# This usually means the submodule is on a different commit than what's recorded
cd wh40k-10e
git status  # Check if you're ahead/behind
git checkout <commit-hash-from-main-repo>  # Reset to what main repo expects
```

### Submodule is out of sync

```bash
# Reset submodule to match what's recorded in main repo
git submodule update --init --recursive
```

### Remove and re-add submodule (if needed)

```bash
# Remove submodule
git submodule deinit wh40k-10e
git rm wh40k-10e
rm -rf .git/modules/wh40k-10e

# Re-add it
git submodule add https://github.com/BSData/wh40k-10e.git wh40k-10e
```

