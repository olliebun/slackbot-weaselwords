# Weasel Bot

A Slack bot for warning users about weasel words in their Slack messaging.

For example, if you post a message like this:

	Sorry, I just wanted to see if it would be okay if I could just work from home for a little while next week one day in the morning? Let me know.

The bot will DM you like this:

![DM screenshot](http://i.imgur.com/47Fx6Hh.png)

In no time you'll be sending these bad boys:

	WFH today. Have slack and phone.

Inspired by [the chrome extension](http://www.slate.com/blogs/xx_factor/2015/12/29/new_chrome_app_helps_women_stop_saying_just_and_sorry_in_emails.html) to help stop women from saying "just" and "sorry" in emails.

## Building

[Install gb](https://getgb.io) then run:

	gb build all

The binary will be in `bin`.

## Configuration

The bot is configured at runtime with environment variables. It looks for these variables:

* `SERVER_ADDR`
* `WORDS_FILE`
* `USERS_FILE`
* `SLACK_TOKEN`

`SERVER_ADDR` is just passed to Go's `http.ListenAndServe` function. It should be of the form $HOST:$PORT. Examples:

* `:8080`
* `192.168.0.1:9090`

WORDS_FILE and USERS_FILE must both point at existing files (relative to the working directory). These files must each have a word match or user per line.

Each line in WORDS_FILE is just a string to search for.

Each line in USERS_FILE is a username.

SLACK_TOKEN must be a valid API token for a bot user in your team.

## Testing

If you want to help test, [contact me on twitter](https://twitter.com) and I'll add you to the Slack team with a test instance of the bot, and probably give you a token to use.