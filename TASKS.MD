# Task List

---

## Задачи по инфраструктуре
* <em>~~Сделать более информативные и полезные логи, убрать
весь мусор оставшийся от поиска ошибок. Покрыть ключевые показатели логами(сообщение от пользователя,
отправленное уведомление)~~</em>
* <em>Prometheus + Grafana(Future)</em>
* Написать sh скрипт для деплоя - должен включать в себя git pull + docker-compose down && docker-compose up --build -d


---

## Задачи по коду

* <em>Доработка админ команды SET: бывает ситуация, 
когда у пользователя нет username, тогда предусмотреть варианты 
взаимодействия с таким пользователем, через номер телефона(следовательно, 
в БД под это надо добавить поле phone_number), либо еще какие-то варианты
 -- посмотреть в документации Telegram </em>
* <em>~~Доработка админ команды DELETE: в handleDelete предусмотреть когда ник пишут через @
в таком случае нужно эту @ убрать, потому что в бд username хранится всегда без @~~</em>
* <em>~~Вынести всю информацию об психологе в константы(например, 
время психолога, текст уведомления)~~</em>
* <em>Добавить вариацию текста уведомлений, когда они приходят за час, или же за 24 часа. Должны быть разные сообщения</em>

---

## Задачи по архитектуре

* <em>Внедрить Dependency Injection с помощью интерфейсов, после этого покрыть
 50% кода тестами</em>