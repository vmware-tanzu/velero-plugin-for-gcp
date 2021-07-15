module github.com/vmware-tanzu/velero-plugin-for-gcp

go 1.13

require (
	cloud.google.com/go/storage v1.8.0
	github.com/hashicorp/go-hclog v0.9.2 // indirect
	github.com/hashicorp/go-plugin v1.0.1-0.20190610192547-a1bc61569a26 // indirect
	github.com/hashicorp/yamux v0.0.0-20190923154419-df201c70410d // indirect
	github.com/pkg/errors v0.9.1
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.5.1
	github.com/vmware-tanzu/velero v1.6.1
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	google.golang.org/api v0.24.0
	k8s.io/api v0.19.12
	k8s.io/apimachinery v0.19.12
)

replace github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.2
