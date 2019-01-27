# fsm
Ridiculously hacky tool to help me keep in touch with my friends.

## Problem

Interacting with my friends makes me happy, but I have an inertial tendency to lose touch. Also, I live away from 99% of them at the moment, so it's hard.

## Solution

I hacked together a Google Sheet with my list of friends, the frequency I want to talk to them, and a chat link. Then I wrote this go script that reads the sheet, calculates if the frequencies are met or not, choses a random friend out of the matched ones and opens the chat link. There should be a cron on the system that runs every now and then, but one should also be able to run it manually.

## Chat links look like:

- https://www.facebook.com/messages/t/%name_or_id%
- whatsapp://send?text=Time%20to%20talk!&phone=%number%
- https://mail.google.com/mail/?extsrc=mailto&url=%email%

## Hacky things that make it hacky

- I just call "open %chat_link%". Beware of what's in that Sheet cell, open doesn't work for everyone, have to be logged in to FB, WhatsApp Web and Gmail.
- I call `Get` and `Update` on the Sheets API. For `Get` I read the whole sheet (btw Sheet name hardcoded). For `Update` I calculate the cell to write extremely hackily.
- In order to make Google Sheets API work I created a service account that has access to that particular sheet and downloaded a `credentials.json`, so this is not a `git clone` away from working.
- The auth permission I created is ridiculously permissive. I think it's like "this app can delete all your spreadsheets".

## If you want to use it

The implementation is terrible, but it works quite well for me at the moment. If you want to try it out feel free to contact me via email or create an issue.
