package cmd

import(
	"github.com/jenkins-x/jx/pkg/gits"
	"github.com/spf13/cobra"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/google/uuid"
	"github.com/jenkins-x/jx/pkg/cmd/opts"
)

const (
	dependencyCheckImage         = "owasp/dependency-check"
	dependencyCheckContainerName = "jx-tools-dependency-check"
	dependencyCheckNamespace     = "jx-tools-dependency-check"
	dependencyCheckJobName       = "jx-tools-dependency-check-job"
	outputFormatYAML = "yaml"
)

type scanDependencyResult struct {
	Dependencies		[]string
	Analyzers			[]string
	Vulnerabilities		[]string
}

// NewCmdScanDependencyCheck creates a command object for "scan dependency-check" command
func NewCmdScanDependencyCheck(commonOpts *opts.CommonOptions) *cobra.Command {
	options := &ScanClusterOptions{
		ScanOptions: ScanOptions{
			CommonOptions: commonOpts,
		},
	}

	cmd := &cobra.Command{
		Use:   "dependency-check",
		Short: "Performs dependency check scans ",
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

func (o *ScanDependencyOptions) Run() error {
	cloneurl := o.Git().CloneURL
	// First, clone the repo so we can analyze it.
	err = o.Git().Clone(cloneurl, "~/")
	if err != nil {
		return err
	}
	kubeClient, err := o.KubeClient()
	if err != nil {
		return errors.Wrap(err, "creating kube client")
	}
	ns, err := createDependendencyCheckNamespace()
	if err != nil {
		return err
	}
	container := o.dependencyCheckContainer()
	job := o.createDependencyCheckScanJob(dependencyCheckJobName, ns, container)
	job, err = kubeClient.BatchV1().Jobs(ns).Create(job)
	if err != nil {
		return err
	}
}
// Create a dedicated namespace for DependencyCheck scan
func createDepedendencyCheckNamespace()(string, err) {
	ns := fmt.Sprintf("%s-%s", dependencyCheckNamespace, uuid.New())
	namespace := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: ns,
		},
	}
	_, err = kubeClient.CoreV1().Namespaces().Create(namespace)
	if err != nil {
		return "", errors.Wrapf(err, "creating namespace '%s'", ns)
	}
	return ns, nil
}

func createDependencyCheckScanJob(name string, namespace string, container *v1.Container) *batchv1.Job {
	podTmpl := v1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1.PodSpec{
			Containers:    []v1.Container{*container},
			RestartPolicy: v1.RestartPolicyNever,
		},
	}

	return &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: batchv1.JobSpec{
			Template: podTmpl,
		},
	}
}

// ScanDependencyOptions the options for 'scan dependency' command
type ScanDependencyOptions struct {
	ScanOptions

	Output string
}

func (o *ScanDependencyOptions) dependencyCheckContainer() *v1.Container {
	return &v1.Container{
		Name:            DependencyCheckContainerName,
		Image:           o.dependencyCheckImage(),
		ImagePullPolicy: v1.PullAlways,
		Command:         []string{"./bin/dependency-check.sh"},
		Args:            []string{"--project depscan", "--out .",fmt.Fprintf("--scan %s", o.Git().Name)},
	}
}