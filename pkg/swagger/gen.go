package swagger

//go:generate rm -rf server/models server/restapi
//go:generate mkdir -p server
//go:generate swagger generate server --quiet --target server --name filer-gateway --spec swagger.yaml --exclude-main -P models.Principle
//go:generate rm -rf client/models client/client
//go:generate mkdir -p client
//go:generate swagger generate client --quiet --target client --name filer-gateway --spec swagger.yaml
