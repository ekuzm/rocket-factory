module github.com/ekuzm/rocket-factory/inventory

go 1.25.6

require (
	github.com/brianvoe/gofakeit/v7 v7.14.0
	github.com/caarlos0/env/v11 v11.4.0
	github.com/ekuzm/rocket-factory/shared v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	github.com/samber/lo v1.52.0
	github.com/stretchr/testify v1.11.1
	go.mongodb.org/mongo-driver v1.17.9
	google.golang.org/grpc v1.78.0
	google.golang.org/protobuf v1.36.10
)

replace github.com/ekuzm/rocket-factory/shared => ../shared

replace github.com/ekuzm/rocket-factory/platform => ../platform

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/klauspost/compress v1.18.1 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	golang.org/x/crypto v0.47.0 // indirect
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251029180050-ab9386a59fda // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
