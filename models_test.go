package QuickHTTPServerLauncher

import (
	"testing"
)

func TestConfig_check(t *testing.T) {

	tests := []struct {
		name    string
		fields  Config
		wantErr bool
	}{
		{"ok", Config{domain: "google.com"}, false},
		{"!ok1", Config{domain: "--2"}, true},
		{"!ok2", Config{domain: ""}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				email:      tt.fields.email,
				db:         tt.fields.db,
				launchMode: tt.fields.launchMode,
				port:       tt.fields.port,
				domain:     tt.fields.domain,
				Logger:     tt.fields.Logger,
			}
			if err := config.check(); (err != nil) != tt.wantErr {
				t.Errorf("check() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
