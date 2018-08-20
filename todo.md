# Сервер


## Архитектура сервисов

Начитавшись [хабра](https://habrahabr.ru/company/mailru/blog/220359/) пришёл к выводу, что чат с игровой механикой не стоит держать в одном месте. Более того - это скорее даже мешает - всё валится в одну кучу. Также авторизация остаётся незакрытым вопросом. Пожалуй стоит оформить каждый из этих фрагментов как отдельный сервис.

Сейчас сервер - это набор сервисов, которые могут обмениваться друг с другом сообщениями через брокер. На данный момент есть следующие сервисы:
- пул соединений (принимает соединения через TCP и WS, парсит передаваемые данные и пробрасывает эти сообщения в брокер)
- авторизация
- чат
- игровая логика
- _статистика?_

Каждый сервис отправляет брокеру сообщения, а тот их отправляет куда следует. На данный момент предполагается следующее поведение:
- клиент подключается к пулу соединений и отправляет запрос на вторизацию
- пул пробрасывает запрос в сервис авторизации
- сервис авторизации пока ничего не проверяет - он **создаёт нового пользователя в системе и делает 2 вещи**: 
    - отправляет **ответ обратно в пул соединений**, чтобы можно было что-то ответить пользователю
    - отправляет **уведомление в игровую логику** о присоединении нового игрока
- соответственно логика, как только получает уведомление - добавляет этого пользователя в список, создаёт под него новый объект и отправляет пользователю "welcome!"

## Брокер
~~Надо приделать ему поддержку fiber или как его там - добавлять в service message requestId и отдавать ответ на это сообщение с таким же requestId~~ (а может и не надо.. по идее на данный момент логика задумывается совсем простая и никаких запросов клиент присылать не будет).

## Логика

Симуляция производится строго по шагам с фиксированным периодом. В идеале стараемся, чтобы шаги симуляции совпадали с реальынм временем. Все приходящие от пользователей сообщения клаёдм в массив с метками времени, в которое они пришли. Этот массив потом должен позволить "воспроизвести" симуляцию. Надо бы приделать сохранение всего стейта симуляции, чтобы потом его можно было загрузить и применить пользовательский ввод и проверить, что всё точно считается. 

## Пул соединений

Пользователь в системе может существовать в разных контекстах. Пул соединений нумерует по порядку всех пользователей и сообщения им отправляются именно по этим id. У авторизации пользователи должны быть свои, но с привязкой к id пула.
У логики пользователи тоже свои, но с привязкой к id пула и ссылкой на пользователя в авторизации. По идее это всё должны быть прямо разные структуры данных, чтобы избежать concurrent modification в случае чего

## Авторизация




## Заметки

[Проблема с рефлексией](http://play.golang.org/p/AlQ9rOdXJU)
    > Объяснил добрый дядя на stackoverflow
  
Instead of sending throw broker we can pass client description with it's connection straight into logic - less overhead for sending.
Users should be registered globally in the server with all info (including connection). Each service can send info independently (chat doesn't collide with logic).
  

    MAIN STEP
        get inputs from queue and apply to world objects
        simulate N times
            wide phase
                foreach activeObject - find possible collisions (AABB)
            foreach possiblyCollidingObject
                check collision
            before collision
            resolve collisions
            after collision?
            update state
            should we remember objects in viewport in that place? (possibly less checks for distance)
    
        foreach player: 
            find objects in viewport (static objects are easy to find - in quadtree, what to do with active objects?)
            count new state
            count diff
            send diff (or full state if old)
    
    
    Logic main loop 
        wait for inputs from players
            put them into queue
        wait for simulation time
            MAIN STEP
    
    
    sending info to client
        add to send buffer
        send averything from buffer
        on received acknowledge - remove from buffer
        on diff lifetime expire - remove old from buffer