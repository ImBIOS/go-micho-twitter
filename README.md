# Twitter API with Go 👋

Twitter-like API with Go, using [Echo](https://github.com/labstack/echo) and [Go-Micro](https://github.com/go-micro/go-micro) framework.

## Features 🪶

- [x] Authentication
  - [x] Sign up
  - [x] Sign in
- [ ] Follow System
  - [x] Follow
  - [ ] Unfollow
- [x] Public Tweet
  - [x] Send Tweet
  - [x] Feed
- [ ] Private Message
  - [ ] Send Message
  - [ ] View Private Message
- [ ] User Profile
  - [ ] View Profile
  - [ ] Edit Profile
- [ ] Search
- [ ] Security
  - [ ] JWT
    - [x] Access Token
    - [ ] Refresh Token
  - [ ] Rate Limiting
  - [ ] CORS
- [ ] DevOps
  - [ ] Docker
  - [ ] Kubernetes
  - [x] Continuous Integration
  - [ ] Containerized Development

## Tools ⚒️

- Dev Watcher => [Air](https://github.com/cosmtrek/air)
- Linter => [GolangCI-Lint](https://github.com/golangci/golangci-lint)
- Code Formatter => [Golines](https://github.com/segmentio/golines)

## VSCode Extension 👾

- [Run on Save](https://marketplace.visualstudio.com/items?itemName=emeraldwalk.RunOnSave)

## Project Structure 🌊

```bash
TODO:
```

- **.vscode** => TODO:
- **configs** => TODO:
- **controllers** => TODO:
- **helpers** => TODO:
- **models** => TODO:
- **routes** => TODO:

## Get Started 👟

1. Create a MongoDB database, and name it `twitter`
2. Copy the example environment variable (`.env.example`) to `.env`

```bash
cp .env.example .env
```

3. Configure the `.env` file

4. Make sure Go (min v1.18) and [Air](https://github.com/cosmtrek/air) installed.

5. Run the Air 🛫

```bash
air
```
