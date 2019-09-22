package cmd

import
	"github.com/jenkins-x/jx/pkg/gits"
	"github.com/spf13/cobra"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

const (
	rbacCheckImage         = "nvisium/jxtools-rbac-check"
	rbacCheckContainerName = "jx-tools-rbac-check"
	rbacCheckNamespace     = "jx-tools-rbac-check"
	rbacCheckJobName       = "jx-tools-rbac-check-job"
	outputFormatYAML = "yaml"
)

type ScanRbacOptions struct {
	ScanOptions

	Output string
}

func (o *ScanRbacOptions) Run() error {

}

func createRbacScanJob(name string, namespace string, container *v1.Container) *batchv1.Job {

}

// NewCmdScanRbac creates a command object for "scan rbac" command
func NewCmdScanRbac(commonOpts *opts.CommonOptions) *cobra.Command {
	options := &ScanClusterOptions{
		ScanOptions: ScanOptions{
			CommonOptions: commonOpts,
		},
	}

	cmd := &cobra.Command{
		Use:   "rbac",
		Short: "Performs rbac scans ",
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := options.Run()
			helper.CheckErr(err)
		},
	}

	cmd.Flags().StringVarP(&options.Output, "output", "o", "plain", "output format is one of: yaml|plain")

	return cmd
}