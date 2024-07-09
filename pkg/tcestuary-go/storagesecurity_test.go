package tcestuary

import (
	"reflect"
	"testing"

	"git.code.oa.com/tce-config/tcestuary-go/v4/configcenter"
)

func Test_parseStorageSecretConfig(t *testing.T) {
	tests := []struct {
		name    string
		want    configcenter.SecretConfig
		wantErr bool
	}{
		{
			name: "parse storage-secret",
			want: configcenter.SecretConfig{
				Method:     "aes-256-gcm",
				PublicKey:  "",
				PrivateKey: "",
				AesKey:     "5c2bd12683ceefb8830abba988339e67",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseStorageSecretConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseStorageSecretConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseStorageSecretConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
