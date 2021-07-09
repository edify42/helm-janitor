package scan

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/time"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func TestNameList(t *testing.T) {
	type args struct {
		r []*release.Release
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "simple test",
			args: args{
				r: []*release.Release{
					{Name: "this"},
				},
			},
			want: []string{
				"this",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NameList(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NameList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func annotationsGen(timestamp string) []byte {
	return []byte(fmt.Sprintf(`{"janitorAnnotations": {
		"janitor/expires": "%s"}}`, timestamp))
}

func TestCheckReleaseExpired(t *testing.T) {
	type args struct {
		r release.Release
	}
	exp := annotationsGen("2021-07-07T10:56:09Z") // Expired (past)
	var expiredDat map[string]interface{}
	if err := json.Unmarshal(exp, &expiredDat); err != nil {
		panic(err)
	}
	layout := "2006-01-02T15:04:05Z"
	now := time.Now()                                           // should always be in the future unless the date is wrong on the system.
	pastTime, err := time.Parse(layout, "2021-06-07T10:56:09Z") // month behind of above dat
	t.Log(pastTime)
	if err != nil {
		panic(err)
	}

	futureTime := now.AddDate(0, 1, 0)
	fut := annotationsGen(futureTime.Format(layout)) // future date
	var futureDate map[string]interface{}
	if err := json.Unmarshal(fut, &futureDate); err != nil {
		panic(err)
	}

	t.Log(pastTime)
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "first test expired - success as expired",
			args: args{
				r: release.Release{
					Name: "release",
					Info: &release.Info{
						LastDeployed: now,
					},
					Config: expiredDat,
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "second test expired - not expired",
			args: args{
				r: release.Release{
					Name: "release",
					Info: &release.Info{
						LastDeployed: now,
					},
					Config: futureDate,
				},
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckReleaseExpired(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckReleaseExpired() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckReleaseExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}