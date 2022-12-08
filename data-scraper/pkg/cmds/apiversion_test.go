package cmds

import (
	"reflect"
	"testing"
)

func TestParseAPIVersion(t *testing.T) {
	type args struct {
		version string
	}
	tests := []struct {
		name    string
		args    args
		want    *APIVersion
		wantErr bool
	}{
		{
			name:    "v1",
			args:    args{version: "v1"},
			want:    &APIVersion{Major: 1, Level: "stable", Revision: 0},
			wantErr: false,
		}, {
			name:    "v2",
			args:    args{version: "v2"},
			want:    &APIVersion{Major: 2, Level: "stable", Revision: 0},
			wantErr: false,
		}, {
			name:    "v1alpha1",
			args:    args{version: "v1alpha1"},
			want:    &APIVersion{Major: 1, Level: "alpha", Revision: 1},
			wantErr: false,
		}, {
			name:    "v1alpha20",
			args:    args{version: "v1alpha20"},
			want:    &APIVersion{Major: 1, Level: "alpha", Revision: 20},
			wantErr: false,
		}, {
			name:    "v1beta1",
			args:    args{version: "v1beta1"},
			want:    &APIVersion{Major: 1, Level: "beta", Revision: 1},
			wantErr: false,
		}, {
			name:    "v1beta20",
			args:    args{version: "v1beta20"},
			want:    &APIVersion{Major: 1, Level: "beta", Revision: 20},
			wantErr: false,
		}, {
			name: "v2alpha3",
			args: args{version: "v2alpha3"},
			want: &APIVersion{Major: 2, Level: "alpha", Revision: 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAPIVersion(tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAPIVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseAPIVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIVersion_LessThan(t *testing.T) {
	type fields struct {
		Major    int
		Level    string
		Revision int
	}
	type args struct {
		another APIVersion
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{{
		name: "v1 < v2",
		fields: fields{
			Major:    1,
			Level:    "stable",
			Revision: 0,
		},
		args: args{
			another: APIVersion{
				Major:    2,
				Level:    "stable",
				Revision: 0,
			},
		},
		want: true,
	}, {
		name: "v1alpha1 < v1",
		fields: fields{
			Major:    1,
			Level:    "alpha",
			Revision: 1,
		},
		args: args{
			another: APIVersion{
				Major:    1,
				Level:    "stable",
				Revision: 0,
			},
		},
		want: true,
	}, {
		name: "v1alpha1 < v1alpha2",
		fields: fields{
			Major:    1,
			Level:    "alpha",
			Revision: 1,
		},
		args: args{
			another: APIVersion{
				Major:    1,
				Level:    "alpha",
				Revision: 2,
			},
		},
		want: true,
	}, {
		name: "v1alpha1 < v1beta1",
		fields: fields{
			Major:    1,
			Level:    "alpha",
			Revision: 1,
		},
		args: args{
			another: APIVersion{
				Major:    1,
				Level:    "beta",
				Revision: 1,
			},
		}, want: true,
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := APIVersion{
				Major:    tt.fields.Major,
				Level:    tt.fields.Level,
				Revision: tt.fields.Revision,
			}
			if got := v.LessThan(tt.args.another); got != tt.want {
				t.Errorf("LessThan() = %v, want %v", got, tt.want)
			}
		})
	}
}
