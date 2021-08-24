package delete

import (
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name string
		want *Client
	}{
		{
			name: "basic",
			want: &Client{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClient(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
