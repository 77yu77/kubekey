package binary

import (
	"fmt"
	"os/exec"

	kubekeyapiv1alpha2 "github.com/kubesphere/kubekey/apis/kubekey/v1alpha2"
	"github.com/kubesphere/kubekey/pkg/common"
	"github.com/kubesphere/kubekey/pkg/core/cache"
	"github.com/kubesphere/kubekey/pkg/core/connector"
	"github.com/kubesphere/kubekey/pkg/core/util"
	"github.com/kubesphere/kubekey/pkg/files"
	"github.com/pkg/errors"
)

type GetBinaryPath struct {
	common.KubeAction
	BinaryNames []string
}

func (g *GetBinaryPath) Execute(runtime connector.Runtime) error {
	cfg := g.KubeConf.Cluster

	var kubeVersion string
	if cfg.Kubernetes.Version == "" {
		kubeVersion = kubekeyapiv1alpha2.DefaultKubeVersion
	} else {
		kubeVersion = cfg.Kubernetes.Version
	}

	archMap := make(map[string]bool)
	for _, host := range cfg.Hosts {
		switch host.Arch {
		case "amd64":
			archMap["amd64"] = true
		case "arm64":
			archMap["arm64"] = true
		default:
			return errors.New(fmt.Sprintf("Unsupported architecture: %s", host.Arch))
		}
	}

	for arch := range archMap {
		if err := setK8sBinaryPath(g.KubeConf, runtime.GetWorkDir(), kubeVersion, arch, g.PipelineCache, g.BinaryNames); err != nil {
			return err
		}
	}
	return nil
}

func setK8sBinaryPath(kubeConf *common.KubeConf, path, version, arch string, pipelineCache *cache.Cache, binaryNames []string) error {
	allBinaryNames := map[string]bool{
		"etcd":    true,
		"kubeadm": true,
		"kubelet": true,
		"kubectl": true,
		"kubecni": true,
		"helm":    true,
		"crictl":  true,
	}

	binariesMap := make(map[string]*files.KubeBinary)
	for _, binary := range binaryNames {
		if _, ok := allBinaryNames[binary]; !ok {
			return errors.New(fmt.Sprintf("Unsupported binary name to get path: %s", binary))
		}
		kubeBinary := files.NewKubeBinary(binary, arch, kubekeyapiv1alpha2.DefaultEtcdVersion, path, kubeConf.Arg.DownloadCommand)
		binariesMap[kubeBinary.ID] = kubeBinary

		if !util.IsExist(kubeBinary.BaseDir) {
			return errors.New("BaseDir of download binary is not exist")
		}
		if util.IsExist(kubeBinary.Path()) {
			if err := kubeBinary.SHA256Check(); err != nil {
				p := kubeBinary.Path()
				_ = exec.Command("/bin/sh", "-c", fmt.Sprintf("rm -f %s", p)).Run()
				if err := kubeBinary.Download(); err != nil {
					return fmt.Errorf("Failed to download %s binary: %s error: %w ", kubeBinary.ID, kubeBinary.GetCmd(), err)
				}
			}
		}

	}

	pipelineCache.Set(common.KubeBinaries+"-"+arch, binariesMap)
	return nil
}
