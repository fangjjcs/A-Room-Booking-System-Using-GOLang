
# A Room Booking System Using GOLang

<img width="1279" alt="screenshot" src="https://user-images.githubusercontent.com/33279791/126778592-c3664ae0-2436-468e-bee0-c66583e59d0d.png">


This is a repository for my bookings and reservation project.

- Biult in Go version 1.16
- Uses the [chi routing](https://github.com/go-chi/chi)
- Uses [Alex Edwards SCS session](https://github.com/alexedwards/scs/v2) management
- Uses [nosurf](https://github.com/justinas/nosurf)
- Uses [go-simple-mail](https://github.com/xhit/go-simple-mail) & mailhog
- Uses Postgres/DBeaver
- Uses [Buffalo(soda)](https://gobuffalo.io/en/docs/db/migrations) for table migrations

</br>

#### demo vedio
https://drive.google.com/file/d/12pndL436i3igD9vb0p-VOoLtduYIvoiY/view?usp=sharing

#### Client
<img width="1000" alt="screenshot" src="https://user-images.githubusercontent.com/33279791/127741984-4a510179-0dcd-4ed6-8311-1b1671d85c30.png">

#### Admin Back-end
<img width="1000" alt="screenshot" src="https://user-images.githubusercontent.com/33279791/127741878-01c9c2ed-3533-437c-8c42-51106b996248.png">

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
