# Sample Functional API test project (Golang)

## Steps

### 1. update hosts file if run on local machine

   vi /etc/hosts

- 127.0.0.1    influxdb2
- 127.0.0.1    grafana

### 2. run InfluxDB2, Grafana with docker compose file

   docker-compose -f docker-compose.yml up -d

### 3. setup InfluxDB2 with api if not set environment propertyies in step 2

~~~bash
curl -v POST \
  http://influxdb2:8086/api/v2/setup \
  --header 'Content-type: application/json' \
  --data '{
  "username": "admin",
  "password": "admin123",
  "token": "FuncTestToken",
  "org": "FuncTestOrg",
  "bucket": "FuncTestBucket",
  "retentionPeriodHrs": 0,
  "retentionPeriodSeconds": 0
}'
~~~

### 4. build docker image and run

~~~docker
1. docker build -t functest:golang .
2. docker run -it --name functest --network=sample-functest_functest functest:golang ./buildtest.sh env.prod
~~~

#### Optional Go test commands

~~~bash
go test -v ./functest # run tests under functest package
go test -v -run TestDemo # run test specific to "TestDemo"
go clean -testcache # clean test cache
~~~

### 5. check data in InfluxDB2

- InfluxDB2: check test data be generated in bucket = FuncTestBucket & measurement = FuncTest

### 6. config Grafana + InfluxDB2 (Optinal)

- Grafana: <http://grafana:3000> admin / admin -> skip change password
- config -> data source -> add data source
- select InfluxDB as template
- select Flux as language
- close basic auth
- Add URL: <http://influxdb2:8086> -> Org = FuncTestOrg -> Token = "FuncTestToken" -> Bucket = FuncTestBucket -> save & test
- create new dashboard and select InfluxDB data source
