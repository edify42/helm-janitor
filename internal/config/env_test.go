package config

import (
	"testing"
)

func TestEnvConfig_setDefaultEnv(t *testing.T) {
	type fields struct {
		Cluster         string
		Region          string
		JanitorLabel    string
		TmpFileLocation string
		TmpFilePrefix   string
		DebugFlag       bool
	}
	type args struct {
		field string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "first test",
			fields: fields{
				Cluster:         "test-only",
				Region:          "ap-southeast-2",
				JanitorLabel:    "teddy=test",
				TmpFileLocation: "/tmp",
				TmpFilePrefix:   "fs-*",
				DebugFlag:       false,
			},
			args: args{
				field: "Cluster",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &EnvConfig{
				Cluster:         tt.fields.Cluster,
				Region:          tt.fields.Region,
				JanitorLabel:    tt.fields.JanitorLabel,
				TmpFileLocation: tt.fields.TmpFileLocation,
				TmpFilePrefix:   tt.fields.TmpFilePrefix,
				DebugFlag:       tt.fields.DebugFlag,
			}
			c.setDefaultEnv(tt.args.field)
		})
	}
}

func TestEnvConfig_setFromEnv(t *testing.T) {
	type fields struct {
		Cluster         string
		Region          string
		JanitorLabel    string
		TmpFileLocation string
		TmpFilePrefix   string
		DebugFlag       bool
	}
	type args struct {
		field string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "first test",
			fields: fields{
				Cluster:         "test-only",
				Region:          "ap-southeast-2",
				JanitorLabel:    "teddy=test",
				TmpFileLocation: "/tmp",
				TmpFilePrefix:   "fs-*",
				DebugFlag:       false,
			},
			args: args{
				field: "Cluster",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &EnvConfig{
				Cluster:         tt.fields.Cluster,
				Region:          tt.fields.Region,
				JanitorLabel:    tt.fields.JanitorLabel,
				TmpFileLocation: tt.fields.TmpFileLocation,
				TmpFilePrefix:   tt.fields.TmpFilePrefix,
				DebugFlag:       tt.fields.DebugFlag,
			}
			c.setFromEnv(tt.args.field)
		})
	}
}

func TestEnvConfig_Init(t *testing.T) {
	type fields struct {
		Cluster         string
		Region          string
		JanitorLabel    string
		TmpFileLocation string
		TmpFilePrefix   string
		DebugFlag       bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "first test",
			fields: fields{
				Cluster:         "test-only",
				Region:          "ap-southeast-2",
				JanitorLabel:    "teddy=test",
				TmpFileLocation: "/tmp",
				TmpFilePrefix:   "fs-*",
				DebugFlag:       false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &EnvConfig{
				Cluster:         tt.fields.Cluster,
				Region:          tt.fields.Region,
				JanitorLabel:    tt.fields.JanitorLabel,
				TmpFileLocation: tt.fields.TmpFileLocation,
				TmpFilePrefix:   tt.fields.TmpFilePrefix,
				DebugFlag:       tt.fields.DebugFlag,
			}
			c.Init()
		})
	}
}
