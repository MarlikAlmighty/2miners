# Monitoring Bot for 2miners.com

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white) &nbsp;
![React](https://img.shields.io/badge/react-%2320232a.svg?style=for-the-badge&logo=react&logoColor=%2361DAFB) &nbsp;
![Swagger](https://img.shields.io/badge/-Swagger-%23Clojure?style=for-the-badge&logo=swagger&logoColor=white) &nbsp;
![Gmail](https://img.shields.io/badge/Gmail-D14836?style=for-the-badge&logo=gmail&logoColor=white) &nbsp;

---

## What is it?

This is my vision of how this bot should work https://t.me/Pool2MinersBot

My application monitors the API of the 2miners.com website.

I implemented a workers pool in it. This allows using one goroutine for a request.
I used boltDB as a data storage. This will save us from unnecessary services.
I use React for visual data management in the application.

The application supports registration, cookie-based authentication. Email should be used as a login.
As soon as you add your wallet, the application starts a goroutine that monitors the API.
In order not to receive 429, there is a parameter in the configuration, the number of seconds between requests.

The application sends an email as a notification.
*This part should be rewritten, since I use a personal email account.*

<img src="https://github.com/MarlikAlmighty/2miners_bot/blob/main/images/2024-07-30-04-02-29-372.jpg" width="200" alt="Authentication"> &nbsp;
<img src="https://github.com/MarlikAlmighty/2miners_bot/blob/main/images/2024-07-30-04-03-09-550.jpg" width="200" alt="Dashboard"> &nbsp;
<img src="https://github.com/MarlikAlmighty/2miners_bot/blob/main/images/2024-07-30-06-04-24-431.jpg" width="200" alt="Add purse"> &nbsp;

<br />

<img src="https://github.com/MarlikAlmighty/2miners_bot/blob/main/images/2024-07-30-04-04-30-217.jpg" width="200" alt="Dashboard after"> &nbsp;
<img src="https://github.com/MarlikAlmighty/2miners_bot/blob/main/images/2024-07-30-04-04-39-420.jpg" width="200" alt="Stats"> &nbsp;
<img src="https://github.com/MarlikAlmighty/2miners_bot/blob/main/images/2024-07-30-04-04-48-607.jpg" width="200" alt="Settings"> &nbsp;

## History

Initially, I generated the project via swagger. But then I refused, nevertheless I got beautiful models and methods))
I use my personal email account in this project, it is not professional, this solution was for development,
so this part should be rewritten.

## How to run

```shell
$ export MAX_ADDR=3
```

Limit for users, the number of addresses for each user that he can monitor.

```shell
$ export REQUEST_OVER_TIME=3
```

Time in seconds, during which time the application sleeps. This is important in order not to receive a 429 response from the server.

```shell
$ export COOKIE_BLOCK_KEY="something_similar_to_hash"
$ export COOKIE_HASH_KEY="another_similar_to_hash"
```

This is for encrypting cookies. Cookies contain the user ID, so it is important to encrypt them.

```shell
$ export COIN_MARKET_CAP_API_KEY="token_coin_market_cap"
```

With the help of this token we get the value of the coin, this is necessary for conversion into dollars.

```shell
$ export DOMAIN=marlikalmighty.uk
$ export SMTP_PASSWORD="*** *** *** ***" // secret
$ export SMTP_USER="user@gmail.com" // real user
$ export SMTP_PORT="587"
$ export SMTP_HOST="smtp.gmail.com"
```

This is for sending mail and generating a link. When registering, we send a link to the user, after clicking on it, he becomes authorized.

I do not recommend doing this. This was needed for development. This should be rewritten.

---

This story lacks an admin panel to manage users, I almost started, but unfortunately I have no time. There are no tests either))
