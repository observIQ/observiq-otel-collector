module github.com/observiq/observiq-otel-collector

go 1.17

require (
	github.com/GoogleCloudPlatform/opentelemetry-operations-collector v0.0.3-0.20220711143229-08f2752ed367
	github.com/google/uuid v1.3.0
	github.com/observiq/observiq-otel-collector/exporter/googlecloudexporter v1.3.0
	github.com/observiq/observiq-otel-collector/packagestate v0.0.0
	github.com/observiq/observiq-otel-collector/processor/resourceattributetransposerprocessor v1.3.0
	github.com/observiq/observiq-otel-collector/receiver/pluginreceiver v1.3.0
	github.com/open-telemetry/opamp-go v0.2.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/alibabacloudlogserviceexporter v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awscloudwatchlogsexporter v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awsemfexporter v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awskinesisexporter v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awsprometheusremotewriteexporter v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awsxrayexporter v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/azuremonitorexporter v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/carbonexporter v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/elasticsearchexporter v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/f5cloudexporter v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/fileexporter v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/googlecloudpubsubexporter v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/googlemanagedprometheusexporter v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/influxdbexporter v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/jaegerexporter v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/kafkaexporter v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/loadbalancingexporter v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/logzioexporter v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/lokiexporter v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/observiqexporter v0.55.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/opencensusexporter v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusexporter v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/sapmexporter v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/splunkhecexporter v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/zipkinexporter v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/bearertokenauthextension v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/healthcheckextension v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/oidcauthextension v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/pprofextension v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/storage v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/cumulativetodeltaprocessor v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/deltatorateprocessor v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/filterprocessor v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbyattrsprocessor v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbytraceprocessor v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/logstransformprocessor v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricsgenerationprocessor v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricstransformprocessor v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/probabilisticsamplerprocessor v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourceprocessor v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/routingprocessor v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/spanmetricsprocessor v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/spanprocessor v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/tailsamplingprocessor v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/activedirectorydsreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/aerospikereceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/apachereceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awscontainerinsightreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsecscontainermetricsreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsfirehosereceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsxrayreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/bigipreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/carbonreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/cloudfoundryreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/collectdreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/couchdbreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/dockerstatsreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/dotnetdiagnosticsreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/elasticsearchreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filelogreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/flinkmetricsreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/fluentforwardreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/googlecloudpubsubreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/googlecloudspannerreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/iisreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/influxdbreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/jaegerreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/jmxreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/journaldreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8sclusterreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8seventsreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kafkametricsreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kafkareceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kubeletstatsreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/memcachedreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mongodbatlasreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mongodbreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mysqlreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/nginxreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/opencensusreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/podmanreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/postgresqlreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/rabbitmqreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/redisreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/riakreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sapmreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/simpleprometheusreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sqlserverreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/statsdreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/syslogreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/tcplogreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/udplogreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/vcenterreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/windowseventlogreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/windowsperfcountersreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/zipkinreceiver v0.56.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/zookeeperreceiver v0.56.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.8.0
	go.opentelemetry.io/collector v0.56.0
	go.uber.org/multierr v1.8.0
	go.uber.org/zap v1.21.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	cloud.google.com/go v0.102.1 // indirect
	cloud.google.com/go/compute v1.7.0 // indirect
	cloud.google.com/go/iam v0.3.0 // indirect
	cloud.google.com/go/logging v1.4.2 // indirect
	cloud.google.com/go/monitoring v1.5.0 // indirect
	cloud.google.com/go/pubsub v1.23.1 // indirect
	cloud.google.com/go/spanner v1.34.1 // indirect
	cloud.google.com/go/trace v1.2.0 // indirect
	code.cloudfoundry.org/clock v1.0.0 // indirect
	code.cloudfoundry.org/go-diodes v0.0.0-20211115184647-b584dd5df32c // indirect
	code.cloudfoundry.org/go-loggregator v7.4.0+incompatible // indirect
	code.cloudfoundry.org/rfc5424 v0.0.0-20201103192249-000122071b78 // indirect
	contrib.go.opencensus.io/exporter/prometheus v0.4.1 // indirect
	contrib.go.opencensus.io/exporter/stackdriver v0.13.13 // indirect
	github.com/Azure/azure-sdk-for-go v65.0.0+incompatible // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest v0.11.27 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.20 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/to v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp v0.32.3 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/collector v0.32.3 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/collector/googlemanagedprometheus v0.32.3 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v1.8.3 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping v0.32.3 // indirect
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/ReneKroon/ttlcache/v2 v2.11.0 // indirect
	github.com/Shopify/sarama v1.34.1 // indirect
	github.com/Showmax/go-fqdn v1.0.0 // indirect
	github.com/VividCortex/gohistogram v1.0.0 // indirect
	github.com/aerospike/aerospike-client-go/v5 v5.9.0 // indirect
	github.com/alecthomas/participle/v2 v2.0.0-alpha9 // indirect
	github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137 // indirect
	github.com/aliyun/aliyun-log-go-sdk v0.1.37 // indirect
	github.com/andybalholm/brotli v1.0.1 // indirect
	github.com/antonmedv/expr v1.9.0 // indirect
	github.com/apache/thrift v0.16.0 // indirect
	github.com/armon/go-metrics v0.3.10 // indirect
	github.com/aws/aws-sdk-go v1.44.57 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/bmatcuk/doublestar/v3 v3.0.0 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/cenkalti/backoff/v4 v4.1.3 // indirect
	github.com/census-instrumentation/opencensus-proto v0.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/checkpoint-restore/go-criu/v5 v5.3.0 // indirect
	github.com/cilium/ebpf v0.7.0 // indirect
	github.com/cloudfoundry-incubator/uaago v0.0.0-20190307164349-8136b7bbe76e // indirect
	github.com/cncf/udpa/go v0.0.0-20210930031921-04548b0d99d4 // indirect
	github.com/cncf/xds/go v0.0.0-20220314180256-7f1daf1720fc // indirect
	github.com/containerd/console v1.0.3 // indirect
	github.com/containerd/ttrpc v1.1.0 // indirect
	github.com/coreos/go-oidc v2.2.1+incompatible // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/cyphar/filepath-securejoin v0.2.3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dennwc/varint v1.0.0 // indirect
	github.com/digitalocean/godo v1.81.0 // indirect
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/docker v20.10.17+incompatible // indirect
	github.com/docker/go-connections v0.4.1-0.20210727194412-58542c764a11 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/dsnet/compress v0.0.2-0.20210315054119-f66993602bf5 // indirect
	github.com/eapache/go-resiliency v1.2.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20180814174437-776d5712da21 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/elastic/go-elasticsearch/v7 v7.17.1 // indirect
	github.com/elastic/go-structform v0.0.9 // indirect
	github.com/emicklei/go-restful v2.9.5+incompatible // indirect
	github.com/envoyproxy/go-control-plane v0.10.3 // indirect
	github.com/envoyproxy/protoc-gen-validate v0.6.7 // indirect
	github.com/euank/go-kmsg-parser v2.0.0+incompatible // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/go-kit/kit v0.12.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.20.0 // indirect
	github.com/go-openapi/swag v0.21.1 // indirect
	github.com/go-redis/redis/v7 v7.4.1 // indirect
	github.com/go-resty/resty/v2 v2.1.1-0.20191201195748-d7b97669fe48 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/go-zookeeper/zk v1.0.2 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/godbus/dbus/v5 v5.0.6 // indirect
	github.com/gofrs/uuid v4.2.0+incompatible // indirect
	github.com/gogo/googleapis v1.4.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.2.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/cadvisor v0.44.1 // indirect
	github.com/google/gnostic v0.5.7-v3refs // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.1.0 // indirect
	github.com/googleapis/gax-go/v2 v2.4.0 // indirect
	github.com/gophercloud/gophercloud v0.25.0 // indirect
	github.com/gorilla/handlers v1.5.1 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/grafana/regexp v0.0.0-20220304095617-2e8d9baf4ac2 // indirect
	github.com/grobie/gomemcache v0.0.0-20180201122607-1f779c573665 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/hashicorp/consul/api v1.13.1 // indirect
	github.com/hashicorp/cronexpr v1.1.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.2.1 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.1 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/go-uuid v1.0.2 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/nomad/api v0.0.0-20220718181929-9457a13c7ccf // indirect
	github.com/hashicorp/serf v0.9.7 // indirect
	github.com/hetznercloud/hcloud-go v1.35.0 // indirect
	github.com/iancoleman/strcase v0.2.0 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/influxdata/go-syslog/v3 v3.0.1-0.20210608084020-ac565dc76ba6 // indirect
	github.com/influxdata/influxdb-observability/common v0.2.21-0.20220517160208-05f925d616de // indirect
	github.com/influxdata/influxdb-observability/influx2otel v0.2.21-0.20220517160208-05f925d616de // indirect
	github.com/influxdata/influxdb-observability/otel2influx v0.2.21-0.20220517160208-05f925d616de // indirect
	github.com/influxdata/line-protocol/v2 v2.2.1 // indirect
	github.com/ionos-cloud/sdk-go/v6 v6.1.0 // indirect
	github.com/jaegertracing/jaeger v1.36.0 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.0.0 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.2 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/karrick/godirwalk v1.16.1 // indirect
	github.com/klauspost/pgzip v1.2.5 // indirect
	github.com/knadh/koanf v1.4.2 // indirect
	github.com/kolo/xmlrpc v0.0.0-20201022064351-38db28db192b // indirect
	github.com/leoluk/perflib_exporter v0.1.0 // indirect
	github.com/lib/pq v1.10.6 // indirect
	github.com/linode/linodego v1.8.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/microsoft/ApplicationInsights-Go v0.4.4 // indirect
	github.com/miekg/dns v1.1.50 // indirect
	github.com/mindprince/gonvml v0.0.0-20190828220739-9ebdce4bb989 // indirect
	github.com/mistifyio/go-zfs v2.1.2-0.20190413222219-f784269be439+incompatible // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/hashstructure v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/moby/sys/mountinfo v0.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mongodb-forks/digest v1.0.4 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/mostynb/go-grpc-compression v1.1.16 // indirect
	github.com/mrunalp/fileutils v0.5.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f // indirect
	github.com/nginxinc/nginx-prometheus-exporter v0.8.1-0.20201110005315-f5a5f8086c19 // indirect
	github.com/nwaples/rardecode v1.1.0 // indirect
	github.com/observiq/ctimefmt v1.0.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/googlecloudexporter v0.56.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/awsutil v0.56.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/containerinsight v0.56.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/cwlogs v0.56.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/ecsutil v0.56.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/k8s v0.56.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/metrics v0.56.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/proxy v0.56.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/xray v0.56.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/common v0.56.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal v0.56.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/docker v0.56.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/k8sconfig v0.56.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/kubelet v0.56.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/metadataproviders v0.56.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/sharedcomponent v0.56.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/splunk v0.56.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/batchperresourceattr v0.56.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/batchpersignal v0.56.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/experimentalmetricmetadata v0.56.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/resourcetotelemetry v0.56.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza v0.56.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/telemetryquerylanguage v0.56.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/jaeger v0.56.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/opencensus v0.56.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/prometheus v0.56.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/prometheusremotewrite v0.56.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/signalfx v0.56.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/zipkin v0.56.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/winperfcounters v0.56.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.3-0.20211202183452-c5a74bcca799 // indirect
	github.com/opencontainers/runc v1.1.2 // indirect
	github.com/opencontainers/runtime-spec v1.0.3-0.20210326190908-1c3f411f0417 // indirect
	github.com/opencontainers/selinux v1.10.1 // indirect
	github.com/openlyinc/pointy v1.1.2 // indirect
	github.com/openshift/api v0.0.0-20210521075222-e273a339932a // indirect
	github.com/openshift/client-go v0.0.0-20210521082421-73d9475a9142 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/openzipkin/zipkin-go v0.4.0 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.1 // indirect
	github.com/philhofer/fwd v1.1.1 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pierrec/lz4/v4 v4.1.14 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/pquerna/cachecontrol v0.1.0 // indirect
	github.com/prometheus/client_golang v1.12.2 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/common/sigv4 v0.1.0 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/prometheus/prometheus v0.37.0 // indirect
	github.com/prometheus/statsd_exporter v0.21.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/rs/cors v1.8.2 // indirect
	github.com/scaleway/scaleway-sdk-go v1.0.0-beta.9 // indirect
	github.com/seccomp/libseccomp-golang v0.9.2-0.20210429002308-3879420cc921 // indirect
	github.com/signalfx/com_signalfx_metrics_protobuf v0.0.3 // indirect
	github.com/signalfx/gohistogram v0.0.0-20160107210732-1ccfd2ff5083 // indirect
	github.com/signalfx/golib/v3 v3.3.13 // indirect
	github.com/signalfx/sapm-proto v0.11.0 // indirect
	github.com/signalfx/signalfx-agent/pkg/apm v0.0.0-20201202163743-65b4fa925fc8 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/spf13/afero v1.8.2 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/viper v1.12.0 // indirect
	github.com/stretchr/objx v0.4.0 // indirect
	github.com/subosito/gotenv v1.3.0 // indirect
	github.com/syndtr/gocapability v0.0.0-20200815063812-42c35b437635 // indirect
	github.com/tidwall/gjson v1.10.2 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/tidwall/tinylru v1.1.0 // indirect
	github.com/tidwall/wal v1.1.7 // indirect
	github.com/tinylib/msgp v1.1.6 // indirect
	github.com/tklauser/go-sysconf v0.3.10 // indirect
	github.com/tklauser/numcpus v0.4.0 // indirect
	github.com/uber/jaeger-client-go v2.30.0+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	github.com/ulikunitz/xz v0.5.9 // indirect
	github.com/vishvananda/netlink v1.1.1-0.20210330154013-f5de75959ad5 // indirect
	github.com/vishvananda/netns v0.0.0-20210104183010-2eb08e3e575f // indirect
	github.com/vmware/govmomi v0.28.0 // indirect
	github.com/vultr/govultr/v2 v2.17.2 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.1 // indirect
	github.com/xdg-go/stringprep v1.0.3 // indirect
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	github.com/yuin/gopher-lua v0.0.0-20220504180219-658193537a64 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	go.etcd.io/bbolt v1.3.6 // indirect
	go.mongodb.org/atlas v0.16.0 // indirect
	go.mongodb.org/mongo-driver v1.10.0 // indirect
	go.opencensus.io v0.23.0 // indirect
	go.opentelemetry.io/collector/pdata v0.56.0 // indirect
	go.opentelemetry.io/collector/semconv v0.56.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.33.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.33.0 // indirect
	go.opentelemetry.io/contrib/zpages v0.33.0 // indirect
	go.opentelemetry.io/otel v1.8.0 // indirect
	go.opentelemetry.io/otel/exporters/prometheus v0.31.0 // indirect
	go.opentelemetry.io/otel/metric v0.31.0 // indirect
	go.opentelemetry.io/otel/sdk v1.8.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v0.31.0 // indirect
	go.opentelemetry.io/otel/trace v1.8.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/goleak v1.1.12 // indirect
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d // indirect
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4 // indirect
	golang.org/x/net v0.0.0-20220630215102-69896b714898 // indirect
	golang.org/x/oauth2 v0.0.0-20220628200809-02e64fa58f26 // indirect
	golang.org/x/sync v0.0.0-20220601150217-0de741cfad7f // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20220609170525-579cf78fd858 // indirect
	golang.org/x/tools v0.1.11 // indirect
	golang.org/x/xerrors v0.0.0-20220609144429-65e65417b02f // indirect
	gonum.org/v1/gonum v0.11.0 // indirect
	google.golang.org/api v0.88.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220630174209-ad1d48641aa7 // indirect
	google.golang.org/grpc v1.48.0 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.66.4 // indirect
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/api v0.24.3 // indirect
	k8s.io/apimachinery v0.24.3 // indirect
	k8s.io/client-go v0.24.3 // indirect
	k8s.io/klog v1.0.0 // indirect
	k8s.io/klog/v2 v2.70.0 // indirect
	k8s.io/kube-openapi v0.0.0-20220328201542-3ee0da9b0b42 // indirect
	k8s.io/kubelet v0.24.0 // indirect
	k8s.io/utils v0.0.0-20220210201930-3a6ce19ff2f9 // indirect
	sigs.k8s.io/json v0.0.0-20211208200746-9f7c6b3444d2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.1 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)

require (
	github.com/containerd/containerd v1.6.6 // indirect
	github.com/klauspost/compress v1.15.8 // indirect
	github.com/mholt/archiver/v3 v3.5.1
	github.com/shirou/gopsutil/v3 v3.22.6
	github.com/spf13/cobra v1.5.0 // indirect
	golang.org/x/sys v0.0.0-20220704084225-05e143d24a9e
)

replace github.com/observiq/observiq-otel-collector/processor/resourceattributetransposerprocessor => ./processor/resourceattributetransposerprocessor

replace github.com/observiq/observiq-otel-collector/receiver/pluginreceiver => ./receiver/pluginreceiver

replace github.com/observiq/observiq-otel-collector/exporter/googlecloudexporter => ./exporter/googlecloudexporter

replace github.com/observiq/observiq-otel-collector/packagestate => ./packagestate
