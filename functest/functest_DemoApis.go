package functest

import (
	"testing"

	resty "github.com/go-resty/resty/v2"
	"github.com/influxdata/influxdb-client-go/v2/domain"
)

func defaultRequest() *resty.Request {
	client := resty.New()
	client.SetDebug(true)
	req := client.R()
	req.EnableTrace().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Token "+authToken)

	return req
}

func postCreateTask(t *testing.T, taskReq domain.TaskCreateRequest) (taskRsp domain.Task, rawRsp *resty.Response) {
	req := defaultRequest()
	rawRsp, err := req.
		SetBody(taskReq).
		SetResult(&taskRsp).
		Post(serviceURL + "/api/v2/tasks")

	printWrite(t, rawRsp, err)
	return
}

func getRetrieveTask(t *testing.T, taskId string) (taskRsp domain.Task, rawRsp *resty.Response) {
	req := defaultRequest()
	rawRsp, err := req.
		SetResult(&taskRsp).
		Get(serviceURL + "/api/v2/tasks/" + taskId)

	printWrite(t, rawRsp, err)
	return
}

func patchUpdateTask(t *testing.T, taskId string, taskUpdateReq domain.TaskUpdateRequest) (taskRsp domain.Task, rawRsp *resty.Response) {
	req := defaultRequest()
	rawRsp, err := req.
		SetBody(taskUpdateReq).
		SetResult(&taskRsp).
		Patch(serviceURL + "/api/v2/tasks/" + taskId)

	printWrite(t, rawRsp, err)
	return
}

func delDeleteTask(t *testing.T, taskId string) *resty.Response {
	req := defaultRequest()
	rawRsp, err := req.
		Delete(serviceURL + "/api/v2/tasks/" + taskId)

	printWrite(t, rawRsp, err)
	return rawRsp
}
