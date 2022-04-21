# Sample Functional API test project (Golang)

## Steps

### 1. run InfluxDB2, Grafana with docker compose file

   docker-compose -f docker-compose.yml up -d

### 2. update hosts file if run on localmachine

   vi /etc/hosts

- 127.0.0.1    influxdb2
- 127.0.0.1    grafana

### 3. get InfluxDB2 token and set into ./config/main.yaml > influxdb2.token

- InfluxDB2: <http://influxdb2:8086> admin / admin123
- InfluxDB2: Data -> API Tokens (tab) -> Click "admin's Token" -> Copy the token to clipboard: "$COPIED_TOKEN"

### 4. build docker image and run

1. docker build -t functest:golang .
2. docker run -it --name functest --network=sample-functest_functest functest:golang ./buildtest.sh env.prod "$COPIED_TOKEN"

### 5. check data in InfluxDB2

- InfluxDB2: check test data be generated in bucket = FuncTestBucket & measurement = FuncTest

### 6. config Grafana + InfluxDB2 (Optinal)

- Grafana: <http://grafana:3000> admin / admin -> skip change password
- config -> data source -> add data source
- select InfluxDB as template
- select Flux as language
- close basic auth
- Add URL: <http://influxdb2:8086> -> Org = FuncTestOrg -> Token = "$COPIED_TOKEN" -> Bucket = FuncTestBucket -> save & test
- create new dashboard and select InfluxDB data source
