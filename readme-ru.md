# GophKeeper

![GophKeeper](https://pictures.s3.yandex.net/resources/gophkeeper_2x_1650456239.png)

## Описание

GophKeeper — это клиент-серверная система для безопасного хранения логинов, паролей и других приватных данных. Основной функционал сосредоточен в клиентской части, которая позволяет пользователям управлять своими данными. Клиентская часть может работать без синхронизации с сервером.

## Функциональные возможности

### Клиентская часть:
- **Аутентификация и авторизация**: Регистрация и вход в систему для доступа к данным.
- **Работа в автономном режиме**: Возможность управления данными локально без синхронизации с сервером.
- **Запрос данных по ключу**: Пользователи могут получать доступ к своим данным, запрашивая их по уникальному ключу.

### Серверная часть:
- Регистрация и аутентификация пользователей.
- Синхронизация данных между клиентами.
- Хранение и управление приватными данными.

## Установка

### Установка с использованием Makefile
1. Клонируйте репозиторий:  
   ```
   git clone <URL>  
   cd gophKeeper
   ```  
2. Убедитесь, что у вас установлен Makefile.
3. Запустите команды для сборки и запуска приложения:
   - Для сборки серверной части:
   ```
   make build_server
   ```
- Для сборки клиентской части:
  ```
  make build_client
  ```
- Для сборки обоих компонентов:
  ```
  make build_all
  ```
- Для сборки сервера:
  ```
  make run_app
  ```
- Для сборки клиента:
  ```
  make run_client
  ```
### Тестирование
- Для запуска тестов:
  ```
  make test
  ```
- Для проверки на наличие гонок данных:
  ```
  make race
  ```
- Для создания отчета о покрытии:
  ```
  make cover
  ```
## Архитектура

GophKeeper состоит из клиентской и серверной частей, реализованных на языке Go. Клиент предоставляет интерфейс командной строки для управления данными, а сервер использует gRPC для обработки запросов и синхронизации данных. Все данные шифруются для обеспечения безопасности.

## Лицензия

Этот проект лицензируется в соответствии с условиями [MIT License](LICENSE).