package eks

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"sigs.k8s.io/aws-iam-authenticator/pkg/token"
)

func TestNew(t *testing.T) {
	yep := `LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUV1VENDQTZHZ0F3SUJBZ0lRUUJyRVpDR3pF
eUVERHJ2a0VockZIVEFOQmdrcWhraUc5dzBCQVFzRkFEQ0IKdlRFTE1Ba0dBMVVFQmhNQ1ZWTXhG
ekFWQmdOVkJBb1REbFpsY21sVGFXZHVMQ0JKYm1NdU1SOHdIUVlEVlFRTApFeFpXWlhKcFUybG5i
aUJVY25WemRDQk9aWFIzYjNKck1Ub3dPQVlEVlFRTEV6RW9ZeWtnTWpBd09DQldaWEpwClUybG5i
aXdnU1c1akxpQXRJRVp2Y2lCaGRYUm9iM0pwZW1Wa0lIVnpaU0J2Ym14NU1UZ3dOZ1lEVlFRREV5
OVcKWlhKcFUybG5iaUJWYm1sMlpYSnpZV3dnVW05dmRDQkRaWEowYVdacFkyRjBhVzl1SUVGMWRH
aHZjbWwwZVRBZQpGdzB3T0RBME1ESXdNREF3TURCYUZ3MHpOekV5TURFeU16VTVOVGxhTUlHOU1R
c3dDUVlEVlFRR0V3SlZVekVYCk1CVUdBMVVFQ2hNT1ZtVnlhVk5wWjI0c0lFbHVZeTR4SHpBZEJn
TlZCQXNURmxabGNtbFRhV2R1SUZSeWRYTjAKSUU1bGRIZHZjbXN4T2pBNEJnTlZCQXNUTVNoaktT
QXlNREE0SUZabGNtbFRhV2R1TENCSmJtTXVJQzBnUm05eQpJR0YxZEdodmNtbDZaV1FnZFhObElH
OXViSGt4T0RBMkJnTlZCQU1UTDFabGNtbFRhV2R1SUZWdWFYWmxjbk5oCmJDQlNiMjkwSUVObGNu
UnBabWxqWVhScGIyNGdRWFYwYUc5eWFYUjVNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUYKQUFPQ0FR
OEFNSUlCQ2dLQ0FRRUF4MkUzWHJFQk5OdGkxeFdiLzFoYWpDTWoxbUNPa2RlUW1JTjY1bGdaT0l6
Rgo5dVZraGJTaWNmdnR2Ym5helUwQXRNZ3RjNlhIYVhHVkh6azhza1FIbk9nTytrMUt4Q0hmS1dH
UE1pSmhnc1dICkgyNk1mRjhXSUZGRTBYQlBWK3JqSE9QTWVlNVkyQTdDczBXVHdDem5taGNyZXdB
M2VrRXplT0V6NHZNUUduK0gKTEw3MjlmZEM0dVcvaDJLSlh3QkwzOFhkNUhWRU1rRTZIbkZ1YWNz
TGRVWUkwY3JTSzVYUXovdTVRR3RrakZkTgovQk1SZVlUdFhsVDJOSjhJQWZNUUpRWVhTdHJ4SFhw
bWE1aGdacVRaNzlJdWd2SHc3d25xUk1rVmF1SURialBUCnJKOVZBTWYyQ0dxVXVWL2M0RFB4aEdE
NVd5Y1J0UHdXOHJ0V2FvQWxqUUlEQVFBQm80R3lNSUd2TUE4R0ExVWQKRXdFQi93UUZNQU1CQWY4
d0RnWURWUjBQQVFIL0JBUURBZ0VHTUcwR0NDc0dBUVVGQndFTUJHRXdYNkZkb0ZzdwpXVEJYTUZV
V0NXbHRZV2RsTDJkcFpqQWhNQjh3QndZRkt3NERBaG9FRkkvbDB4cUdySTJPYThQUGdHclVTQmdz
CmV4a3VNQ1VXSTJoMGRIQTZMeTlzYjJkdkxuWmxjbWx6YVdkdUxtTnZiUzkyYzJ4dloyOHVaMmxt
TUIwR0ExVWQKRGdRV0JCUzJkL3BwU0VlZlV4TFZ3dW9ITW5ZSDBaY0hHVEFOQmdrcWhraUc5dzBC
QVFzRkFBT0NBUUVBU3ZqNApzQVBtTEdkNzVKUjNZOHh1VFBsOURnM2N5TGsxdVhCUFkvb2srbXlE
akVlZE8yUHptdmwyTXBXUnNYZThySnErCnNlUXhJY2FCbFZaYURySEMxTEdtV2F6eFk4dTRUQjFa
a0VydmtCWW9IMXF1RVB1QlVEZ01iTXp4UGNQMVkrT3oKNHlISkpEbnAvUlZtUnZRYkVkQk5jNk45
UnZrOTdhaGZZdFR4UC9qZ2RGY3JHSjJCdE1RbzJwU1hwWERyckIyKwpCeEh3MWR2ZDVZencxVEt3
ZytaWDRvKy92cUdxdnowZHRkUTQ2dGV3WERwUGFqK1B3R1pzWTZycDJhUVc5SUhSCmxSUU9mYzJW
Tk5uU2ozQnpnWHVjZnIyWVlkaEZoNWlReGV1R01NWTF2L0QvdzFXSWcwdnZCWklHY2ZLNG1KTzMK
N00yQ1lmRTQ1aytYbUNwYWpRPT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=` // cat /etc/ssl/certs/VeriSign_Universal_Root_Certification_Authority.pem | base64
	endpoint := "http://localhost"

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
					Endpoint:             &endpoint,
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
