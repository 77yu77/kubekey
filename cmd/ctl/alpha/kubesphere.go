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

package alpha

import (
	"fmt"
	"time"

	"github.com/kubesphere/kubekey/cmd/ctl/options"
	"github.com/kubesphere/kubekey/cmd/ctl/util"
	alpha "github.com/kubesphere/kubekey/pkg/alpha/kubesphere"
	"github.com/kubesphere/kubekey/pkg/common"
	"github.com/kubesphere/kubekey/pkg/version/kubesphere"
	"github.com/spf13/cobra"
)

type UpgradeKubeSphereOptions struct {
	CommonOptions    *options.CommonOptions
	ClusterCfgFile   string
	EnableKubeSphere bool
	KubeSphere       string
	DownloadCmd      string
}

func NewUpgradeKubeSphereOptions() *UpgradeKubeSphereOptions {
	return &UpgradeKubeSphereOptions{
		CommonOptions: options.NewCommonOptions(),
	}
}

// NewCmdUpgradeKubeSphere creates a new UpgradeKubeSphere command
func NewCmdUpgradeKubeSphere() *cobra.Command {
	o := NewUpgradeKubeSphereOptions()
	cmd := &cobra.Command{
		Use:   "kubesphere",
		Short: "upgrade your kubesphere to a newer version with this command",
		Run: func(cmd *cobra.Command, args []string) {
			util.CheckErr(o.Complete(cmd, args))
			util.CheckErr(o.Run())
		},
	}
	o.CommonOptions.AddCommonFlag(cmd)
	o.AddFlags(cmd)

	if err := ksCompletionSetting(cmd); err != nil {
		panic(fmt.Sprintf("Got error with the completion setting"))
	}
	return cmd
}

func (o *UpgradeKubeSphereOptions) Complete(cmd *cobra.Command, args []string) error {
	var ksVersion string
	if o.EnableKubeSphere && len(args) > 0 {
		ksVersion = args[0]
	} else {
		ksVersion = kubesphere.Latest().Version
	}
	o.KubeSphere = ksVersion
	return nil
}

func (o *UpgradeKubeSphereOptions) Run() error {
	arg := common.Argument{
		FilePath:         o.ClusterCfgFile,
		KsEnable:         o.EnableKubeSphere,
		KsVersion:        o.KubeSphere,
		SkipConfirmCheck: o.CommonOptions.SkipConfirmCheck,
		Debug:            o.CommonOptions.Verbose,
	}
	return alpha.UpgradeKubeSphere(arg, o.DownloadCmd)
}

func (o *UpgradeKubeSphereOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&o.ClusterCfgFile, "filename", "f", "", "Path to a configuration file")
	cmd.Flags().BoolVarP(&o.EnableKubeSphere, "with-kubesphere", "", false, fmt.Sprintf("Deploy a specific version of kubesphere (default %s)", kubesphere.Latest().Version))
	cmd.Flags().StringVarP(&o.DownloadCmd, "download-cmd", "", "curl -L -o %s %s",
		`The user defined command to download the necessary binary files. The first param '%s' is output path, the second param '%s', is the URL`)
}

func ksCompletionSetting(cmd *cobra.Command) (err error) {
	cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) (
		strings []string, directive cobra.ShellCompDirective) {
		versionArray := kubesphere.VersionsStringArr()
		versionArray = append(versionArray, time.Now().Add(-time.Hour*24).Format("nightly-20060102"))
		return versionArray, cobra.ShellCompDirectiveNoFileComp
	}

	return
}
