# Bookings and Reservations System

This is a repository for my bookings and reservation project.

- Biult in Go version 1.16
- Uses the [chi routing](https://github.com/go-chi/chi)
- Uses [Alex Edwards SCS session](https://github.com/alexedwards/scs/v2) management
- Uses [nosurf](https://github.com/justinas/nosurf)
- Uses [go-simple-mail](https://github.com/xhit/go-simple-mail) & mailhog
- Uses Postgres/DBeaver
- Uses [Buffalo(soda)](https://gobuffalo.io/en/docs/db/migrations) for table migrations

</br>

migrate table
```bash=
soda migrate
```

</br>

start mail server
```bash=
brew services start mailhog
```

</br>

start web server
```bash=
go build -o bookings cmd/web/*.go
./bookings -dbhost=localhost -dbname= -dbuser= -dbport= -cache= -production=
```