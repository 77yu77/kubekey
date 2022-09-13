package etcd

import (
	"github.com/kubesphere/kubekey/pkg/alpha/binary"
	"github.com/kubesphere/kubekey/pkg/common"
	"github.com/kubesphere/kubekey/pkg/core/task"
	"github.com/kubesphere/kubekey/pkg/etcd"
)

type PreCheckModule struct {
	common.KubeModule
	Skip bool
}

func (p *PreCheckModule) IsSkip() bool {
	return p.Skip
}

func (p *PreCheckModule) Init() {
	p.Name = "ETCDPreCheckModule"
	p.Desc = "Get ETCD cluster status"

	getStatus := &task.RemoteTask{
		Name:     "GetETCDStatus",
		Desc:     "Get etcd status",
		Hosts:    p.Runtime.GetHostsByRole(common.ETCD),
		Action:   new(etcd.GetStatus),
		Parallel: false,
		Retry:    0,
	}

	setBinaryCache := &task.LocalTask{
		Name:   "SetEtcdBinaryCache",
		Desc:   "Set Etcd Binary Path in PipelineCache",
		Action: new(binary.GetEtcdBinaryPath),
	}

	p.Tasks = []task.Interface{
		getStatus,
		setBinaryCache,
	}
}
