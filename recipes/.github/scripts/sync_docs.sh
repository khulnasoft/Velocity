#!/usr/bin/env bash
set -e

# Some env variables
BRANCH="main"
REPO_URL="github.com/khulnasoft/docs.git"
AUTHOR_EMAIL="github-actions[bot]@users.noreply.github.com"
AUTHOR_USERNAME="github-actions[bot]"
REPO_DIR="recipes"
COMMIT_URL="https://go.khulnasoft.com/velocity/recipes"

# Set commit author
git config --global user.email "${AUTHOR_EMAIL}"
git config --global user.name "${AUTHOR_USERNAME}"

git clone https://${TOKEN}@${REPO_URL} velocity-docs

latest_commit=$(git rev-parse --short HEAD)

# remove all files in the docs directory
rm -rf velocity-docs/docs/${REPO_DIR}/*

# Find and copy relevant files
find . \
  -type f \
  \( -iname '*.md' -o -iname '*.png' -o -iname '*.jpg' -o -iname '*.jpeg' -o -iname '*.gif' -o -iname '*.bmp' -o -iname '*.svg' -o -iname '*.webp' \) \
  -not -path "./velocity-docs/*" \
  -not -path "*/vendor/*" \
  -not -path "*/.github/*" \
  -not -path "*/.*" |
while IFS= read -r f; do
  log_output=$(git log --oneline "${BRANCH}" HEAD~1..HEAD --name-status -- "${f}" || true)

  if [[ -n $log_output || ! -f "velocity-docs/docs/${REPO_DIR}/${f}" ]]; then
    mkdir -p velocity-docs/docs/${REPO_DIR}/$(dirname "$f")
    cp "$f" velocity-docs/docs/${REPO_DIR}/$f
  fi
done

# Push changes
cd velocity-docs/ || true
git add .

git commit -m "Add docs from ${COMMIT_URL}/commit/${latest_commit}"

MAX_RETRIES=5
DELAY=5
retry=0

while ((retry < MAX_RETRIES))
do
    git push https://${TOKEN}@${REPO_URL} && break
    retry=$((retry + 1))
    git pull --rebase
    sleep $DELAY
done

if ((retry == MAX_RETRIES))
then
    echo "Failed to push after $MAX_RETRIES attempts. Exiting with 1."
    exit 1
fi
