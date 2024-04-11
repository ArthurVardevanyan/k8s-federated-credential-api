# Kubernetes Federated Credential Controller

```bash
go install goa.design/goa/v3/cmd/goa@v3
go get goa.design/goa/v3/http@v3.16.0

mkdir -p kubernetes-federated-credential-controller/design
cd kubernetes-federated-credential-controller
go mod init kubernetes-federated-credential-controller

~/go/bin/goa gen kubernetes-federated-credential-controller/design
~/go/bin/goa example kubernetes-federated-credential-controller/design

go build -C cmd/kfcc -o /tmp/kfcc

/tmp/kfcc

export JSON='{
  "jwt": "'$(kubectl create token default --duration=1h -n default)'",
  "namespace": "default",
  "serviceAccount": "default"
}'

curl "http://localhost:8088/exchangeToken" -X POST \
  -H "Content-type: application/json" \
  -H "Accept: application/json" \
  -d "${JSON}"
```
