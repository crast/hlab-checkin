# Anime game automated checkin


## Build it
```shell
go build .
```

## Run it


Log in on firefox or edge (chrome/chromium uses encrypted cookies, may not work)
Alternately edit `internal/getcookie/get_cookie.go` and comment out browsers you don't want

You can also get cookies from a cookies.txt style file (export from browser of choice, or from apps like curl).
To do this, the `-file` option parameter must end in `.txt`


```shell
# Can run for each game
./hlab-checkin -game genshin
./hlab-checkin -game honkai
./hlab-checkin -game zzz
```

(Windows users may need to do `.\hlab-checkin.exe` )

You can also use `-file` to store multiple account profiles.  

automation and GUI maybe in the future
