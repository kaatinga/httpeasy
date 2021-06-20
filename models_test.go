package httpeasy

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
		HTTP: HTTP{
			Port: 8089,
		},
	}

	sslForgotten = Config{
		ProductionMode: true,
		HTTP:           HTTP{Port: 8089},
	}

	dbForgotten = Config{
		HasDB: true,
		HTTP:  HTTP{Port: 8089},
	}
)
