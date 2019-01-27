# fsm
Ridiculously hacky tool to help me keep in touch with my friends.

## Problem

Interacting with my friends makes me happy, but I have an inertial tendency to lose touch. Also, I live away from 99% of them at the moment, so it's hard.

## Solution

I hacked together a Google Sheet with my list of friends, the frequency I want to talk to them, and a chat link. Then I wrote this go script that reads the sheet, calculates if the frequencies are met or not, choses a random friend out of the matched ones and opens the chat link. There should be a cron on the system that runs every now and then, but one should also be able to run it manually.

## Chat links look like:

- https://www.facebook.com/messages/t/%name_or_id%
- whatsapp://send?text=Time%20to%20talk!&phone=%number%
- https://mail.google.com/mail/?view=cm&to=%email%

## Hacky things that make it hacky

- I just call "open %chat_link%". Beware of what's in that Sheet cell, open doesn't work for everyone, have to be logged in to FB, WhatsApp Web and Gmail.
- I call `Get` and `Update` on the Sheets API. For `Get` I read the whole sheet (btw Sheet name hardcoded). For `Update` I calculate the cell to write extremely hackily.
- In order to make Google Sheets API work I created a service account that has access to that particular sheet and downloaded a `credentials.json`, so this is not a `git clone` away from working.
- The auth permission I created is ridiculously permissive. I think it's like "this app can delete all your spreadsheets".

## Usage

- Run it in your crontab and let it annoy you
- Run it directly i.e. `fsm` and annoy yourself, chat-roulette style
- Run `fsm sheet` to open the Sheet if you want to make adjustments
- Run it with a hardcoded person e.g. `fsm chat "John Doe"`

## If you want to use it

The implementation is terrible, but it works quite well for me at the moment. If you want to try it out feel free to contact me via email or create an issue.

It's not super hard to setup: 

- `go get -u github.com/marianogappa/fsm`
- need Google Cloud account (there's a 1-year free tier)
- create a Service Account
- create a Google Sheet with a table with `name` (e.g. John Doe), `chat_link` (e.g. https://mail.google.com/mail/?view=cm&to=johndoe@gmail.com), `frequency` (e.g. 30 means 30 days) and `last_comm` (cli will update so leave blank but e.g. 2019-01-27).
- create credentials for the service account to update the sheet; save as credentials.json in the root of the code where `main.go` is.
- Update code hardcoded things in `main.go` like spreadsheet id, sheet id, my repo address if you forked it, etc.
- Check to see if it works with `fsm`. It won't. Debug it or ask me.
- Add to crontab so it actually forces you to be happy. See crontab example below.

## crontab

```
0 * * * * export GOPATH=%your_gopath% && $GOPATH/bin/fsm
```
