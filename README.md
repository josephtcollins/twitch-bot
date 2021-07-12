## Overview

An automated [Twitch](https://dev.twitch.tv/docs/irc) chat bot console application that runs from a command line interface (CLI).

## How It Works

TBW

## Development

### Requirements

The bot application should be able to:

- [ ] Console output all interactions - legibly formatted, with timestamps. Strip cruft.
- [x] Connect to Twitch IRC over SSL.
- [x] Join a channel.
- [x] Read a channel.
- [x] Read a private message.
- [x] Write to a channel
- [ ] Reply to a private message.
- [x] Avoid premature disconnections by handling Twitch courier ping / pong requests - a request is sent every five minutes from twitch with PING, expects PONG, otherwise terminates.
- [ ] Publicly reply to a user-issued string command within a channel (!YOUR_COMMAND_NAME).
- [ ] Reply to the "!chucknorris" command by dynamically returning a random fact about Chuck Norris using the [Chuck Norris API](https://api.chucknorris.io).

### Qualities

- [ ] All interactions in this project are asynchronous.
- [ ] The application accounts for Twitch API rate limits. You can increase limits on your bot by becoming ["known" or "verified"](https://dev.twitch.tv/docs/irc/guide#known-and-verified-bots), but it looks like that [verification process is currently paused](https://discuss.dev.twitch.tv/t/an-update-for-the-delayed-bot-verification-request-process/32325).
- [ ] The application does not exit prematurely.
- [ ] The project is appropriately documented.
- [ ] The project has 100% test coverage.

### Notes

There are [limits](https://dev.twitch.tv/docs/irc/guide#command--message-limits) of the number of IRC commands or messages you are allowed to send to the server. If you exceed these limits, you are locked out of chat for 30 minutes.