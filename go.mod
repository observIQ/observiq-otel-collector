module github.com/observiq/observiq-collector

go 1.17

require (
	github.com/Azure/azure-event-hubs-go/v3 v3.3.13
	github.com/aws/aws-sdk-go v1.41.19
	github.com/hashicorp/go-multierror v1.1.0
	github.com/jpillora/backoff v1.0.0
	github.com/json-iterator/go v1.1.12
	github.com/mitchellh/mapstructure v1.4.2
	github.com/observiq/goflow/v3 v3.4.4
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/fileexporter v0.39.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/googlecloudexporter v0.39.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/observiqexporter v0.39.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/bearertokenauthextension v0.39.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/healthcheckextension v0.39.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/oidcauthextension v0.39.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/pprofextension v0.39.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/storage v0.39.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/stanza v0.39.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor v0.39.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbyattrsprocessor v0.39.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor v0.39.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricstransformprocessor v0.39.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filelogreceiver v0.39.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/jmxreceiver v0.39.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8sclusterreceiver v0.39.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kubeletstatsreceiver v0.39.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mysqlreceiver v0.39.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/syslogreceiver v0.39.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/tcplogreceiver v0.39.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/udplogreceiver v0.39.0
	github.com/open-telemetry/opentelemetry-log-collection v0.22.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/testcontainers/testcontainers-go v0.11.1
	go.etcd.io/bbolt v1.3.6
	go.opentelemetry.io/collector v0.39.0
	go.opentelemetry.io/collector/model v0.39.0
	go.uber.org/multierr v1.7.0
	go.uber.org/zap v1.19.1
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/apimachinery v0.22.4
	k8s.io/client-go v0.22.4
)

require (
	github.com/gorilla/mux v1.8.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/probabilisticsamplerprocessor v0.39.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourceprocessor v0.39.0
)

require (
	cloud.google.com/go v0.97.0 // indirect
	cloud.google.com/go/monitoring v1.1.0 // indirect
	cloud.google.com/go/trace v1.0.0 // indirect
	contrib.go.opencensus.io/exporter/prometheus v0.4.0 // indirect
	contrib.go.opencensus.io/exporter/stackdriver v0.13.10 // indirect
	github.com/Azure/azure-amqp-common-go/v3 v3.0.1 // indirect
	github.com/Azure/azure-sdk-for-go v55.2.0+incompatible // indirect
	github.com/Azure/go-amqp v0.13.12 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest v0.11.19 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.14 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/to v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v1.0.0 // indirect
	github.com/Microsoft/go-winio v0.4.17 // indirect
	github.com/Microsoft/hcsshim v0.8.21 // indirect
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/antonmedv/expr v1.9.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bmatcuk/doublestar/v3 v3.0.0 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/cenkalti/backoff/v4 v4.1.2 // indirect
	github.com/census-instrumentation/opencensus-proto v0.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/cncf/udpa/go v0.0.0-20210930031921-04548b0d99d4 // indirect
	github.com/cncf/xds/go v0.0.0-20211011173535-cb28da3451f1 // indirect
	github.com/containerd/cgroups v1.0.1 // indirect
	github.com/containerd/containerd v1.5.7 // indirect
	github.com/coreos/go-oidc v2.2.1+incompatible // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/devigned/tab v0.1.1 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v20.10.10+incompatible // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/envoyproxy/go-control-plane v0.9.10-0.20210907150352-cf90f659a021 // indirect
	github.com/envoyproxy/protoc-gen-validate v0.1.0 // indirect
	github.com/felixge/httpsnoop v1.0.2 // indirect
	github.com/form3tech-oss/jwt-go v3.2.3+incompatible // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/go-kit/log v0.2.0 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-logr/logr v0.4.0 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/gax-go/v2 v2.1.1 // indirect
	github.com/googleapis/gnostic v0.5.5 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/iancoleman/strcase v0.2.0 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jaegertracing/jaeger v1.28.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/knadh/koanf v1.3.2 // indirect
	github.com/libp2p/go-reuseport v0.0.1 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/moby/sys/mount v0.2.0 // indirect
	github.com/moby/sys/mountinfo v0.4.1 // indirect
	github.com/moby/term v0.0.0-20210610120745-9d4ed1856297 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/mostynb/go-grpc-compression v1.1.14 // indirect
	github.com/observiq/ctimefmt v1.0.0 // indirect
	github.com/observiq/go-syslog/v3 v3.0.2 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/common v0.39.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal v0.39.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/k8sconfig v0.39.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/kubelet v0.39.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/sharedcomponent v0.39.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/experimentalmetricmetadata v0.39.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/opencensus v0.39.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/opencontainers/runc v1.0.2 // indirect
	github.com/openshift/api v0.0.0-20210521075222-e273a339932a // indirect
	github.com/openshift/client-go v0.0.0-20210521082421-73d9475a9142 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/pquerna/cachecontrol v0.1.0 // indirect
	github.com/prometheus/client_golang v1.11.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/prometheus/statsd_exporter v0.21.0 // indirect
	github.com/rs/cors v1.8.0 // indirect
	github.com/shirou/gopsutil/v3 v3.21.10 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/cobra v1.2.1 // indirect
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/tklauser/go-sysconf v0.3.9 // indirect
	github.com/tklauser/numcpus v0.3.0 // indirect
	go.opencensus.io v0.23.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.26.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.26.1 // indirect
	go.opentelemetry.io/contrib/zpages v0.26.1 // indirect
	go.opentelemetry.io/otel v1.1.0 // indirect
	go.opentelemetry.io/otel/exporters/prometheus v0.24.0 // indirect
	go.opentelemetry.io/otel/internal/metric v0.24.0 // indirect
	go.opentelemetry.io/otel/metric v0.24.0 // indirect
	go.opentelemetry.io/otel/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel/sdk/export/metric v0.24.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v0.24.0 // indirect
	go.opentelemetry.io/otel/trace v1.1.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	golang.org/x/crypto v0.0.0-20210920023735-84f357641f63 // indirect
	golang.org/x/net v0.0.0-20210917221730-978cfadd31cf // indirect
	golang.org/x/oauth2 v0.0.0-20211005180243-6b3c2da341f1 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20211025201205-69cdffdb9359 // indirect
	golang.org/x/term v0.0.0-20210220032956-6a3ed077a48d // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
	gonum.org/v1/gonum v0.9.3 // indirect
	google.golang.org/api v0.60.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20211021150943-2b146023228c // indirect
	google.golang.org/grpc v1.42.0 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	k8s.io/api v0.22.4 // indirect
	k8s.io/klog/v2 v2.9.0 // indirect
	k8s.io/kubelet v0.22.3 // indirect
	k8s.io/utils v0.0.0-20210819203725-bdf08cb9a70a // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.1.2 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace (
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer v0.0.0-00010101000000-000000000000 => github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer v0.39.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/storage v0.0.0-00010101000000-000000000000 => github.com/open-telemetry/opentelemetry-collector-contrib/extension/storage v0.39.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/k8sconfig v0.0.0-00010101000000-000000000000 => github.com/open-telemetry/opentelemetry-collector-contrib/internal/k8sconfig v0.39.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/stanza v0.0.0-00010101000000-000000000000 => github.com/open-telemetry/opentelemetry-collector-contrib/internal/stanza v0.39.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/stanza v0.24.1-0.20210408210148-736647af91e1 => github.com/open-telemetry/opentelemetry-collector-contrib/internal/stanza v0.39.0
)
