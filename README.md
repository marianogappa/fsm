# fsm
Ridiculously hacky tool to help me keep in touch with my friends. `fsm` means "Friend Sheet Manager" btw.

## Problem

Interacting with my friends makes me happy, but I have an inertial tendency to lose touch. Also, I live away from 99% of them at the moment, so it's hard.

## Solution

I hacked together a Google Sheet with my list of friends (like [this one](https://docs.google.com/spreadsheets/d/15a5P0xrPdOwuxhYpBqTfIqvhjz3DJKy8fyVa-1cykgE/edit#gid=0)), the frequency I want to talk to them, and a chat link. Then I wrote this go script that reads the sheet, calculates if the frequencies are met or not, choses a random friend out of the matched ones and opens the chat link. There should be a cron on the system that runs every now and then, but one should also be able to run it manually.

## Chat links look like

- https://www.facebook.com/messages/t/%name_or_id%
- whatsapp://send?text=Time%20to%20talk!&phone=%number%
- https://mail.google.com/mail/?view=cm&to=%email%

Those links work on web as long as you're signed in. If you're gonna add a lot of people, I recommend some hackery to scrape the ids/addresses. I went to my `FB Profile > Friends`, loaded all of my friends and then played around in the Dev Console until I got all the chat links.

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

The implementation code is terrible, but it works quite well for me at the moment and it's surprisingly effective. If you want to try it out feel free to contact me via email or create an issue. I get bonus happiness points by helping other people be happier.

It's not super hard to setup, but it won't take less than 10 minutes: 

- `go get -u github.com/marianogappa/fsm`
- need Google Cloud account (there's a 1-year free tier: https://cloud.google.com/gcp)
- create a Service Account (https://cloud.google.com/iam/docs/creating-managing-service-accounts)
- create a Google Sheet like this one: https://docs.google.com/spreadsheets/d/15a5P0xrPdOwuxhYpBqTfIqvhjz3DJKy8fyVa-1cykgE/edit#gid=0
- create credentials for the service account to update the sheet; save as credentials.json in the root of the code where `main.go` is. Follow this guide: https://developers.google.com/sheets/api/guides/authorizing
- Update hardcoded things in `main.go` like spreadsheet id, sheet id, my repo address if you forked it, etc.
- Check to see if it works with `fsm`. It won't. Debug it or ask me.
- Add to crontab so it actually forces you to be happy. See crontab example below.
- Clap along. Because you're happy.

## crontab

```
0 * * * * export GOPATH=%your_gopath% && $GOPATH/bin/fsm
```

On the latest Mac OS `crontab -e` doesn't work anymore unless you give your terminal `Full Disk Access`. `System Preferences > Security & Privacy > Full Disk Access`.

## "Why not just use a local file for storage" chat I had

```
- wtf! Why didn't you just use a local file for storage?!
- Huh... -stares with dude face for a while- ...but like...Sheets is cool.
- It would have been so easy to set up! Wouldn't you have helped people be happy 100x more?!
- Umm...yeah but like... -stares with dude face for a while- ... yeah fair enough.
- Your architectural decision-making process gave me cancer.
- But like -stares with dude face for a while
*Leaves*
```

Dramatization. May not have happened.


