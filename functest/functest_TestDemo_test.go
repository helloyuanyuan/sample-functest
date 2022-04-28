package functest

import (
	"testing"
	"time"

	"github.com/influxdata/influxdb-client-go/v2/domain"
	"github.com/stretchr/testify/assert"
)

const (
	INVALID_TASK_ID     = "InvalidTaskId123"
	INVALID_TASK_ID_ERR = "failed to decode request: invalid ID"
)

func TestPost(t *testing.T) {
	name := "TestTask" + time.Now().String()
	flux := "option task = { \n  name: \"" + name + "\",\n  every: 10h,\n  offset: 10m\n}"
	taskReq := domain.TaskCreateRequest{Org: &org, Flux: flux}

	taskRsp, rawRsp := postCreateTask(t, taskReq)
	assert.Equal(t, 201, rawRsp.StatusCode())
	assert.NotNil(t, taskRsp.Id)
}

func TestGet(t *testing.T) {
	taskRsp, rawRsp := getRetrieveTask(t, INVALID_TASK_ID)
	assert.Equal(t, 400, rawRsp.StatusCode())
	assert.Contains(t, string(rawRsp.Body()), INVALID_TASK_ID_ERR)
	assert.Empty(t, taskRsp)
}

func TestUpdate(t *testing.T) {
	nameUpdated := "Updated Name"
	descriptionNew := "Test Description New Added"
	taskUpdateReq := domain.TaskUpdateRequest{Name: &nameUpdated, Description: &descriptionNew}

	taskRsp, rawRsp := patchUpdateTask(t, INVALID_TASK_ID, taskUpdateReq)
	assert.Equal(t, 400, rawRsp.StatusCode())
	assert.Contains(t, string(rawRsp.Body()), INVALID_TASK_ID_ERR)
	assert.Empty(t, taskRsp)
}

func TestDelete(t *testing.T) {
	rawRsp := delDeleteTask(t, INVALID_TASK_ID)
	assert.Equal(t, 400, rawRsp.StatusCode())
	assert.Contains(t, string(rawRsp.Body()), INVALID_TASK_ID_ERR)
}

func TestPostGet(t *testing.T) {
	name := "TestTask" + time.Now().String()
	flux := "option task = { \n  name: \"" + name + "\",\n  every: 10h,\n  offset: 10m\n}"
	taskReq := domain.TaskCreateRequest{Org: &org, Flux: flux}

	taskRsp, _ := postCreateTask(t, taskReq)
	taskId := taskRsp.Id

	taskRsp, rawRsp := getRetrieveTask(t, taskId)
	assert.Equal(t, 200, rawRsp.StatusCode())
	assert.NotNil(t, taskRsp.Name)
}

func TestPostUpdate(t *testing.T) {
	name := "TestTask" + time.Now().String()
	flux := "option task = { \n  name: \"" + name + "\",\n  every: 10h,\n  offset: 10m\n}"
	taskReq := domain.TaskCreateRequest{Org: &org, Flux: flux}

	taskRsp, _ := postCreateTask(t, taskReq)
	taskId := taskRsp.Id

	nameUpdated := "Updated" + name
	descriptionNew := "Test Description New Added"
	taskUpdateReq := domain.TaskUpdateRequest{Name: &nameUpdated, Description: &descriptionNew}

	taskRsp, rawRsp := patchUpdateTask(t, taskId, taskUpdateReq)
	assert.Equal(t, 200, rawRsp.StatusCode())
	assert.Equal(t, &nameUpdated, &taskRsp.Name)
	assert.Equal(t, &descriptionNew, taskRsp.Description)
}

func TestPostDelete(t *testing.T) {
	name := "TestTask" + time.Now().String()
	flux := "option task = { \n  name: \"" + name + "\",\n  every: 10h,\n  offset: 10m\n}"
	taskReq := domain.TaskCreateRequest{Org: &org, Flux: flux}

	taskRsp, _ := postCreateTask(t, taskReq)
	taskId := taskRsp.Id

	rawRsp := delDeleteTask(t, taskId)
	assert.Equal(t, 204, rawRsp.StatusCode())
}
