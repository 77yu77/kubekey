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

package phase

import (
	"github.com/kubesphere/kubekey/cmd/ctl/options"
	"github.com/kubesphere/kubekey/cmd/ctl/util"
	"github.com/kubesphere/kubekey/pkg/alpha/etcd"
	"github.com/kubesphere/kubekey/pkg/common"
	"github.com/spf13/cobra"
)

type CreateEtcdOptions struct {
	CommonOptions  *options.CommonOptions
	ClusterCfgFile string
	DownloadCmd    string
}

func NewCreateEtcdOptions() *CreateEtcdOptions {
	return &CreateEtcdOptions{
		CommonOptions: options.NewCommonOptions(),
	}
}

// NewCmdCreateEtcd creates a new install etcd command
func NewCmdCreateEtcd() *cobra.Command {
	o := NewCreateEtcdOptions()
	cmd := &cobra.Command{
		Use:   "etcd",
		Short: "install the etcd on the master",
		Run: func(cmd *cobra.Command, args []string) {
			util.CheckErr(o.Run())
		},
	}

	o.CommonOptions.AddCommonFlag(cmd)
	o.AddFlags(cmd)
	return cmd
}

func (o *CreateEtcdOptions) Run() error {
	arg := common.Argument{
		FilePath: o.ClusterCfgFile,
		Debug:    o.CommonOptions.Verbose,
	}
	return etcd.CreateEtcd(arg, o.DownloadCmd)
}

func (o *CreateEtcdOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&o.ClusterCfgFile, "filename", "f", "", "Path to a configuration file")
	cmd.Flags().StringVarP(&o.DownloadCmd, "download-cmd", "", "curl -L -o %s %s",
		`The user defined command to download the necessary binary files. The first param '%s' is output path, the second param '%s', is the URL`)

}
