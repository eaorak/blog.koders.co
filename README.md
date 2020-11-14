# [Koders Blog](https://blog.koders.co)

# How to Configure

- You can update cron time from `.github/workflows/update-blog.yml` (`on > schedule > cron`)
- You can update daily post count from `update/main.go` (`dailyPostCount`)
- You can update maximum post count from `update/main.go` (`maxPostCount`)

# Update Blog Posts

Install [golang](hhttps://golang.org/doc/install) and run the following command. This command is fetches last days top posts from dev.to API and updates the repo. Don't forget to commit changes after command finish.

```
go run update/main.go
```

# Run blog

To update dependencies: `yarn`

To build: `yarn build`

To run locally: `yarn dev`
