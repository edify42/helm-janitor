package scan

import "testing"

func TestValidateScanArg(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "failing test - no k=v pairing",
			args:    args{[]string{"test"}},
			want:    false,
			wantErr: true,
		},
		{
			name:    "valid test - empty string array",
			args:    args{[]string{}},
			want:    false,
			wantErr: false,
		},
		{
			name:    "valid test - one k=v pair",
			args:    args{[]string{"label=yo"}},
			want:    true,
			wantErr: false,
		},
		{
			name:    "failing test - too much k=v pairing",
			args:    args{[]string{"aoeu", "aoeu", "aoeu"}},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateScanArg(tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateScanArg() error %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateScanArg() = %v, want %v", got, tt.want)
			}
		})
	}
}
