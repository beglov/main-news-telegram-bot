# Main News Telegram Bot

This Telegram bot sends the day's top news stories in a few messages.

## Usage

Copy env file and fill it with your credentials:

```bash
cp .env.sample .env
nano .env
```

Run bot to save chat IDs of interacting users:

```bash
go run cmd/bot/bot.go
```

Periodically send news to these chats by run:

```bash
go run cmd/send_news/send_news.go
```

## Deploy to VDS

To deploy a Golang program to a VDS (Virtual Dedicated Server), you can follow these steps:

1. Build your golang program - Before deploying your program, you need to build it.
   To do that, first, make sure you have golang installed on your local machine.
   Then, navigate to your project directory and run the following command:

```bash
go build -o bot cmd/bot/bot.go
go build -o send_news cmd/send_news/send_news.go
```

This will create an executable files in your project directory.

2. Transfer executable files to your VDS. You can do this using an FTP client or SCP. SCP example:

```bash
scp bot deploy@127.0.0.1:~/main-news-telegram-bot
scp send_news deploy@127.0.0.1:~/main-news-telegram-bot
```

### Deploy bot

3. Create a new systemd file:

```bash
sudo nano /etc/systemd/system/main-news-telegram-bot.service
```

Here is an example systemd service file for a Go program:

```
[Unit]
Description=Main News Telegram bot
After=network.target

[Service]
User=deploy
WorkingDirectory=/home/deploy/main-news-telegram-bot
ExecStart=/home/deploy/main-news-telegram-bot/bot

Environment=TELEGRAM_BOT_TOKEN=secret
Environment=RAPID_API_KEY=secret
Environment=POSTS_COUNT=50

Restart=on-failure
RestartSec=1

[Install]
WantedBy=multi-user.target

# See these pages for lots of options:
#   https://www.freedesktop.org/software/systemd/man/systemd.service.html
#   https://www.freedesktop.org/software/systemd/man/systemd.exec.html
```

_Make sure to replace the environment variables with your own values._

4. Save the file and reload systemd:

```bash
sudo systemctl daemon-reload
```

5. Start the service:

```bash
sudo systemctl start main-news-telegram-bot
```

6. Verify that the service is running:

```bash
sudo systemctl status main-news-telegram-bot
```

7. If everything is working correctly, enable the service to start at boot:

```bash
sudo systemctl enable main-news-telegram-bot
```

That's it! Your Golang program should now be running on your VDS and will automatically start up whenever your server reboots.

### Deploy send_news

8. Add .env file with your credentials. For an example, see the .env.sample file.

9. Add crone task to run send_news file every hour. To do that you can follow these steps:

   - Edit the crontab file using the command: `crontab -e`
   - In the crontab file, add a new line with the following syntax:
      ```
      0 * * * * cd /home/deploy/main-news-telegram-bot && ./send_news > cron_output.txt 2>&1
      ```
   - Save and exit the crontab file.

The cron task is now set up to execute the `/home/deploy/main-news-telegram-bot/send_news` file every hour.