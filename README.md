# teawiki | simple git based wiki

![](https://img.shields.io/github/actions/workflow/status/ngn13/teawiki/test.yml?label=tests)
![](https://img.shields.io/github/actions/workflow/status/ngn13/teawiki/docker.yml?label=build)

a simple HTTP web application that let's you create wikis using git and markdown

![](assets/showcase.png)

I created this program for my own wiki, so I specifically designed it around my
personal needs, however I documented everything as well as I can and tried to
make everything as much configurable as possible, so you can also use it if it
fulfills your needs as well

## features

- free software (free as in freedom),
  [learn more](https://www.gnu.org/philosophy/free-sw.en.html)
- easy installation and configuration with docker compose
- supports both local and remote git repos
- simple and minimal web interface inspired by
  [MediaWiki's MonoBook skin](https://www.mediawiki.org/wiki/Skin:MonoBook)
- configurable light and dark theme, you can also use custom themes
- YAML and markdown based article/page format (with code syntax highlighting)
- very simple, static [infobox](https://en.wikipedia.org/wiki/Infobox) support
  (also inspired by MediaWiki)
- page tags and easy search functionality
- webhook support for syncing with the remote repos instantly
- [sitemap](https://www.sitemaps.org/) generation

## missing features (that other wikis usually have)

- no user account system
- no web-based editor/manager
- not easily extensible (no plugins, extensions etc.)

personally, I don't really care about any of these features, but you might care,
in this case I suggest you look for an another wiki software (there are
[plenty](https://awesome-selfhosted.net/tags/wikis.html) of them)

## installation

to install and run teawiki, I suggest you use docker compose, this is the
_intended_ deployment option

an [example compose file](compose.example.yml) can be found in the repo, copy
this and read the documentation to configure everything properly

I also suggest you use a reverse proxy server instead of directly exposing the
docker container to the internet, so you can configure stuff like SSL, CORS etc.

## documentation

- [configuration](/docs/config.md)
- [wiki setup](/docs/wiki.md)
- [webhook setup](/docs/webhook.md)
- [customization](/docs/custom.md)

## development

for development, clone the repository and switch to `dev` branch, all of your
pull requests should also target this branch

to build the application you'll need GNU make, `go` and the SASS compiler
`sassc`, after obtaining these, you can build the app by running:

```bash
make
```

to build the application in release mode, run:

```bash
make RELEASE=1
```

to check format the code properly, you'll also need `gofmt` (part of the go
toolkit, if you installed `go` you probably already have this) and prettier.
After obtaining these, you can format the code by running:

```bash
make format
```

to check for any syntax/formatting errors, you'll also need GNU grep and sed:

```bash
make check
```

I also wrote few test scripts to make my life easier, to run these you will also
need PCRE2, `curl`, `openssl` and `htmlq`. After obtaining these tools you can
run the tests by running:

```bash
make test
```

### adding translations

all the translations are in the `locale` directory, fork the repo and create the
translation file using the
[ISO 639](https://en.wikipedia.org/wiki/List_of_ISO_639_language_codes) code of
the language you want

then copy the contents of `en.yaml` locale, translate all the stuff and create a
pull request

there are roughly 150 words to translate (maybe even less) so this should only
take few minutes

to check for any missing locales, run the `checks/locale.sh` script. This script
requires the `yq` tool, so make sure to install it before running the script:

```bash
./checks/locale.sh
```

### reporting issues

if you encounter a problem, please create an issue with your docker compose file
(**after removing any sensitive information**) and explain the problem you are
encountering in detail

### other contributions

if you are planning to add a new feature, please first create an issue to
discuss it - you don't need to do this for bug fixes

before creating a pull request make sure you run `make format` and `make check`
to format the code and check the code for any formatting errors, as described in
the [development section](##development). I also suggest you run the tests to
make sure everything is working properly
