#!/bin/bash

VERBOSE=${VERBOSE:-"0"}
V=""
if [[ "${VERBOSE}" == "1" ]]; then
  V="-x"
  set -x
fi

PathPwd="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
PathDist=${1:?"output path"}
PathBuild=${3:?"path to build"}
PackageProject="$(cd "$(dirname "${PathBuild}")/.." && go list -m)"
PackageVersion="${PackageProject}/${2:?"package/version"}"

set -e

# compiler
GOPKG="$GOPATH/pkg"
GOBINARY=${GOBINARY:-go}
# compile options
GOOS=${4:-$(go env GOOS)} # linux
GOARCH=${5:-$(go env GOARCH)}
GOBUILDFLAGS="-trimpath"
GOBUILDFLAGS=${GOBUILDFLAGS:-""}
LDFLAGS="-extldflags -static -w"
GCFLAGS=${GCFLAGS:-}

BUILDINFO=${BUILDINFO:-""}
STATIC=${STATIC:-1}
export CGO_ENABLED=0

if [[ "${STATIC}" != "1" ]]; then
  LDFLAGS=""
fi

# gather buildinfo if not already provided
# For a release build BUILDINFO should be produced
# at the beginning of the build and used throughout
if [[ -z ${BUILDINFO} ]]; then
  BUILDINFO=$(mktemp)
  bash ${PathPwd}/version.sh >${BUILDINFO}
fi

# BUILD LD_VERSIONFLAGS
LD_VERSIONFLAGS=""
while read line; do
  read SYMBOL VALUE < <(echo $line)
  # shellcheck disable=SC2089
  LD_VERSIONFLAGS="${LD_VERSIONFLAGS} -X '${PackageVersion}.${SYMBOL}=${VALUE}'"
done <"${BUILDINFO}"

# forgoing -i (incremental build) because it will be deprecated by tool chain.
time GOOS=${GOOS} GOARCH=${GOARCH} ${GOBINARY} build ${V} ${GOBUILDFLAGS} ${GCFLAGS:+-gcflags "${GCFLAGS}"} -o "${PathDist}" \
  -pkgdir=${GOPKG}/${GOOS}_${GOARCH} -ldflags "${LDFLAGS} ${LD_VERSIONFLAGS}" "${PathBuild}"
