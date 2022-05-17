module github.com/jmcarp/prometheus-config

go 1.16

require (
	github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137
	github.com/aws/aws-sdk-go v1.43.31 // indirect
	github.com/cespare/xxhash/v2 v2.1.2
	github.com/go-kit/log v0.2.0
	github.com/go-openapi/strfmt v0.21.2
	github.com/gogo/protobuf v1.3.2
	github.com/golang/glog v1.0.0 // indirect
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/grafana/regexp v0.0.0-20220304095617-2e8d9baf4ac2
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/alertmanager v0.24.0
	github.com/prometheus/client_golang v1.12.1
	github.com/prometheus/common v0.34.0
	github.com/prometheus/common/sigv4 v0.1.0
	github.com/stretchr/testify v1.7.1 // indirect
	go.uber.org/atomic v1.9.0
	golang.org/x/net v0.0.0-20220412020605-290c469a71a5 // indirect
	golang.org/x/oauth2 v0.0.0-20220411215720-9780585627b5 // indirect
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad // indirect
	golang.org/x/tools v0.1.10
	golang.org/x/xerrors v0.0.0-20220411194840-2f41105eb62f // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220414192740-2d67ff6cf2b4 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v2 v2.4.0
)

// Exclude linodego v1.0.0 as it is no longer published on github.
exclude github.com/linode/linodego v1.0.0

// Exclude grpc v1.30.0 because of breaking changes. See #7621.
exclude (
	github.com/grpc-ecosystem/grpc-gateway v1.14.7
	google.golang.org/api v0.30.0
)

// Exclude pre-go-mod kubernetes tags, as they are older
// than v0.x releases but are picked when we update the dependencies.
exclude (
	k8s.io/client-go v1.4.0
	k8s.io/client-go v1.4.0+incompatible
	k8s.io/client-go v1.5.0
	k8s.io/client-go v1.5.0+incompatible
	k8s.io/client-go v1.5.1
	k8s.io/client-go v1.5.1+incompatible
	k8s.io/client-go v10.0.0+incompatible
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/client-go v2.0.0+incompatible
	k8s.io/client-go v2.0.0-alpha.1+incompatible
	k8s.io/client-go v3.0.0+incompatible
	k8s.io/client-go v3.0.0-beta.0+incompatible
	k8s.io/client-go v4.0.0+incompatible
	k8s.io/client-go v4.0.0-beta.0+incompatible
	k8s.io/client-go v5.0.0+incompatible
	k8s.io/client-go v5.0.1+incompatible
	k8s.io/client-go v6.0.0+incompatible
	k8s.io/client-go v7.0.0+incompatible
	k8s.io/client-go v8.0.0+incompatible
	k8s.io/client-go v9.0.0+incompatible
	k8s.io/client-go v9.0.0-invalid+incompatible
)

retract (
	v2.5.0+incompatible
	v2.5.0-rc.2+incompatible
	v2.5.0-rc.1+incompatible
	v2.5.0-rc.0+incompatible
	v2.4.3+incompatible
	v2.4.2+incompatible
	v2.4.1+incompatible
	v2.4.0+incompatible
	v2.4.0-rc.0+incompatible
	v2.3.2+incompatible
	v2.3.1+incompatible
	v2.3.0+incompatible
	v2.2.1+incompatible
	v2.2.0+incompatible
	v2.2.0-rc.1+incompatible
	v2.2.0-rc.0+incompatible
	v2.1.0+incompatible
	v2.0.0+incompatible
	v2.0.0-rc.3+incompatible
	v2.0.0-rc.2+incompatible
	v2.0.0-rc.1+incompatible
	v2.0.0-rc.0+incompatible
	v2.0.0-beta.5+incompatible
	v2.0.0-beta.4+incompatible
	v2.0.0-beta.3+incompatible
	v2.0.0-beta.2+incompatible
	v2.0.0-beta.1+incompatible
	v2.0.0-beta.0+incompatible
	v2.0.0-alpha.3+incompatible
	v2.0.0-alpha.2+incompatible
	v2.0.0-alpha.1+incompatible
	v2.0.0-alpha.0+incompatible
	v1.8.2
	v1.8.1
	v1.8.0
	v1.7.2
	v1.7.1
	v1.7.0
	v1.6.3
	v1.6.2
	v1.6.1
	v1.6.0
	v1.5.3
	v1.5.2
	v1.5.1
	v1.5.0
	v1.4.1
	v1.4.0
	v1.3.1
	v1.3.0
	v1.3.0-beta.0
	v1.2.3
	v1.2.2
	v1.2.1
	v1.2.0
	v1.1.3
	v1.1.2
	v1.1.1
	v1.1.0
	v1.0.2
	v1.0.1
	v1.0.0
	v1.0.0-rc.0
)
