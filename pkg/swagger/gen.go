package swagger

//go:generate rm -rf server/models server/restapi
//go:generate mkdir -p server
//go:generate swagger generate server --quiet --target server --name filer-gateway --spec swagger.yaml --exclude-main -P models.Principle
