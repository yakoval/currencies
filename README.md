# Currency viewer

Обновление валют и их просмотр

## Установка

### MySQL
Для работы необходимо создать базу данных 
```Mysql
CREATE DATABASE IF NOT EXISTS `immo-test` /*! возможно другое имя*/;
USE `immo-test`;

CREATE TABLE IF NOT EXISTS `cache` (
  `count_total` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `currency` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) DEFAULT NULL,
  `rate` decimal(10,4) DEFAULT NULL,
  `insert_dt` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
```
### Конфигурирование
Приложение конфигурируется при помощи файла `.env`, расположенного в корне проекта.
Его переменные:

Переменная | Тип | Назначение | Значение по умолчанию
--- | --- | --- | --- 
UPDATER_URI|Строка|адрес истоника данных|http://www.cbr.ru/scripts/XML_daily.asp
DB_HOST|IP-адрес|адрес хоста используемой БД |127.0.0.1
DB_PORT|Целое число|порт используемой БД|3306
DB_LOGIN|Строка|логин, для доступа к БД|root
DB_PASSWORD|Строка|пароль для доступа к БД|
DB_NAME|Строка|имя БД|immo-test
WEB_SERVER_HOST|IP-адрес|адрес хоста, на котором будет работать REST-сервер|127.0.0.1:8081
REST_BY_PAGES|Bool|если `true`, делалть ли вывод постраничным.|true
REST_ITEMS_PER_PAGE|Целое число|количество элементов на страницу|10

## Использование
### Обновление валют
Под linux делается сборка 
```bash
cd cmd/update
go build -o updater
```
и запускается
```bash
./updater
```

### REST-сервер
Под linux делается сборка 
```bash
cd cmd/rest
go build -o rest
```
и запускается
```bash
./rest
```

### Запросы к REST
Для постраничной навигации параметр конфигнурации `REST_BY_PAGES` должен быть `true`.

Навигация осуществляется путём передачи параметра `page` в строке запроса, например:
```
http://127.0.0.1:8081/currencies?page=2
```