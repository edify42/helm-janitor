package delete

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/lendi-au/helm-janitor/internal/eks"
	client "github.com/lendi-au/helm-janitor/internal/eks"
	internalhelm "github.com/lendi-au/helm-janitor/internal/helm"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
)

type mock struct{}

func (m *mock) Config() *Client {
	return &Client{}
}

func (m *mock) Init() {}

func (m *mock) Makeawscfg() aws.Config {
	return aws.Config{}
}

func (m *mock) Getekscluster(aws.Config, client.Generator) client.EKSCluster {
	return client.EKSCluster{}
}

func (m *mock) Deleterelease(client.EKSCluster, *action.Configuration, *release.Release, internalhelm.HelmDelete) error {
	return nil
}

func (m *mock) Makeekscfg() client.Generator {
	return &eks.GeneratorType{}
}

func TestRunV2(t *testing.T) {
	type args struct {
		sr InputRun
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "first test",
			args: args{
				sr: &mock{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RunV2(tt.args.sr)
		})
	}
}
