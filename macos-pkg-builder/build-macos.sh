#!/bin/bash

# SPDX-License-Identifier: Apache-2.0
#
# The OpenSearch Contributors require contributions made to
# this file be licensed under the Apache-2.0 license or a
# compatible open source license.
#
# Modifications Copyright OpenSearch Contributors. See
# GitHub history for details.

SCRIPTPATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
PRODUCT="OpenSearchCLI"
DATE=`date +%Y-%m-%d`
TIME=`date +%H:%M:%S`
LOG_PREFIX="[$DATE $TIME]"

usage() {
    echo ""
    echo "This script is used to build the macos package for OpenSearch CLI."
    echo "--------------------------------------------------------------------------"
    echo "Usage: $0 [args]"
    echo ""
    echo "Required arguments:"
    echo -e "-v VERSION\tSpecify the OpenSearch CLI version number that you are building, e.g. '1.0.0' or '1.0.0-beta1'."
    echo -e "-a ARCHITECTURE\tSpecify the architecture type your signed binary is based on, e.g. 'x64,arm64'."
    echo -e "-t TARGET DIRECTORY\tSpecify the location where you like to generate package artifact."
    echo "--------------------------------------------------------------------------"
}

while getopts ":hv:a:t:" arg; do
    case $arg in
        h)
            usage
            exit 1
            ;;
        v)
            VERSION=$OPTARG
            ;;
        a)
            ARCHITECTURE=$OPTARG
            ;;
        t)
            TARGET_DIRECTORY=$OPTARG
            ;;
        :)
            echo "-${OPTARG} requires an argument"
            usage
            exit 1
            ;;
        ?)
            echo "Invalid option: -${arg}"
            exit 1
            ;;
    esac
done

# Validate the required parameters to present
if [ -z "$VERSION" ] || [ -z "$ARCHITECTURE" ] || [ -z "$TARGET_DIRECTORY" ]; then
  echo "You must specify '-v VERSION', '-a ARCHITECTURE', '-t TARGET_DIRECTORY'"
  usage
  exit 1
fi
if [ "$ARCHITECTURE" != "x64" ] && [ "$ARCHITECTURE" != "arm64" ]; then
	echo "We only support 'x64' and 'arm64' as architecture name for -a parameter"
  exit 1
fi

INSTALLER_NAME="opensearch-cli-${VERSION}-macos-${ARCHITECTURE}"
#Functions
go_to_dir() {
    pushd $1 >/dev/null 2>&1
}

log_info() {
    echo "${LOG_PREFIX}[INFO]" $1
}

log_error() {
    echo "${LOG_PREFIX}[ERROR]" $1
}

deleteInstallationDirectory() {
    log_info "Cleaning $TARGET_DIRECTORY directory."
    rm -rf $TARGET_DIRECTORY

    if [[ $? != 0 ]]; then
        log_error "Failed to clean $TARGET_DIRECTORY directory" $?
        exit 1
    fi
}

createInstallationDirectory() {
    if [ -d ${TARGET_DIRECTORY} ]; then
        deleteInstallationDirectory
    fi
    mkdir $TARGET_DIRECTORY

    if [[ $? != 0 ]]; then
        log_error "Failed to create $TARGET_DIRECTORY directory" $?
        exit 1
    fi
}

copyDarwinDirectory(){
  createInstallationDirectory
  cp -r $SCRIPTPATH/darwin ${TARGET_DIRECTORY}/
  chmod -R 755 ${TARGET_DIRECTORY}/darwin/scripts
  chmod -R 755 ${TARGET_DIRECTORY}/darwin/Resources
  chmod 755 ${TARGET_DIRECTORY}/darwin/Distribution
}

copyBuildDirectory() {
    sed -i '' -e 's/__VERSION__/'${VERSION}'/g' ${TARGET_DIRECTORY}/darwin/scripts/postinstall
    sed -i '' -e 's/__PRODUCT__/'${PRODUCT}'/g' ${TARGET_DIRECTORY}/darwin/scripts/postinstall

    sed -i '' -e 's/__VERSION__/'${VERSION}'/g' ${TARGET_DIRECTORY}/darwin/Distribution
    sed -i '' -e 's/__PRODUCT__/'${PRODUCT}'/g' ${TARGET_DIRECTORY}/darwin/Distribution

    rm -rf ${TARGET_DIRECTORY}/darwinpkg
    mkdir -p ${TARGET_DIRECTORY}/darwinpkg

    #Copy opensearch-cli product to /Library/OpenSearchCLI
    mkdir -p ${TARGET_DIRECTORY}/darwinpkg/Library/${PRODUCT}
    cp -a $SCRIPTPATH/application/. ${TARGET_DIRECTORY}/darwinpkg/Library/${PRODUCT}
    chmod -R 755 ${TARGET_DIRECTORY}/darwinpkg/Library/${PRODUCT}

    rm -rf ${TARGET_DIRECTORY}/package
    mkdir -p ${TARGET_DIRECTORY}/package
    chmod -R 755 ${TARGET_DIRECTORY}/package

    rm -rf ${TARGET_DIRECTORY}/pkg
    mkdir -p ${TARGET_DIRECTORY}/pkg
    chmod -R 755 ${TARGET_DIRECTORY}/pkg
}

buildPackage() {
    log_info "Application installer package building started.(1/3)"
    pkgbuild --identifier org.opensearch.${PRODUCT} \
    --version ${VERSION} \
    --scripts ${TARGET_DIRECTORY}/darwin/scripts \
    --root ${TARGET_DIRECTORY}/darwinpkg \
    ${TARGET_DIRECTORY}/package/${PRODUCT}.pkg > /dev/null 2>&1
}

buildProduct() {
    log_info "Application installer product building started.(2/3)"
    productbuild --distribution ${TARGET_DIRECTORY}/darwin/Distribution \
    --resources ${TARGET_DIRECTORY}/darwin/Resources \
    --package-path ${TARGET_DIRECTORY}/package \
    ${TARGET_DIRECTORY}/pkg/$1 > /dev/null 2>&1
}

createInstaller() {
    log_info "Application installer generation process started.(3 Steps)"
    buildPackage
    buildProduct ${INSTALLER_NAME}.pkg
    log_info "Application installer generation steps finished."
}

#Main script
log_info "Installer generating process started."

copyDarwinDirectory
copyBuildDirectory
createInstaller

log_info "Installer generating process finished"
exit 0
