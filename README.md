![tbot](https://user-images.githubusercontent.com/25406248/125493306-565d7b6c-7582-4a74-a22b-639f33681b5b.gif)

- [Overview](#overview)
- [How It Works](#how-it-works)
- [Development](#development)
  - [Requirements](#requirements)
  - [Qualities](#qualities)
  - [Notes](#notes)
  - [Other Potential Enhancements](#other-potential-enhancements)

## Overview

An automated [Twitch](https://dev.twitch.tv/docs/irc) chat bot console application that runs from a command line interface (CLI).

Implemented with a Chuck Norris random fact function; listens for "!chucknorris" in currently joined chat room and sends a random fun fact about chuck norris in return.

## How It Works

1. [Download and Install Go](https://golang.org/doc/install).
2. Export values for `TWITCH_USERNAME` and `TWITCH_OAUTH_TOKEN` (can quickly generate an oauth token [here](https://twitchapps.com/tmi/)).
3. `make run`.

## Development

### Requirements

The bot application should be able to:

- [x] Console output all interactions - legibly formatted, with timestamps.
- [x] Connect to Twitch IRC over SSL.
- [x] Join a channel.
- [x] Read a channel.
- [x] Read a private message.
- [x] Write to a channel
- [ ] Reply to a private message.
- [x] Avoid premature disconnections by handling Twitch courier ping / pong requests - a request is sent every five minutes from twitch with PING, expects PONG, otherwise terminates.
- [x] Publicly reply to a user-issued string command within a channel (!YOUR_COMMAND_NAME).
- [x] Reply to the "!chucknorris" command by dynamically returning a random fact about Chuck Norris using the [Chuck Norris API](https://api.chucknorris.io).
- [x] Add gif to readme
- [x] Add Table of Contents to readme

### Qualities

- [ ] All interactions in this project are asynchronous.
- [ ] The application accounts for Twitch API rate limits. You can increase limits on your bot by becoming ["known" or "verified"](https://dev.twitch.tv/docs/irc/guide#known-and-verified-bots), but it looks like that [verification process is currently paused](https://discuss.dev.twitch.tv/t/an-update-for-the-delayed-bot-verification-request-process/32325).
- [x] The application does not exit prematurely.
- [x] The project is appropriately documented.
- [ ] The project has 100% test coverage. (WIP)

### Notes

There are [limits](https://dev.twitch.tv/docs/irc/guide#command--message-limits) of the number of IRC commands or messages you are allowed to send to the server. If you exceed these limits, you are locked out of chat for 30 minutes.

### Other Potential Enhancements

- Split CLI/GUI into two windows, one for reads and one for writes, for a better UX.
- Handle empty return "command not found" and general other unrecognized command wonkiness.
- Fix wonky printing of timestamp/name.
