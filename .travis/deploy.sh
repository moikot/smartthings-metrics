#!/bin/bash

set -e

#
# Creates and pushes a multi-platform manifest.
#
# $1 - The image name.
# $2 - The platform-agnostic tag.
# $3 - The platform-specific tags.
#
# Examples:
#
#   pushManifest "foo/bar" "1.0.0" "foo/bar:1.0.0-linux-arm64-v8"
#
pushManifest() {
  declare -r image="${1}"
  declare -r tag="${2}"
  declare -r manifests=(${3})

  docker manifest create --amend "${image}:${tag}" ${manifests[@]}

  for manifest in "${manifests[@]}"; do
     local name_parts=($(echo "${manifest}" | tr ':' '\n'))
     local annotations=($(echo "${name_parts[1]}" | tr '-' '\n'))

     if [[ "${annotations[3]}" != "" ]]; then
       local variant="--variant ${annotations[3]}"
     fi

     docker manifest annotate \
       --os "${annotations[1]}" \
       --arch "${annotations[2]}" \
       ${variant} "${image}:${tag}" "${manifest}"

  done

  docker manifest push --purge "${image}:${tag}"
}

#
# Creates and pushes multi-platform manifests using semantic versioning.
#
# $1 - The image name.
# $2 - The version.
# $3 - The target platforms.
#
# Examples:
#
#   pushManifests "foo/bar" "1.0.0" "linux/amd64,linux/arm64/v8"
#
pushManifests() {
  declare -r image="${1}"
  declare -r version="${2}"
  declare -r platforms=($(echo "${3}" | tr ',' '\n'))

  local manifests
  for platform in "${platforms[@]}"; do
    local platform_tag="${version}-${platform//\//-}"
    manifests="${manifests} ${image}:${platform_tag}"
  done

  if [[ "${version}" =~ ^([0-9]+)\.([0-9]+)\.([0-9]+) ]]; then
    local major="${BASH_REMATCH[1]}"
    local minor="${BASH_REMATCH[2]}"
    local patch="${BASH_REMATCH[3]}"

    pushManifest "${image}" "${major}" "${manifests}"
    pushManifest "${image}" "${major}.${minor}" "${manifests}"
    pushManifest "${image}" "${major}.${minor}.${patch}" "${manifests}"
    pushManifest "${image}" latest "${manifests}"
  else
    printf "Version %s is not a semantic version\n" "${version}"
    exit 1
  fi
}

#
# Deletes a tag on a remote server using Docker API v2.
#
# $1 - The image name.
# $2 - The image tag.
# $3 - The JWT.
#
# Examples:
#
#   deleteTag "foo/bar" "1.0.0-linux-amd64" "token"
#
deleteTag() {
  declare -r image="${1}"
  declare -r tag="${2}"
  declare -r token="${3}"

  local code=$(curl -s -o /dev/null -LI -w "%{http_code}" \
    https://hub.docker.com/v2/repositories/"${image}"/tags/"${tag}"/ \
    -X DELETE \
    -H "Authorization: JWT ${token}")

  if [[ "${code}" = "204" ]]; then
    printf "Successfully deleted %s\n" "${image}:${tag}"
  else
    printf "Unable to delete %s, response code: %s\n" "${image}:${tag}" "${code}"
    exit 1
  fi
}

#
# Deletes platform-spceific image tags.
#
# $1 - The image name.
# $2 - The version.
# $3 - The JWT.
# $4 - The target platforms.
#
# Examples:
#
#   deleteTags "foo/bar" "1.0.0" "token" "linux/amd64,linux/arm64/v8"
#
deleteTags() {
  declare -r image="${1}"
  declare -r version="${2}"
  declare -r token="${3}"
  declare -r platforms=($(echo "${4}" | tr ',' '\n'))

  for platform in "${platforms[@]}"; do
    local platform_tag="${version}-${platform//\//-}"
    deleteTag "${image}" "${platform_tag}" "${token}"
  done
}

#
# Pushes multi-platfomrm set of images and creates manifests.
#
# $1 - The image name.
# $2 - The tag.
# $3 - The target platforms.
#
# Examples:
#
#   push_images "foo/bar" "1.0.0" "linux/amd64,linux/arm64"
#
pushImages() {
  declare -r image="${1}"
  declare -r tag="${2}"
  declare -r platforms=($(echo "${3}" | tr ',' '\n'))

  # Login into Docker repository
  echo "${DOCKER_PASSWORD}" | docker login -u "${DOCKER_USERNAME}" --password-stdin

  if [[ "${tag}" =~ ^v?([0-9]+\.[0-9]+\.[0-9]+) ]]; then
    local version="${BASH_REMATCH[1]}"
  else
    printf "Tag %s is not a semantic version\n" "${tag}"
    exit 1
  fi

  for platform in "${platforms[@]}"; do
    local platform_tag="${version}-${platform//\//-}"
    docker push "${image}:${platform_tag}"
  done

  pushManifests "${image}" "${version}" "${3}"

  local token=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d '{"username": "'"$DOCKER_USERNAME"'", "password": "'"$DOCKER_PASSWORD"'"}' \
    https://hub.docker.com/v2/users/login/ | jq -r .token)

  deleteTags "${image}" "${version}" "${token}" "${3}"
}

#
# Pushes description to Docker Hub.
#
# $1 - The image name.
# $2 - The file name.
#
# Examples:
#
#   push_images "foo/bar" "README.md"
#
pushDescription() {
  declare -r image="${1}"
  declare -r file_name="${2}"

  local token=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d '{"username": "'"$DOCKER_USERNAME"'", "password": "'"$DOCKER_PASSWORD"'"}' \
    https://hub.docker.com/v2/users/login/ | jq -r .token)

  local code=$(jq -n --arg msg "$(<${file_name})" \
    '{"registry":"registry-1.docker.io","full_description": $msg }' | \
        curl -s -o /dev/null  -L -w "%{http_code}" \
           https://cloud.docker.com/v2/repositories/"${image}"/ \
           -d @- -X PATCH \
           -H "Content-Type: application/json" \
           -H "Authorization: JWT ${token}")

  if [[ "${code}" = "200" ]]; then
    printf "Successfully pushed %s to Docker Hub\n" "${file_name}"
  else
    printf "Unable to push %s to Docker Hub, response code: %s\n" "${file_name}" "${code}"
    exit 1
  fi
}

"$@"
