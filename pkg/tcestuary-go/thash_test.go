package tcestuary

import (
	"reflect"
	"testing"

	"git.code.oa.com/tce-config/tcestuary-go/v4/configcenter"
	"git.code.oa.com/tce-config/tcestuary-go/v4/tcesecurity"
)

func Test_thasher_New(t *testing.T) {
	type fields struct {
		f tcesecurity.HashFunc
	}
	tests := []struct {
		name    string
		fields  fields
		want    THash
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &thasher{
				f: tt.fields.f,
			}
			got, err := h.New()
			if (err != nil) != tt.wantErr {
				t.Errorf("thasher.New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("thasher.New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseHashSecretConfig(t *testing.T) {
	tests := []struct {
		name    string
		want    configcenter.HashConfig
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseHashSecretConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseHashSecretConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseHashSecretConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTHasher(t *testing.T) {
	tests := []struct {
		name    string
		want    THasher
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTHasher()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTHasher() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTHasher() = %v, want %v", got, tt.want)
			}
		})
	}
}
