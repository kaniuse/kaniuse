package main

import (
	"github.com/Masterminds/semver"
	"testing"
)

func TestNewRepoSwaggerFetcher(t *testing.T) {
	type args struct {
		version string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "v1.24.0",
			args: args{
				version: "v1.24.0",
			},
			wantErr: true,
		}, {
			name: "1.24.3",
			args: args{
				version: "1.24.3",
			},
			wantErr: false,
		},
		{
			// equivalent to 1.24.0
			name: "1.24",
			args: args{
				version: "1.24",
			},
			wantErr: false,
		},
		{
			name: "1.24.0-alpha.1",
			args: args{
				version: "1.24.0-alpha.1",
			},
			wantErr: false,
		}, {
			name: "not a semver",
			args: args{
				version: "not a semver",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewRepoSwaggerFetcher(tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRepoSwaggerFetcher() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestRepoSwaggerFetcher_KubernetesMinorVersion(t *testing.T) {
	type fields struct {
		version *semver.Version
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "1.24.0",
			fields: fields{
				version: semver.MustParse("1.24.0"),
			},
			want:    "1.24",
			wantErr: false,
		}, {
			name: "1.24.3",
			fields: fields{
				version: semver.MustParse("1.24.3"),
			},
			want:    "1.24",
			wantErr: false,
		}, {
			name: "1.24.0-rc.0",
			fields: fields{
				version: semver.MustParse("1.24.0-rc.0"),
			},
			want:    "1.24",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RepoSwaggerFetcher{
				version: tt.fields.version,
			}
			got, err := r.KubernetesMinorVersion()
			if (err != nil) != tt.wantErr {
				t.Errorf("KubernetesMinorVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("KubernetesMinorVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}
