package functest

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	resty "github.com/go-resty/resty/v2"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const (
	MAIN_YAML = "main.yaml"
)

var mainConfig, envConfig *viper.Viper
var influxdb2URL, authToken, org, bucket, measurement string
var serviceURL string

func init() {
	cw, _ := os.Getwd()

	depth := 0
	found := false
	configDir := cw

	for {
		configDir, found = isConfigDir(cw)

		if !found {
			cw = filepath.Join(cw, "..")
		}

		if depth > 10 || found {
			break
		}
	}

	mainConfig = LoadConfig(configDir, MAIN_YAML)
	envConfig = LoadConfig(configDir, mainConfig.GetString("env")+".yaml")

	influxdb2URL = mainConfig.GetString("influxdb2.url")
	authToken = mainConfig.GetString("influxdb2.token")
	org = mainConfig.GetString("influxdb2.org")
	bucket = mainConfig.GetString("influxdb2.bucket")
	measurement = mainConfig.GetString("influxdb2.measurement")

	serviceURL = envConfig.GetString("service.url")

}

func LoadConfig(path, file string) *viper.Viper {
	viper := viper.New()
	viper.SetConfigFile(filepath.Join(path, file))
	viper.ReadInConfig()

	return viper
}

func isConfigDir(target string) (string, bool) {
	configDir := filepath.Join(target, "config")
	fi, statErr := os.Stat(filepath.Join(configDir, MAIN_YAML))
	if statErr != nil {
		return "", false
	}

	return configDir, fi.Mode().IsRegular()
}

func printWrite(t *testing.T, rawRsp *resty.Response, err error) {
	assert.Nil(t, err)
	printResult(rawRsp)
	writeResult(t, rawRsp)
}

func printResult(rawRsp *resty.Response) {
	// Explore trace info
	fmt.Println("Request Trace Info:")
	ti := rawRsp.Request.TraceInfo()
	fmt.Println("  DNSLookup     :", ti.DNSLookup)
	fmt.Println("  ConnTime      :", ti.ConnTime)
	fmt.Println("  TCPConnTime   :", ti.TCPConnTime)
	fmt.Println("  TLSHandshake  :", ti.TLSHandshake)
	fmt.Println("  ServerTime    :", ti.ServerTime)
	fmt.Println("  ResponseTime  :", ti.ResponseTime)
	fmt.Println("  TotalTime     :", ti.TotalTime)
	fmt.Println("  IsConnReused  :", ti.IsConnReused)
	fmt.Println("  IsConnWasIdle :", ti.IsConnWasIdle)
	fmt.Println("  ConnIdleTime  :", ti.ConnIdleTime)
	fmt.Println("  RequestAttempt:", ti.RequestAttempt)
	fmt.Println("  RemoteAddr    :", ti.RemoteAddr.String())
}

func writeResult(t *testing.T, rawRsp *resty.Response) {
	client := influxdb2.NewClient(influxdb2URL, authToken)
	writeAPI := client.WriteAPI(org, bucket)

	point := influxdb2.NewPointWithMeasurement(measurement).
		SetTime(time.Now()).
		AddTag("Env", mainConfig.GetString("env")).
		AddField("Api", rawRsp.Request.URL).
		AddTag("Duration", rawRsp.Time().String()).
		AddTag("TestCase", t.Name()).
		AddTag("Status", rawRsp.Status())

	writeAPI.WritePoint(point)
	writeAPI.Flush()

	defer client.Close()
	fmt.Println("==============================================================================")
	fmt.Println("Result Wrote To InfluxDB.")
	fmt.Println("==============================================================================")
}
