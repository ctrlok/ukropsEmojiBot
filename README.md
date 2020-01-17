# Emoji Bot
Emoji bot is a simple code with lambda, which will do actions 
based on reactions in slack. Created for [ukrops](https://ukrops.club)
community. See `Contribution` if you want to use it for your community. 

# Features
- when user react to the message with :to_best: reaction - send this message 
to the special channel.

# Deploy
```bash
terraform apply
```

# Terraform structure
Right now there is two terraform workspaces: `production` and `development`.
You need to select needed workspace first.  
Both of this spaces uses different URLs, so we need a two different slack apps. 
Slack apps configuration also different: production bot are subscribed to events 
in every channel, while development bot are subscribed on events in channel it invited to. 

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