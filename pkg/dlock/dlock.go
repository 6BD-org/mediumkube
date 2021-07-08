package dlock

import (
	"context"
	"encoding/json"
	"mediumkube/pkg/common"
	"mediumkube/pkg/etcd"
	"mediumkube/pkg/models"
	"time"

	"github.com/google/uuid"
	"go.etcd.io/etcd/client/v2"
	"k8s.io/klog/v2"
)

type lockManager interface {
	DoWithLock(lockType string, timeout int64, f func(), fallback func())
}

type etcdLockManager struct {
	config *common.OverallConfig
}

func _lockKey(config *common.OverallConfig, lockType string) string {
	return config.Overlay.DLockEtcdPrefix + "/" + lockType
}

// DoWithLock
// ETCD based distributed lock
// 1. Try to get lock.
// 2. If key not exists, put myLock
// 	  Otherwise check TTL of existing lock. If it's timeout, then put myLock
// 3. Reget Lock and check UUID. If same, do job, otherwize, fallback
func (m *etcdLockManager) DoWithLock(lockType string, timeout int64, f func(), fallback func()) {
	var err error = nil

ERR:
	if err != nil {
		klog.Warning("Error acquiring dlock", err)
		fallback()
		return
	}

	kpi := client.NewKeysAPI(etcd.NewClientOrDie())
	resp, err := kpi.Get(context.TODO(), _lockKey(m.config, lockType), nil)
	myLock := models.Lock{
		UUID:    uuid.New().String(),
		TIMEOUT: time.Now().UnixNano() + timeout,
	}
	lock := models.Lock{}
	payload, err := json.Marshal(&myLock)
	if err != nil {
		goto ERR
	}

	if err != nil {
		// Create if not exist
		_, err = kpi.Set(context.TODO(), _lockKey(m.config, lockType), string(payload), nil)
		if err != nil {
			goto ERR
		}
	} else {

		err = json.Unmarshal([]byte(resp.Node.Value), &lock)
		if err != nil {
			goto ERR
		}
		if lock.TIMEOUT < time.Now().UnixNano() {
			// If timeout, try to overwrite it
			_, err = kpi.Set(context.TODO(), _lockKey(m.config, lockType), string(payload), nil)
			if err != nil {
				goto ERR
			}
		}
	}

	resp, err = kpi.Get(context.TODO(), _lockKey(m.config, lockType), nil)
	if err != nil {
		goto ERR
	}
	err = json.Unmarshal([]byte(resp.Node.Value), &lock)
	if err != nil {
		goto ERR
	}

	if lock.UUID == myLock.UUID {
		f()
		myLock.TIMEOUT = 0
		payload, err = json.Marshal(&myLock)
		if err != nil {
			goto ERR
		}
		resp, err = kpi.Set(context.TODO(), _lockKey(m.config, lockType), string(payload), nil)
		if err != nil {
			goto ERR
		}

	} else {
		klog.Warning("Failed to acquire lock, falling back")
		fallback()
	}
}

func NewEtcdLockManager(config *common.OverallConfig) *etcdLockManager {
	return &etcdLockManager{
		config: config,
	}
}

func init() {
	var _ lockManager = (*etcdLockManager)(nil)
}
