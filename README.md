# Provider API устройств Openhab в Yandex-dialogs

`./ohyandex` собран под AMD \
`GOOS=linux && GOARCH=amd && go build -o ohyandex cmd/ohyandex/main.go`

<h3>.env.local</h3>

CLIENT_ID и CLIENT_SECRET - будут использоваться Навыком Яндекс для oauth авторизации в ohyandex приложении \
OPENHAB_HOST - хост Openhab \

<h3>Установка</h3>

`cd /opt` \
`git clone https://github.com/ebuyan/ohyandex.git` \
`cp .env .env.local` \
`mkdir -p /var/log/ohyandex` \
`touch /var/log/ohyandex/app.log` \
`cp ohyandex.service /etc/systemd/systemd` \
`systemctl daemon-reload` \
`systemctl start ohyandex.service` \
`systemctl enable ohyandex.service`