#!/bin/bash

# Copyright 2019 Cortex Labs, Inc.
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


ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. >/dev/null && pwd)"

source $ROOT/dev/config/cortex.sh

eksctl utils write-kubeconfig --name=$CORTEX_CLUSTER --region=$CORTEX_REGION | grep -v "saved kubeconfig as" | grep -v "using region" || true

echo "Uninstalling Cortex ..."

kubectl delete --ignore-not-found=true namespace istio-system
kubectl delete --ignore-not-found=true namespace $CORTEX_NAMESPACE

echo "✓ Uninstalled Cortex"
