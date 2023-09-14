# go-musthave-devops-tpl

Шаблон репозитория для практического трека «Go в DevOps».

# Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` - адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

# Обновление шаблона

Чтобы получать обновления автотестов и других частей шаблона, выполните следующую команду:

```
git remote add -m main template https://github.com/yandex-praktikum/go-musthave-devops-tpl.git
```

Для обновления кода автотестов выполните команду:

(Для Unix систем)

```
git fetch template && git checkout template/main .github
```

(Для Windows PowerShell)

```
(git fetch template) -and (git checkout template/main .github)
```

Затем добавьте полученные изменения в свой репозиторий.

# Запуск godoc 

Чтобы запустить страницу с документацией воспользуйтесь утилитой godoc, выполните следующую кманду: 
```
godoc -http=:8090 -play
```
По умолчанию godoc не отображает пакеты, расположенные в поддиректориях internal. 
Чтобы увидеть служебные пакеты, добавьте в браузере параметр ?m=all: например, 
```
http://localhost:8090/pkg/?m=all
```
# Компиляция проекта
Для сборки серверной части проекта выполните команду следующего вида (vN.N.NN - версия сборки, например, v1.0.15)
```
	go build -ldflags "-X 'main.buildVersion=vN.N.NN' -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'  -X 'main.buildCommit=YOUR COMMIT TEXT'" cmd/server/main.go

```
Перед компиляцией проекта правильность заполнения параметров команды можно проверить запустив приложение, передав те же параметры. Например: 
```
go run -ldflags "-X 'main.buildVersion=v1.0.15' -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'  -X 'main.buildCommit=LATEST COMMIT'" cmd/server/main.go

 ```
В результате приложение выводит в stdout 

```
Build version: v1.0.15
Build date: 2023/09/11 00:22:46
Build commit: LATEST COMMIT

```
Для сборки клиентской части проекта соответственно
```
	go build -ldflags "-X 'main.buildVersion=vN.N.NN' -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'  -X 'main.buildCommit=YOUR COMMIT TEXT'" cmd/agent/main.go

```
