package tcestuary

import (
	"reflect"
	"testing"

	"git.code.oa.com/tce-config/tcestuary-go/v4/configcenter"
)

func Test_parseTSMSecretConfig(t *testing.T) {
	tests := []struct {
		name    string
		want    configcenter.TSMConfig
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTSMSecretConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseTSMSecretConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTSMSecretConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitTencentSM(t *testing.T) {
	type args struct {
		appid  []byte
		bundle []byte
		cert   []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "init 1",
			args: args{
				appid:  []byte("tencentSmWeb.com.cn"),
				bundle: nil,
				cert: []byte(`-----BEGIN CERTIFICATE-----
MIICqjCCAlGgAwIBAgIJMTAwMDAwMDE2MAoGCCqBHM9VAYN1MIHAMQswCQYDVQQG
EwJDTjESMBAGA1UECAwJR3Vhbmdkb25nMREwDwYDVQQHDAhTaGVuemhlbjE6MDgG
A1UECgwxU2hlbnpoZW4gVGVuY2VudCBDb21wdXRlciBTeXN0ZW1zIENvbXBhbnkg
TGltaXRlZDEMMAoGA1UECwwDRmlUMRowGAYDVQQDDBFUZW5jZW50U00gUm9vdCBD
QTEkMCIGCSqGSIb3DQEJARYVVGVuY2VudFNNQHRlbmNlbnQuY29tMCIYDzIwMjAx
MjA5MDkxODI1WhgPMjAyMTAzMTkwOTE4MjVaMIGNMQswCQYDVQQGEwJDTjESMBAG
A1UECAwJR3VhbmdEb25nMREwDwYDVQQHDAhTaGVuWmhlbjEVMBMGA1UECgwMVGVu
Y2VudCBJbmMuMQwwCgYDVQQLDANmaXQxFDASBgNVBAMMC1RTTSBsaWNlbnNlMRww
GgYDVQQNDBN0ZW5jZW50U21XZWIuY29tLmNuMFkwEwYHKoZIzj0CAQYIKoEcz1UB
gi0DQgAEmVuMCmTg8d6E8erM1BfNOm4YmLPvjLumBjXIuhMvqVam1HHc/kBZam6Y
hSHa7uGa1SdF6aGKmvV1OZdpockskqNhMF8wHwYDVR0jBBgwFoAUoHUwm/JhjsLU
0gyxBmEGzO3vkAwwHQYDVR0OBBYEFHNC+Ci8Mgam8cAaDhn4hrauyKAeMAwGA1Ud
EwEB/wQCMAAwDwYDVR0PAQH/BAUDAwc4ADAKBggqgRzPVQGDdQNHADBEAiAdJTC+
6faaT3SAfAJZZ9DUDu7FrdD4WKb8rT9rcUkc8wIgbPa40xE7lF9RIUe3ZoBH/ibT
fqAFHhS63CWKxWvsExw=
-----END CERTIFICATE-----`),
			},
			wantErr: false,
		},
		{
			name: "init 2",
			args: args{
				appid:  []byte("tencentSmWeb.com.cn"),
				bundle: nil,
				cert: []byte(`-----BEGIN CERTIFICATE-----
MIICqjCCAlGgAwIBAgIJMTAwMDAwMDE2MAoGCCqBHM9VAYN1MIHAMQswCQYDVQQG
EwJDTjESMBAGA1UECAwJR3Vhbmdkb25nMREwDwYDVQQHDAhTaGVuemhlbjE6MDgG
A1UECgwxU2hlbnpoZW4gVGVuY2VudCBDb21wdXRlciBTeXN0ZW1zIENvbXBhbnkg
TGltaXRlZDEMMAoGA1UECwwDRmlUMRowGAYDVQQDDBFUZW5jZW50U00gUm9vdCBD
QTEkMCIGCSqGSIb3DQEJARYVVGVuY2VudFNNQHRlbmNlbnQuY29tMCIYDzIwMjAx
MjA5MDkxODI1WhgPMjAyMTAzMTkwOTE4MjVaMIGNMQswCQYDVQQGEwJDTjESMBAG
A1UECAwJR3VhbmdEb25nMREwDwYDVQQHDAhTaGVuWmhlbjEVMBMGA1UECgwMVGVu
Y2VudCBJbmMuMQwwCgYDVQQLDANmaXQxFDASBgNVBAMMC1RTTSBsaWNlbnNlMRww
GgYDVQQNDBN0ZW5jZW50U21XZWIuY29tLmNuMFkwEwYHKoZIzj0CAQYIKoEcz1UB
gi0DQgAEmVuMCmTg8d6E8erM1BfNOm4YmLPvjLumBjXIuhMvqVam1HHc/kBZam6Y
hSHa7uGa1SdF6aGKmvV1OZdpockskqNhMF8wHwYDVR0jBBgwFoAUoHUwm/JhjsLU
0gyxBmEGzO3vkAwwHQYDVR0OBBYEFHNC+Ci8Mgam8cAaDhn4hrauyKAeMAwGA1Ud
EwEB/wQCMAAwDwYDVR0PAQH/BAUDAwc4ADAKBggqgRzPVQGDdQNHADBEAiAdJTC+
6faaT3SAfAJZZ9DUDu7FrdD4WKb8rT9rcUkc8wIgbPa40xE7lF9RIUe3ZoBH/ibT
fqAFHhS63CWKxWvsExw=
-----END CERTIFICATE-----`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := initTencentSM(tt.args.appid, tt.args.bundle, tt.args.cert); (err != nil) != tt.wantErr {
				t.Errorf("initTencentSM() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInitTencentSMWithConfig(t *testing.T) {
	type args struct {
		tsmConf configcenter.TSMConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitTencentSMWithConfig(tt.args.tsmConf); (err != nil) != tt.wantErr {
				t.Errorf("InitTencentSMWithConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
