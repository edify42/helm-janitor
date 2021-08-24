package eks

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"sigs.k8s.io/aws-iam-authenticator/pkg/token"
)

func TestNew(t *testing.T) {
	yep := "yeps"
	// ca, _ := base64.StdEncoding.DecodeString(yep)
	// cluster, _ := kubernetes.NewForConfig(
	// 	&rest.Config{
	// 		Host:        yep,
	// 		BearerToken: "this",
	// 		TLSClientConfig: rest.TLSClientConfig{
	// 			CAData: ca,
	// 		},
	// 	},
	// )
	type args struct {
		cluster *types.Cluster
		tok     token.Token
	}
	tests := []struct {
		name string
		args args
		// want    *kubernetes.Clientset
		wantErr bool
	}{
		{
			name: "first",
			args: args{
				cluster: &types.Cluster{
					Endpoint:             &yep,
					CertificateAuthority: &types.Certificate{Data: &yep},
				},
				tok: token.Token{
					Token: "this",
				},
			},
			// want:    cluster,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.args.cluster, tt.args.tok)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("New() = %v, want %v", got, tt.want)
			// }
		})
	}
}
