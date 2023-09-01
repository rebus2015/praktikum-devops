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
godoc -http=:8088 -play
```
По умолчанию godoc не отображает пакеты, расположенные в поддиректориях internal. 
Чтобы увидеть служебные пакеты, добавьте в браузере параметр ?m=all: например, 
```
http://localhost:8088/pkg/?m=all
```
