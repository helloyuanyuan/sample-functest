package functest

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	resty "github.com/go-resty/resty/v2"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const (
	PROFILE_YAML = "profile.yaml"
)

var config *viper.Viper
var influxdb2Config *viper.Viper

func init() {
	cw, err := os.Getwd()
	if err != nil {
		panic(err)
	}

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

	if !found {
		panic("no configuration file found")
	}

	viper := LoadConfig(configDir, PROFILE_YAML)
	influxdb2Config = viper
	env := viper.GetString("env")
	if len(strings.TrimSpace(env)) == 0 {
		log.Panic("No environment specified")
	}

	config = LoadConfig(configDir, env+".yaml")
}

func LoadConfig(path, file string) *viper.Viper {
	viper := viper.New()
	// viper.AddConfigPath(path)
	viper.SetConfigFile(filepath.Join(path, file))
	//viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("can not load file %s", filepath.Join(path, file))
	}

	return viper
}

func isConfigDir(target string) (string, bool) {
	configDir := filepath.Join(target, "config")
	fi, statErr := os.Stat(filepath.Join(configDir, PROFILE_YAML))
	if statErr != nil {
		return "", false
	}

	return configDir, fi.Mode().IsRegular()
}

func httpGET(t *testing.T) (rawRsp *resty.Response) {
	client := resty.New()
	client.SetDebug(true)
	rawRsp, err := client.R().
		EnableTrace().
		Get(config.GetString("service.url") + "/get")

	printInfo(rawRsp, err)

	influxdb2Write(t, rawRsp)
	assert.Nil(t, err)
	return
}

func httpPOST(t *testing.T) (rawRsp *resty.Response) {
	client := resty.New()
	client.SetDebug(true)
	rawRsp, err := client.R().
		EnableTrace().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"username":"user@test.com", "password":"12345"}`).
		// SetResult(&AuthSuccess{}).
		Post(config.GetString("service.url") + "/public/users/login")

	printInfo(rawRsp, err)

	influxdb2Write(t, rawRsp)
	assert.Nil(t, err)
	return
}

func printInfo(rawRsp *resty.Response, err error) {
	// Explore response object
	fmt.Println("Response Info:")
	fmt.Println("  Error      :", err)
	fmt.Println("  Status Code:", rawRsp.StatusCode())
	fmt.Println("  Status     :", rawRsp.Status())
	fmt.Println("  Proto      :", rawRsp.Proto())
	fmt.Println("  Time       :", rawRsp.Time())
	fmt.Println("  Received At:", rawRsp.ReceivedAt())
	fmt.Println("  Body       :\n", rawRsp)
	fmt.Println()

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

func influxdb2Write(t *testing.T, rawRsp *resty.Response) {
	serverURL := influxdb2Config.GetString("influxdb2.url")
	authToken := influxdb2Config.GetString("influxdb2.token")
	org := influxdb2Config.GetString("influxdb2.org")
	bucket := influxdb2Config.GetString("influxdb2.bucket")

	client := influxdb2.NewClient(serverURL, authToken)
	writeAPI := client.WriteAPI(org, bucket)

	point := influxdb2.NewPointWithMeasurement(influxdb2Config.GetString("influxdb2.measurement")).
		AddTag("Env", influxdb2Config.GetString("env")).
		AddTag("TestCase", t.Name()).
		AddTag("ApiCall", rawRsp.Request.URL).
		AddTag("StatusCode", strconv.Itoa(rawRsp.StatusCode())).
		AddField("TimeDuration", rawRsp.Time()).
		SetTime(time.Now())

	writeAPI.WritePoint(point)
	writeAPI.Flush()

	defer client.Close()
}
