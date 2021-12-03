#!/bin/bash
# SPDX-License-Identifier: Apache-2.0
#
# The OpenSearch Contributors require contributions made to
# this file be licensed under the Apache-2.0 license or a
# compatible open source license.
#
# Modifications Copyright OpenSearch Contributors. See
# GitHub history for details.

# This is intended to be run from the  root directory. `release/build.sh`
set -e

GIT_ROOT=`git rev-parse --show-toplevel`

cd $GIT_ROOT # We need to start from repository root

#####################################################
#      Generating Artifacts opensearch-cli                      #
#                                                   #
#####################################################


echo 'Generating artifacts'
goreleaser --snapshot --skip-publish --rm-dist

# goreleaser generates folder and binary too. Remove unwanted files to keep only
# relevant files inside dist folder
rm -rf dist/opensearch-cli_*
rm -rf dist/config.yaml

ls -l dist/

