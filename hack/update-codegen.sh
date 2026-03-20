#!/usr/bin/env bash

# Copyright 2025.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
CODEGEN_PKG=${CODEGEN_PKG:-$(cd "${SCRIPT_ROOT}"; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}

MODULE="sigs.k8s.io/rbgs/api"
APIS_PKG="${MODULE}/apis"
CLIENT_PKG="${MODULE}/client-go"
GROUPS="workloads:v1alpha1,v1alpha2"
OUTPUT_BASE="$(dirname "${BASH_SOURCE[0]}")/../.."

# Run deepcopy-gen for API types
echo "Running deepcopy-gen..."
"${GOPATH}/bin/controller-gen" object:headerFile="${SCRIPT_ROOT}/hack/boilerplate.go.txt" \
  paths="${APIS_PKG}/..." \
  output:object:dir="${SCRIPT_ROOT}/apis"

# Run client-gen, informer-gen, lister-gen
echo "Running k8s code generators..."
bash "${CODEGEN_PKG}/generate-groups.sh" \
  "client,lister,informer" \
  "${CLIENT_PKG}" \
  "${APIS_PKG}" \
  "${GROUPS}" \
  --output-base "${OUTPUT_BASE}" \
  --go-header-file "${SCRIPT_ROOT}/hack/boilerplate.go.txt"

echo "Code generation complete."
