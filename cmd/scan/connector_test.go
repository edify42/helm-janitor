package scan

import (
	"reflect"
	"testing"

	janitorconfig "github.com/edify42/helm-janitor/internal/config"
	client "github.com/edify42/helm-janitor/internal/eks"
	internalhelm "github.com/edify42/helm-janitor/internal/helm"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
)

func TestScanClient_Getreleases(t *testing.T) {
	type fields struct {
		Selector          string
		Namespace         string
		AllNamespaces     bool
		IncludeNamespaces string
		ExcludeNamespaces string
		Env               janitorconfig.EnvConfig
	}
	type args struct {
		c    client.EKSCluster
		a    *action.Configuration
		list internalhelm.HelmList
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*release.Release
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &ScanClient{
				Selector:          tt.fields.Selector,
				Namespace:         tt.fields.Namespace,
				AllNamespaces:     tt.fields.AllNamespaces,
				IncludeNamespaces: tt.fields.IncludeNamespaces,
				ExcludeNamespaces: tt.fields.ExcludeNamespaces,
				Env:               tt.fields.Env,
			}
			if got := sc.Getreleases(tt.args.c, tt.args.a, tt.args.list); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ScanClient.Getreleases() = %v, want %v", got, tt.want)
			}
		})
	}
}
