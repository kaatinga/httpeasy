package QuickHTTPServerLauncher

import (
	"testing"
)

func TestConfig_check(t *testing.T) {

	validConfig := Config{
		email:      "info@yandex.ru",
		launchMode: "prod",
		port:       "8089",
		domain:     "yandex.ru",
	}

	portTooSmall := Config{
		email:      "info@yandex.ru",
		launchMode: "prod",
		port:       "500",
		domain:     "yandex.ru",
	}

	portTooBig := Config{
		email:      "info@yandex.ru",
		launchMode: "prod",
		port:       "50000",
		domain:     "yandex.ru",
	}

	badEmail := Config{
		email:      "info",
		launchMode: "prod",
		port:       "5000",
		domain:     "yandex.ru",
	}

	badDomain := Config{
		email:      "info@yandex.ru",
		launchMode: "prod",
		port:       "5000",
		domain:     "-",
	}

	badMode := Config{
		email:      "info@yandex.ru",
		launchMode: "test",
		port:       "5000",
		domain:     "yandex.ru",
	}

	tests := []struct {
		name    string
		fields  Config
		wantErr bool
	}{
		{"ok", validConfig, false},
		{"port is too small", portTooSmall, true},
		{"port is too big", portTooBig, true},
		{"bad email", badEmail, true},
		{"bad domain", badDomain, true},
		{"bad launch mode", badMode, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := tt.fields
			if err := config.check(); (err != nil) != tt.wantErr {
				t.Errorf("check() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
