# astra-test-task
### Настройка переменных окружения. В корне проекта есть `.env` файл со значениями, которые использовала. 
Переменные окружения необходимо настроить для базы данных.  
```text
POSTGRES_USERNAME = <имя_пользователя_для_подключения_к_бд>
POSTGRES_PASSWORD = <пароль_для_подключения_кь_бд>
POSTGRES_HOST = <хост_для_подключения_к_бд>
POSTGRES_PORT = <порт_для_подключения_к_бд>
POSTGRES_DATABASE = <имя_бд>
POSTGRES_CONN_URL = <URL-строка_для_подключения_к_бд>
```

### База данных 
БД security PostgreSQL создана локально:   
```text
createdb -U postgres security;```
```  
Создание таблицы warnings происходит при выполнении программы.

### Запуск приложения
#### Запуск без сборки
```go
go run main.go report.sarif
```

#### с помощью исполняемого файла
Чтобы скомпилировать исполняемый файл для конфигурации системы GOOS=linux GOARCH=amd64 использованы команды:
```text
set GOARCH=amd64
set GOOS=linux
go build -o bin/app-amd64-linux main.go
```  
Скомпилирвоанный файл: ```astra-test-task/bin/app-amd64-linux```

Для запуска используйте команду:
```text
./bin/app-amd64-linux report.sarif
```
#### через docker (не доделала :\ )
Чтобы собрать образ:
```text
docker build --tag astra-test .
```
Запустить (выдает ошибку failed to connect to `user=postgres database=security`):
```text
docker run -e POSTGRES_USERNAME=postgres -e POSTGRES_PASSWORD=postgres1 -e POSTGRES_HOST=localhost -e POSTGRES_PORT=5432 -e POSTGRES_DATABASE=security astra-test
```
*todo: в идеале разобраться и использовать docker-compose

### Проверка есть ли в программе уязвимости
Я использовала GoSec в составе [golangci-lint](https://github.com/golangci/golangci-lint)  

Установка golangci-lint:  
```curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.61.0```  

Запуск с включенным линтером GoSec:  
```golangci-lint run -E gosec```  
Проблемы с безопасностью не выявлены.

### Результаты
1. С помощью парсинга данные получены из json файла и занесены в БД в созданную таблицу. 
2. По результатам парсинга файла в консоль выводится классификация
количества уязвимостей по критичности (с примером для проверки корректности совпадает).

Что заметила: при повторном запуске приложения с тем же json файлом, переданным в качестве аргумента, данные будут опять занесены в таблицу, т.е. появятся дубликаты записей. Это логично, поскольку нет каких-то ограничений уникальности. Я думала, нужно ли с этим что-то делать, но решила, что лучше потом уточнить, необходимо ли решать этот вопрос или нет. Но например, можно изменить таблицу добавив уникальный ключ ```UNIQUE (ruleId, uri, startLine)```. 