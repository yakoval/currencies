package config

import (
	"gitlab.com/knopkalab/gobox"
)

// Config приложения
type Config struct {
	// Конфигурация загрузчика данных.
	Updater UpdaterConfig

	// Конфигурация БД.
	Database DatabaseConfig

	// Конфигурация REST-сервера.
	Web WebConfig
}

// Конфигурация загрузчика данных.
type UpdaterConfig struct {
	// Адрес, по которому располагаются даннные.
	URI string
}

// Конфигурация БД.
type DatabaseConfig struct {
	// Хост.
	Host string

	// Порт.
	Port uint16

	// Логин.
	Login string

	// Пароль.
	Password string

	// Имя БД.
	Name string
}

// Конфигурация rest-сервера.
type WebConfig struct {
	// IP-адрес с портом, по которому доступен сервер.
	Host string

	// Если true, вывод осуществляется постранично.
	ByPages bool

	// Количество отображаемых элементов на 1 странице списка валют.
	ItemsPerPage int
}

// NewConfig - создание новой конфигурации приложения на основе файла .env
func NewConfig() *Config {
	env := gobox.NewEnvParser(".env")

	c := new(Config)

	c.Updater.URI = env.Get("UPDATER_URI", "http://www.cbr.ru/scripts/XML_daily.asp")

	c.Database.Host = env.Get("DB_HOST", "127.0.0.1")
	c.Database.Port = env.Port("DB_PORT", 3306)
	c.Database.Login = env.Get("DB_LOGIN", "root")
	c.Database.Password = env.Get("DB_PASSWORD", "")
	c.Database.Name = env.Get("DB_NAME", "currency")

	c.Web.Host = env.Get("WEB_SERVER_HOST", "127.0.0.1:8081")
	c.Web.ByPages = env.Bool("REST_BY_PAGES", true)
	c.Web.ItemsPerPage = env.Int("REST_ITEMS_PER_PAGE", 10)

	return c
}
