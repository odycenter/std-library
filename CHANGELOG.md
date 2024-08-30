## Change log

### 1.8.5 (08/07/2024 - 08/21/2024)

* executor: add running/free/waiting metrics support
* log:
    - `remove unnecessary level`
    - `integration with slog`
    - `dynamic log level change by /_sys/log`
* remove ipx support
* db: config support
* http: http: server is enabled by default port 18080, even when not configured

### 1.7.49 (05/30/2024 - 07/31/2024)

* api: /_sys/api to get http api definition, and use sys.api.allowCIDR to restrict access
* property: multi line supported
* metric: metrics server supported
* redis/mongo: multi config supported
* redis/cache: pool size supported
* grpc:
    - `client retry 5 times when server response error UNAVAILABLE`
    - `client slow grpc log`
    - `client support migration option to switch grpc.NewClient for create grpc.ClientConn`
    - `ReadinessProbe when register client by default`
    - `maxConnections supported`
* mysql: tweak default ConnMaxIdleTime to 30 minutes, ConnMaxLifetime = ConnMaxIdleTime * 2

### 1.7.28 (05/29/2024 - 05/30/2024)

* mongo: execution elapsed monitor

### 1.7.27 (05/29/2024 - 05/29/2024)

* scheduler: /_sys/job/:name to trigger job

### 1.7.26 (05/23/2024 - 05/28/2024)

* property: /_sys/cache to get cache info
* scheduler: /_sys/job to get info of all jobs

### 1.7.24 (05/21/2024 - 05/23/2024)

* app: added external dependency checking before start, currently only check kafka and redis/cache to be ready
* property: /_sys/property to print properties and env var

### 1.7.22 (05/08/2024 - 05/16/2024)

* mysql: ConnMaxIdletime 1 hour, ConnMaxLifetime 2 hours with default
* redis: db selecting supported
* scheduler: disallow concurrent job support
* pyroscope: disable local startup by default
* kafka: if runtime.NumCPU() * 4 > 4, max 4 consumer by default
* module: mongo support
* mongo: default timeout 120s
* ResponseWriter: fix Hijacker type for websocket upgrade
* property manager: fix LoadProperties issue when values include '='

### 1.7.11 (05/02/2024 - 05/07/2024)

* mysql: ConnMaxIdletime, default 2 hours
* mongo: MaxConnIdleTime, default 30 minutes
* module: redis/cache support

### 1.7.10 (04/01/2024 - 04/29/2024)

* fix kafka producer memory leak issue
* mysql: driver to 1.8.1
* graceful shutdown: http server, grpc server, scheduler and kafka message listener support graceful shutdown now
* package change!
    - `async` to `app/async`
    - `conf` to `app/conf`
    - `property` to `app/property`
    - `scheduler` to `app/scheduler`
    - `web` to `app/web`
    - `actionlog` to `app/log`
    - `grpc` to `app/web/grpc`

### 1.6.13 (03/21/2024 - 03/28/2024)

* kafka: revert!! consumer watch partition by default 60s interval
* kafka: producer force flush when local queue full error
* ipx: remove GeoLite support
* slice: add slice to map util
* grpc: server default with prometheus interceptor
* mongo: add AppName for identify client

### 1.6.10 (03/19/2024 - 03/20/2024)

* web: header ref-id & client support
* grpc: context util support

### 1.6.9 (03/18/2024 - 03/18/2024)

* ipx: add new driver ipdb
* conf: improve conf support

### 1.6.7 (03/15/2024 - 03/15/2024)

* ipx: update default driver

### 1.6.6 (03/14/2024 - 03/14/2024)

* improve conf manager support

### 1.6.4 (03/06/2024 - 03/13/2024)

* kafka: consumer watch partition by default 60s interval
* action log: tweak log context usage
* conf manager support

### 1.5.9-5 (02/21/2024 - 02/29/2024)

* web: beego filter force parseForm before read
* tg: fix SimpleSendMsg parameter not append issue
* property manager support

### 1.5.7 (02/20/2024 - 02/20/2024)

* kafka: support Compression switch, use EnableCompression to enable, default: false

### 1.5.6 (02/15/2024 - 02/19/2024)

* kafka: confluent-kafka-go producer support, use CliV2 for confluent-kafka-go client

### 1.5.3 (02/05/2024 - 02/06/2024)

* grpc: add default timeout support

> default timeout: 120s, use grpc.EnableDefaultTimeout to enable, use grpc.EnableTimeout to set specific timeout.

### 1.5.1 (02/01/2024 - 02/02/2024)

* kafka: add pool size of consumer, default: (Processors * 4)
* ipx: ip geolocation util support
* db: force live check for db connection

### 1.4.1 (01/29/2024 - 01/31/2024) !!! breaking changes

* async task support
  > tweak package, use async.Executor instead to executor.Executor
  > task: support WaitGroup

### 1.3.2 (01/25/2024 - 01/26/2024)

* scheduler supported
* reflects util: for retrieve StructName and FunctionName