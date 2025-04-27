# configuration
all the configuration is environment variable based, which makes the
configuration easy for deployment with docker compose

all the environment variable that are used for configuration use the `TW_`
prefix, any other environment variable will not effect the configuration

here is a list of all the available options:

- [`TW_LISTEN` (required)](#tw_listen-required)
- [`TW_URL`](#tw_url)
- [`TW_NAME` (required)](#tw_name-required)
- [`TW_DESC` (required)](#tw_desc-required)
- [`TW_KEYWORDS`](#tw_keywords)
- [`TW_SOURCE_URL`](#tw_source_url)
- [`TW_COMMIT_URL`](#tw_commit_url)
- [`TW_REPO_PATH` (required)](#tw_repo_path-required)
- [`TW_REPO_URL`](#tw_repo_url)
- [`TW_RELOAD_INTERVAL`](#tw_reload_interval)
- [`TW_WEBHOOK_SECRET`](#tw_webhook_secret)
- [`TW_CHROMA`](#tw_chroma)
- [`TW_THEME` (required)](#tw_theme-required)
- [`TW_LANG`](#tw_lang)
- [`TW_TIME`](#tw_time)
- [`TW_LOGO`](#tw_logo)
- [`TW_ICON`](#tw_icon)

## `TW_LISTEN` (required)
host the web server will listen on

default is `127.0.0.1:8080`

## `TW_URL`
full URL for your wiki

for example if you are hosting your wiki on `wiki.example.com` using HTTPS,
under `/tw`, this URL would be `https://wiki.example.com/tw`

this option is not required but i highly suggest you set it, as it's required
for [sitemap](https://www.sitemaps.org/) generation

## `TW_NAME` (required)
name of the wiki, displayed in the HTML title

default is `my_wiki`

## `TW_DESC` (required)
a short description for your wiki, displayed in HTML meta elements

default is `my personal wiki`

## `TW_KEYWORDS`
keywords for your wiki, will be included the in HTML meta elements, split
multiple words by comma

for example if you wiki is about old audio devices, you may want to set this to
`audio,old,history,wiki`

default is `wiki`

## `TW_SOURCE_URL`
URL for the web interface that is used to display git repo's source code

for example if your wiki lives in `https://github.com/example/wiki`, and uses
the `main` branch, this URL should be set to
`https://github.com/example/wiki/tree/main`

this URL changes across different git hosting software and it also depends on
how the software is deployed, this is why it's configurable

## `TW_COMMIT_URL`
URL for the web interface that is used to display git repo's commits, for
example if your wiki lives in `https://github.com/example/wiki`, this should be
set to `https://github.com/example/wiki/commit`

## `TW_REPO_PATH` (required)
path to the local git repo, no need to modify if you are using a remote repo
with `TW_REPO_URL`

default is `./source`

## `TW_REPO_URL`
URL for the remote git repo

if specified, this remote repo will be pulled/cloned to `TW_REPO_PATH`

## `TW_RELOAD_INTERVAL`
when teawiki starts up, it loads all the pages in memory to speed up rendering
when a user requests a page

however the contents of the repo may change, so this interval essentially
creates a reload timer that reloads the pages from the git repo

this reload timer also triggers a pull from `TW_REPO_URL` if it's configured, so
it can also used to keep the local repo synced with the remote

by default this interval is set to thirty minutes, and you can this by
specifying a number with a `h` (hour), `m` (minute) or `s` (second) suffix, so
for example to set the interval to 3 hours, set this variable to `3h`

if you configured `TW_REPO_URL` and if you are using this option to keep the
local repo synced with the remote, you should also checkout `TW_WEBHOOK_SECRET`
to configure a webhook for better and almost instant syncing

## `TW_WEBHOOK_SECRET`
secret value for the remote git hosting platform's webhook, should be the same
secret that you use while configuring the webhook, this is used by the server
to verify the authenticity of webhook requests

to learn more about setting up a webhook,
[see this documentation](/docs/webhook.md)

## `TW_CHROMA`
style for the [chroma syntax highlighter](https://github.com/alecthomas/chroma),
check out [Chroma Style Gallery](https://xyproto.github.io/splash/docs/) to view
the available styles

by default, if you are using the dark theme, it's set to `rrt`, however if you
are using the light theme, it's set to `emacs`

see `TW_THEME` to see the theme configuration

## `TW_THEME` (required)
name of the color theme, `dark` and `light` are the only available themes,
however [you can also use a custom theme](/docs/custom.md)

default is `dark`

## `TW_LANG`
language for the web application, currently available languages are:

- English (`en`, default)
- Turkish (`tr`)

adding your own language is quite simple, see the translation section of
the [README](/README.md) for more info

## `TW_TIME`
time format that will be used to display stuff like commit times

default format is `02/01/06 15:04:05 GMT-07`, so that's DD-MM-YY HH:MM:SS and
timezone, to use a different format,
[see how time formatting is done in go](https://go.dev/src/time/format.go)

## `TW_LOGO`
name of the logo file, which is displayed in the sidebar of the web interface,
by default it's set to `logo.png`

## `TW_ICON`
name of the icon file, which is displayed as the browser tab icon, by default
it's set to `logo.png`

to change the icon and the logo, see the
[customization documentation](/docs/custom.md)
