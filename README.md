# Twitch bot

This is a test bot, using oauth.

## how to run

create a file called `.env` with the following contents:

```
TWITCH_CLIENT_ID=<id>
TWITCH_CLIENT_SECRET=<secret>
TWITCH_PREFIX=<prefix>
TWITCH_CHANNEL=<channel>
```

TWITCH_CLIENT_ID and TWITCH_CLIENT_SECRET you can get from twitch developer console by creating an app.

TWITCH_PREFIX is the command the bot should listen to.

TWITCH_CHANNEL is the channel to which the bot should connect to.

then:

```sh
go run .

# or

go build .
./bot-ttv
```
