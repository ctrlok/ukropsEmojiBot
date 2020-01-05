# Emoji Bot
Emoji bot is a simple code with lambda, which will do actions 
based on reactions in slack. Created for [ukrops](https://ukrops.club)
community.

# Features
TODO

# Deploy
```bash
terraform apply
```

### Requirements
- `go ≥ 1.13`
- `terraform ≥ 0.12`

# Contribution

Usual fork and ~~exec~~ PR works.
You can create an issue. 
If you want to use this bot in your community - please, create
an issue, and I will remove all ukrops-specific code. 

### Code structure
In the main folder is a terraform files for deploy code and infrastructure
In the `slackConnector` folder is a code for a lambda