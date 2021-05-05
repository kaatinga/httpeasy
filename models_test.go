package QuickHTTPServerLauncher

import (
	"errors"
	"github.com/davecgh/go-spew/spew"
	"testing"
)

var (
	validConfig = Config{
		SSL: &SSL{Email: "info@yandex.ru",
			Domain: "yandex.ru",
		},
		ProductionMode: true,
		HTTP: HTTP{
			Port: 8089,
		},
	}

	portTooSmall = Config{
		SSL: &SSL{Email: "info@yandex.ru",
			Domain: "yandex.ru",
		},
		ProductionMode: true,
		HTTP: HTTP{
			Port: 50,
		},
	}

	portTooBig = Config{
		SSL: &SSL{Email: "info@yandex.ru",
			Domain: "yandex.ru",
		},
		ProductionMode: true,
		HTTP: HTTP{
			Port: 50000,
		},
	}

	badEmail = Config{
		SSL: &SSL{Email: "info",
			Domain: "yandex.ru",
		},
		ProductionMode: true,
		HTTP: HTTP{
			Port: 8089,
		},
	}

	badDomain = Config{
		SSL: &SSL{Email: "info@yandex.ru",
			Domain: "-",
		},
		ProductionMode: true,
		HTTP: HTTP{
			Port: 8089,
		},
	}

	devMode = Config{
		ProductionMode: false,
		HTTP: HTTP{
			Port: 8089,
		},
	}
)

func TestHTTP_check(t *testing.T) {

	tests := []struct {
		name    string
		fields  Config
		wantErr error
	}{
		{"ok1", devMode, nil},
		{"ok2", validConfig, nil},
		{"Port is too small", portTooSmall, errValidationError},
		{"Port is too big", portTooBig, errValidationError},
		{"bad Email", badEmail, errValidationError},
		{"bad Domain", badDomain, errValidationError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.check()
			if err != nil {
				t.Log(err)
			}
			if !errors.Is(err, tt.wantErr) {
				spew.Dump(tt.fields)
				t.Errorf("check() error\nhave %v\nwant %v\n", err, tt.wantErr)
			}
		})
	}
}
