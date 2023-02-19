package httpeasy

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

var (
	validConfig = Config{
		Port:              8089,
		ReadTimeout:       1 * time.Minute,
		ReadHeaderTimeout: 15 * time.Second,
		WriteTimeout:      1 * time.Minute,
	}

	//portTooSmall = Config{
	//	SSL: &SSL{Email: "info@yandex.ru",
	//		Domain: "yandex.ru",
	//	},
	//	ProductionMode: true,
	//	HTTP: HTTP{
	//		Port: 50,
	//	},
	//}
	//
	//portTooBig = Config{
	//	SSL: &SSL{Email: "info@yandex.ru",
	//		Domain: "yandex.ru",
	//	},
	//	ProductionMode: true,
	//	HTTP: HTTP{
	//		Port: 50000,
	//	},
	//}
	//
	//badEmail = Config{
	//	SSL: &SSL{Email: "info",
	//		Domain: "yandex.ru",
	//	},
	//	ProductionMode: true,
	//	HTTP: HTTP{
	//		Port: 8089,
	//	},
	//}
	//
	//badDomain = Config{
	//	SSL: &SSL{Email: "info@yandex.ru",
	//		Domain: "-",
	//	},
	//	ProductionMode: true,
	//	HTTP: HTTP{
	//		Port: 8089,
	//	},
	//}
	//
	//devMode = Config{
	//	HTTP: HTTP{
	//		Port: 8089,
	//	},
	//}
	//
	//sslForgotten = Config{
	//	ProductionMode: true,
	//	HTTP:           HTTP{Port: 8089},
	//}
	//
	//dbForgotten = Config{
	//	HTTP:  HTTP{Port: 8089},
	//}
)

func TestConfig_newWebService(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		httpServer := validConfig.newWebService()
		if !strings.Contains(httpServer.Addr, fmt.Sprintf(":%d", validConfig.Port)) {
			t.Error("incorrect http port")
		}

		if httpServer.ReadTimeout != validConfig.ReadTimeout {
			t.Error("invalid read timeout")
		}

		if httpServer.WriteTimeout != validConfig.WriteTimeout {
			t.Error("invalid write timeout")
		}

		if httpServer.ReadHeaderTimeout != validConfig.ReadHeaderTimeout {
			t.Error("invalid read header timeout")
		}
	})
}
