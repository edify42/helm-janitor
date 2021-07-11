package utils

import (
	"reflect"
	"testing"
)

func TestParseTime(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    Time
		wantErr bool
	}{
		{
			name: "first valid test - 1 day",
			args: args{
				input: "1d",
			},
			want: Time{
				Days:    1,
				Hours:   24,
				Minutes: 1440,
				Seconds: 86400,
			},
			wantErr: false,
		},
		{
			name: "second valid test - 1 hour",
			args: args{
				input: "1h",
			},
			want: Time{
				Days:    0,
				Hours:   1,
				Minutes: 60,
				Seconds: 3600,
			},
			wantErr: false,
		},
		{
			name: "third valid test - 1 minute",
			args: args{
				input: "1m",
			},
			want: Time{
				Days:    0,
				Hours:   0,
				Minutes: 1,
				Seconds: 60,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTime(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
