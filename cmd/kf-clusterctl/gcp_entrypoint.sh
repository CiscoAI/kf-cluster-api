#!/bin/bash

set -e

echo "${APPLICATION_CREDENTIALS}" | base64 -d > /tmp/account.json
gcloud -q auth activate-service-account --key-file=/tmp/account.json --user-output-enabled false
gcloud -q config set project "$PROJECT" --user-output-enabled false

export GOOGLE_APPLICATION_CREDENTIALS=/tmp/account.json
# kops create cluster - creates cluster spec and initializes state
kops create cluster ${CLUSTER_NAME} --zones ${ZONE} --state ${KOPS_STATE_STORE}/ --project=${PROJECT}
# kops update cluster - updates cluster spec, actual step that creates the cluster
kops update cluster ${CLUSTER_NAME} --yes
# Export created cluster kubeconfig
kops export kubecfg ${CLUSTER_NAME} --kubeconfig /mnt/volume/${CLUSTER_NAME}/kubeconfig

# Export Kubeconfig
export KUBECONFIG=/mnt/volume/${CLUSTER_NAME}/kubeconfig
kubectl get ns
# Install Kubeflow
mkdir -p /mnt/volume/${CLUSTER_NAME}/kf-app
cd /mnt/volume/${CLUSTER_NAME}/kf-app 
kfctl apply -V -f ${KF_CONFIG}
sleep 120
kubectl get po -A
