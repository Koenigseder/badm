#!/bin/bash

# Needed env vars:
# GITHUB_REPOSITORY - Name of GitHub repository
# TAG - SemVer tag
# TOKEN - Token used for uploading binary to the release

set -eu

REPO_API="https://api.github.com/repos/$GITHUB_REPOSITORY"
TAGS_API="$REPO_API/releases/tags/$TAG"
AUTH="Authorization: Bearer $TOKEN"

ID=$(curl -L -H "$AUTH" "$TAGS_API" | jq -r '.ID')

if [ "$ID" = "null" ]; then
	echo "Release not found!"
	exit 1
fi

upload() {
		filename="$1"

		# Upload asset
		echo "Uploading asset $filename to release $TAG with ID $ID... "

		# Construct URL
		url="https://uploads.github.com/repos/$GITHUB_REPOSITORY/releases/$ID/assets?name=$(basename "$filename")"

		echo "Constructed URL $url"

		curl -L -X POST -H "$AUTH" -H "Content-Type: application/octet-stream" -H "Accept: application/vnd.github+json" -H "X-GitHub-Api-Version: 2022-11-28" --data-binary "@$filename" "$url"
}

find "build" -name "badm*" | while read -r file; do upload "$file"; done
