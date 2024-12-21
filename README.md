# twitchets

[![Go Report
Card](https://goreportcard.com/badge/github.com/ahobsonsayers/twigots)](https://goreportcard.com/report/github.com/ahobsonsayers/twigots)
[![License - MIT](https://img.shields.io/badge/License-MIT-9C27B0)](LICENSE)

A tool to watch for ticket listings of desired events on [Twickets](https://www.twickets.live) that match certain filters (Discount, Number of tickets, Location), and send you a notification so you can snap them up!

I built this due to the official app notifications having a limit on the number of tracked events, and them not having the features i wanted. See features!

**Note**: this program does **not** buy tickets, reserve them automatically, or do anything unethical. All this does it notify you of new tickets!

Powered by (my similarly poorly named)
[twigots](https://github.com/ahobsonsayers/twigots), a package to help retreive these Twicket ticket listings and filter them watch

## Installation

```bash
go get -u github.com/ahobsonsayers/twigots
```

## Features

- Not limit on the number of Event you can track!
- Only get notified of Tickets with a certain discount
- Show more details in the notification - e.g. Event Date/Time, Number Tickets, Discount
- Faster to notify of new listings than the official Twickets app notificction
- No need to have the Twickets app or even an account.
- Your choice of notification service (Telegram, Ntfy, Gotify currently supported)

## Getting an API key

To use this tool you will need a Twickets API key. Twickets currently do not have a free API
HOWEVER it is possible to easily obtain an API key you can use.

To do this simply visit the [Twickets Live Feed](https://www.twickets.live/app/catalog/browse),
open you browser `Developer Tools` (by pressing `F12`), navigate to the `Network` tab, look for the
`GET` request to `https://www.twickets.live/services/catalogue` and copy the `api_key` query
parameter from the request.

This API key is not provided here due to liability concerns, but the key seems to be fixed/unchanging and
is very easy to get using the instructions above.

## Running

The best way to run abs-tract is to use Docker.

To run twitchets using Docker, use the following command:

```bash
docker run -d \
    --name twitchets \
    -v <path to config>:/twitchets/config.yaml
    --restart unless-stopped \
    arranhs/twitchets:latest
```

or if you are using docker compose:

```yaml
services:
  twichets:
    container_name: twitchets
    image: arranhs/twitchets:latest
    restart: unless-stopped
    volumes:
      - <path to config>:/twitchets
```

## Configuration

To start watching for ticket listings and getting notified, you must configure twickets.

Twickets will look for a `config.yaml` in your current working directory, and will fail to start if i can not be found.

The structure of this config.yaml can be seen in [`config.example.yaml`](./config.example.yaml) in this repo, or below

```yaml
apiKey: <MUST BE SET. See README.md for detail on how to obtain>

country: GB # Currently only GB is supported

# Configure your notification services.
# Delete the ones you do not need.
notification:
  ntfy:
    url: <your url>
    topic: <your topic>
    username: <your username>
    password: <your password>
  telegram:
    token: <your api token>
    chatId: <your chat id>
  gotify:
    url: <your url>
    token: <your api token>

# Global configuration that applies to all tickets.
# These can be overridden on a ticket by ticket basis (see below).
# The below are examples for all global configuration settings.
# Remove the settings you do not need.
global:
  # Regions ticket listing must be for.
  # Defaults to all regions if not set.
  # Full list of regions can be seen here:
  # https://github.com/ahobsonsayers/twigots/blob/main/location.go#L79-L90
  regions:
    - GBLO # Only look for tickets in London

  # How similar the event name in the ticket listing must be to the one you specified.
  # Defaults to 0.85 if not set, to allow for naming inconsistencies.
  # Between 0-1
  eventSimilarity: 0.85

  # How many tickets must be in the listing.
  # Defaults to any number of tickets if not set.
  numTickets: 2

  # Minimum discount (including fees) the tickets in the listing
  # must have against the original price.
  # Defaults to any discount of tickets if not set.
  discount: 15

  # Notification services to send found tickets listings to.
  # Defaults to all configured services above.
  notification:
    - ntfy

# Ticket configuration.
# Settings set here take priority over global settings.
# If a setting is not set/specified here, global setting will be used
# To unset a global setting and use the default, set to a negative value or empty list depends on type
# The below are examples. Tweak them as needed
tickets:
  # Look for Lion King tickets using the global configuration.
  # For global configuration settings not set, defaults will be used.
  - event: Lion King # Required

  # Look for Taylor Swift tickets, overriding the global configuration.
  # of the defaults for settings not specified.
  - event: Taylor Swift # Required

    # Override global configuration setting
    # Event name must be an exact match
    # For example This will NOT match Taylor Swift: The Eras Tour
    eventSimilarity: 1

    # Override global configuration setting
    # Reset to default - watch for tickets from any region
    regions: []

    # Override global configuration setting
    # Reset to default - watch for tickets with ANY discount
    numTickets: -1

    # Override global configuration setting
    # Reset to default - watch for tickets with ANY discount
    discount: -1

    # Override global configuration setting
    # Send notifications to telegram (instead of notify)
    notification:
      - telegram

  # Look for Lion King tickets generally using the global configuration.
  # Notification will be sent to all notification services configures
  - event: Oasis # Required
    notification: []
```

## Why the name twitchets?

Because I feel like sometimes you need to have twitch like reactions to snap up tickets on Twickets before someone else gets them - which this tool helps you do. Therefore the mangling together of **twitch** and **Twickets** seemed fun and appropriate.

[![Hits](https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https%3A%2F%2Fgithub.com%2Fahobsonsayers%2Ftwitchets&count_bg=%2379C83D&title_bg=%23555555&icon=&icon_color=%23E7E7E7&title=visitors+day+%2F+total&edge_flat=false)](https://hits.seeyoufarm.com)
