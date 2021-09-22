package eks

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	janitorconfig "github.com/lendi-au/helm-janitor/internal/config"
	"sigs.k8s.io/aws-iam-authenticator/pkg/token"
)

type mockGeneratorType struct{}

func (m *mockGeneratorType) DescribeCluster(*eks.Client, string) (*eks.DescribeClusterOutput, error) {
	name := "yeah"
	cert := "some-fake-certificate"
	return &eks.DescribeClusterOutput{
		Cluster: &types.Cluster{
			Name:                 &name,
			Endpoint:             &name,
			CertificateAuthority: &types.Certificate{Data: &cert},
		},
	}, nil
}

func (m *mockGeneratorType) GenToken(*string) (token.Token, error) {
	return token.Token{
		Token: "test",
	}, nil
}

func (m *mockGeneratorType) TestCluster(*types.Cluster, token.Token) error {
	return nil
}

func (m *mockGeneratorType) WriteCA(a AwsConfig, e *eks.DescribeClusterOutput) string {
	return "test"
}

func TestAwsConfig_Init(t *testing.T) {
	type fields struct {
		J     janitorconfig.EnvConfig
		Token string
	}
	type args struct {
		cfg aws.Config
		g   Generator
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   EKSCluster
	}{
		{
			name: "first valid test",
			fields: fields{
				J: janitorconfig.EnvConfig{
					Cluster:         "unit-test",
					TmpFileLocation: "/tmp",
					TmpFilePrefix:   "none",
				},
				Token: "a token",
			},
			args: args{
				cfg: aws.Config{
					Region: "ap-southeast-2",
				},
				g: &mockGeneratorType{},
			},
			want: EKSCluster{
				Name:     "yeah",
				Endpoint: "yeah",
				Token:    "test",
				CAFile:   "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AwsConfig{
				J:     tt.fields.J,
				Token: tt.fields.Token,
			}
			if got := a.Init(tt.args.cfg, tt.args.g); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AwsConfig.Init() = %v, want %v", got, tt.want)
			}
		})
	}
}
