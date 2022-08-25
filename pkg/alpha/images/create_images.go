/*
 Copyright 2022 The KubeSphere Authors.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package images

import (
	"errors"
	"fmt"

	"github.com/kubesphere/kubekey/pkg/bootstrap/precheck"
	"github.com/kubesphere/kubekey/pkg/common"
	"github.com/kubesphere/kubekey/pkg/container"
	"github.com/kubesphere/kubekey/pkg/core/module"
	"github.com/kubesphere/kubekey/pkg/core/pipeline"
	"github.com/kubesphere/kubekey/pkg/images"
	"github.com/kubesphere/kubekey/pkg/kubernetes"
)

func NewCreateImagesPipeline(runtime *common.KubeRuntime) error {
	skipPushImages := runtime.Arg.Artifact == "" || runtime.Cluster.Registry.PrivateRegistry == ""
	m := []module.Module{
		&precheck.NodePreCheckModule{},
		&kubernetes.StatusModule{},
		&container.InstallContainerModule{},
		&images.CopyImagesToRegistryModule{Skip: skipPushImages},
		&images.PullModule{},
	}

	p := pipeline.Pipeline{
		Name:    "CreateImagesPipeline",
		Modules: m,
		Runtime: runtime,
	}
	if err := p.Start(); err != nil {
		return err
	}
	return nil
}

func CreateImages(args common.Argument, downloadCmd string) error {
	args.DownloadCommand = func(path, url string) string {
		// this is an extension point for downloading tools, for example users can set the timeout, proxy or retry under
		// some poor network environment. Or users even can choose another cli, it might be wget.
		// perhaps we should have a build-in download function instead of totally rely on the external one
		return fmt.Sprintf(downloadCmd, path, url)
	}
	var loaderType string

	if args.FilePath != "" {
		loaderType = common.File
	} else {
		loaderType = common.AllInOne
	}

	runtime, err := common.NewKubeRuntime(loaderType, args)
	if err != nil {
		return err
	}
	switch runtime.Cluster.Kubernetes.Type {
	case common.Kubernetes:
		if err := NewCreateImagesPipeline(runtime); err != nil {
			return err
		}
	default:
		return errors.New("unsupported cluster kubernetes type")
	}

	return nil
}
