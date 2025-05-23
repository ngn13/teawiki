# webhook setup

when using a remote repo, webhooks are the best way to keep your wiki synced
with the actual git repo

currently teawiki supports webhooks for the following git hosting platforms:

- [Github](https://github.com)
- [Gitea](https://about.gitea.com)
- [Forgejo](https://forgejo.org)

## creating a webhook

to create a webhook, find the webhook settings for your repo, different git
hosting platforms use different interfaces, so the location of these settings
may differ

### gitea & forgejo

visit the settings page for your wiki's repo, then click on "Webhooks", you
should see a "Add Webhook" button, click this and select gitea/forgejo from the
dropdown

- for the target URL
  - enter `[https or http]://[address for your wiki]/_/webhook/gitea` for gitea
  - enter `[https or http]://[address for your wiki]/_/webhook/forgejo` for
    forgejo
- for the HTTP method, select `POST`
- for the POST content type, select `application/json`
- for the secret, create a random and long string using something like a
  password generator, make sure to save this for later
- for the trigger, select "Push Events"

### github

visit the settings page for your wiki's repo, then click on "Webhooks", you
should see a "Add webhook" button, click on this

- for the payload URL, enter
  `[https or http]://[address for your wiki]/_/webhook/github`
- for the content type, select `application/json`
- for the secret, create a random and long string using something like a
  password generator, make sure to save this for later
- disable SSL verification if you are using a self-signed certificate
- for the "Which events would you like to trigger this webhook?" option, select
  "Just the push event"

## using the webhook

provide the webhook secret you saved to teawiki using the `TW_WEBHOOK_SECRET`
option:

```yaml
environment:
  TW_WEBHOOK_SECRET: "topsecret"
```

and that's it, now whenever you do a push to the remote repo, your wiki will be
synced with the remote repo almost instantly
