#!/bin/bash
alias kubectl="kubecolor"
shopt -s expand_aliases

# Also works for OpenShift CLI

jwt-decode() {
  jq -R 'split(".") |.[0:2] | map(@base64d) | map(fromjson)' <<<$1
}

########################
# include the magic
########################
. /home/arthur/demo-magic.sh -n

# hide the evidence
clear

# Put your stuff here

# this command is typed and executed
export KUBECONFIG=~/.kube/okd

pe "kubectl get nodes"

export KUBECONFIG=~/.kube/microshift

pe "kubectl get nodes"

wait

pe "kubectl get serviceaccount -n default argocd -o yaml | kubectl neat | yq"

wait

export KUBECONFIG=~/.kube/microshift

pe "kubectl get nodes"

pe "export JSON='{
  \"namespace\": \"default\",
  \"serviceAccountName\": \"argocd\"
}'"

wait

p 'jwt-decode curl --silent --header "Authorization: Bearer $(kubectl create token argocd-creds --duration=1h -n argocd)"\
  "https://kfca.microshift.arthurvardevanyan.com/exchangeToken" -X POST \
  -H "Content-type: application/json" \
  -H "Accept: application/json" \
  -d "${JSON}" | jq --raw-output .status.token'

jwt-decode $(curl --silent --header "Authorization: Bearer $(kubectl create token argocd-creds --duration=1h -n argocd)" \
  "https://kfca.microshift.arthurvardevanyan.com/exchangeToken" -X POST \
  -H "Content-type: application/json" \
  -H "Accept: application/json" \
  -d "${JSON}" | jq --raw-output '.status.token')
