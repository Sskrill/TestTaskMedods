# В этом проекте есть обработчики на маршрутах 
* auth/sign-up регистрация
* auth/sign-in аутентификация
* auth/refresh рефреш токенов (на этом обработчике если зайти с другого ip будет поссылаться письмо на почту которую регистрировались)
* medods/guid/{guid} взять токены по специальному guid у пользователя
  

*Код частично покрыт тестами и запускаеться в docker*
