#!/bin/bash

set -e

#
# Ensures that:
#   1. The latest version of Docker is installed.
#   2. Experimental mode and buildkit are enabled.
#
# $1 - The Docker version.
#
# Examples:
#
#   ensureDocker
#
ensureDocker() {
  if [[ "${OSTYPE}" != "linux-gnu" ]]; then
    printf "Unsupported host OS %s\n" "${OSTYPE}"
    exit 1
  fi

  curl -fsSL https://get.docker.com -o get-docker.sh
  sudo sh get-docker.sh

  echo '{"experimental":true,"features":{"buildkit":true}}' \
    | sudo tee /etc/docker/daemon.json

  sudo service docker restart
}

#
# Builds platform-specific Docker images.
#
# $1 - The application folder.
# $2 - The image name.
# $3 - The image tag.
# $4 - The target platforms.
#
# Examples:
#
#   build "/go/src/github.com/foo/bar" "foo/bar" "1.0.0" linux/amd64,linux/arm64"
#
build() {
  declare -r app_folder="${1}"
  declare -r image="${2}"
  declare -r tag="${3}"
  declare -r platforms=($(echo "${4}" | tr ',' '\n'))

  if [[ "${tag}" =~ ^v?([0-9]+\.[0-9]+\.[0-9]+) ]]; then
    local version="${BASH_REMATCH[1]}"
  else
    printf "Tag %s is not a semantic version\n" "${tag}"
    exit 1
  fi

  for platform in "${platforms[@]}"; do
    # Form a platform tag, e.g. "1.0.0-linux-amd64".
    local platform_tag="${version}-${platform//\//-}"

    # Build a platform spceific Docker image.
    docker build --platform="${platform}" \
      --build-arg APP_FOLDER="${app_folder}" \
      --tag "${image}:${platform_tag}" \
      .

  done
}

"$@"
