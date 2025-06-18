# wiki setup

to actually use teawiki, you'll need a git repo that will contain all the pages
of your wiki, this documentation explains how you can setup this repo

## creating a git repo

i'm not gonna walk you through the basics of git or anything, if you want to use
teawiki, you are most likely already familiar with git (otherwise why even use a
wiki that's literally built for creating wikis with git?)

all the content of your wiki (pages, images etc.) will be contained in a single
repository, which you can simply create with `git init`. You do not need any
other setup, your entire wiki is literally stored in a single, standard git repo

in this repo you can create directories, files etc. Teawiki only treats markdown
files with `.md` extension as pages and that's pretty much it, anything else
will be served as a static file, and they will be accessible with their relative
paths

## creating pages

pages are simply markdown files, with optional metadata. Metadata let's you
change the page title, let's you add different tags and let's you create simple
infoboxes (all of which are covered later in this document)

the metadata is stored inside a "leaf block", which basically a text block
contained inside two `---` separators. This block is located at the start of the
file

so basically a page with has the following format:

```md
---
(metadata here)
---

(content here)
```

If you ever used [Hugo](https://gohugo.io/) or some other program that uses the
[goldmark-meta](https://github.com/yuin/goldmark-meta) extension, you might have
already seen this syntax (to be clear, teawiki does not use goldmark)

### metadata

page metadata is formatted with YAML, and it can contain few different keys

the most important one (and the one you'll mostly likely wanna use in every
page) is the `title` key. This key let's you specify a title for your page. Here
is an example:

```yaml
title: i386
```

**if no title is specified**, teawiki will check if the markdown content starts
with a h1 heading (heading with a single `#`). If so, the contents of the
heading (plain text contents, not rendered) will be used as the title and as a
result the heading will not be rendered as a part of the markdown content

**if no such heading exists**, teawiki will just use the file name as the title
after removing the `.md` extension. So if the markdown file has the name
`hello_world.md`, the page will have the title `hello_world`

another key is the `tags` key, which let's you specify a list of tags for a
page. Users can see these tags right under the title of the page, and click on
them to view pages with the same tags. They can also search for these tags,
which makes it easier for users to discover similar pages:

```yaml
tags:
  - cpu
  - intel
  - x86
```

another key you can use is the `image` key, which when specified will create an
[infobox](https://en.wikipedia.org/wiki/Infobox) with the specified image, for
example:

```yaml
image: images/i386.jpg
```

this looks like:

![](/assets/infobox1.png)

lastly you can use the `fields` key to create an infobox with custom key-value
paired fields:

```yaml
fields:
  - name: launched
    value: October 1985

  - name: discontinued
    value: September 28, 2007
```

you can also add links using to these fields:

```yaml
- name: predecessor
  value: Intel 80286
  link: i80286.md

- name: successor
  value: i486
  link: i486.md
```

this looks like:

![](/assets/infobox2.png)

### content

actual content of the pages is formatted with markdown. The markdown renderer
teawiki uses, [blackfriday](https://github.com/russross/blackfriday), supports
multiple extensions that are enabled:

- **no intra-word emphasis**: emphasis markers are ignored inside words
- **tables**: support for
  [markdown tables](https://www.markdownguide.org/extended-syntax/#tables)
- **fenced code**: support for
  [fenced code blocks](https://www.markdownguide.org/extended-syntax/#fenced-code-blocks)
- **definition list**: support for
  [definition lists](https://www.markdownguide.org/extended-syntax/#definition-lists)
- **footnotes**: support for
  [footnotes](https://www.markdownguide.org/extended-syntax/#footnotes)
- **auto URL linking**: automatically detects URLs and link them
- **heading IDs**: automatically creates IDs for headings, you can also specify
  custom IDs using `{#id}` syntax, this is useful because IDs are used link the
  headings
- **backslash line breaks**: converts trailing backslashes to line breaks

please note that these extensions are not part of the
[CommonMark standart](https://commonmark.org/), you should avoid using these if
you want your wiki pages to stay compatible with other markdown processors

## special pages

some pages have special uses in the wiki, the most important one is the
`README.md` page

if you create a page in the `README.md` file, this page will be used as the
default/index page for whatever directory it's in, and the directory will be
visible in the page listing section of the sidebar

for example, if you create a page in `/README.md`, that will be the index page
for your wiki and it will be displayed whenever users visit `/`

similarly, if you create a page in `/unix/README.md`, that will be the index
page for the `/unix` directory, and it will be displayed whenever users visit
`/unix`

the other special page is `/LICENSE.md`, if you create page at this location, a
link for it will be added to the navigation bar, and users will easily be able
to access it from anywhere

## other static content

as explained in the first section, anything that is not a page, will be served
as a static file without any rendering whatsoever

this means you can access & link all the images and other static content,
directly from the markdown content

for example let's say you have images under the `/images` directory, and you
want to use `my_image.png` in the page content, then all you need to do is:

```md
![my image](/images/my_image.png)
```

of course, you can also use relative paths

## using a local repo

if you are using a local git repo, you'll need make it accessible to the wiki
container

to do so you'll need to mount the repo to `TW_REPO_PATH`, for example if your
git repo is named `repo` and is located in the same directory with the compose
file, you can mount it like so:

```yaml
volumes:
  - ./repo:/tw/source:ro
```

for local repos, it's okay to use a read-only mount since teawiki does not need
to do any modification to the repo in order to load it's contents

when using a local repo, depending on how regularly you update the content, i
suggest you use a smaller `TW_RELOAD_INTERVAL` to keep the wiki content
synchronized with the actual git repo

## using a remote repo

if you are using a remote git repo, use the URL of the repo to configure
`TW_REPO_URL`, when the the wiki is started, the remote repo will be
cloned/pulled into `TW_REPO_PATH`, meaning you do not need to create a mount for
this path

however if you prefer to create a mount for some reason, make sure it's not a
read-only mount, as teawiki will need to modify the contents of the repo while
pulling content from the remote

when using a remote repo, i suggest you use a larger `TW_RELOAD_INTERVAL` to
prevent wasting bandwidth by sending too many requests to the remote server, and
ideally you should [configure a webhook](/docs/webhook.md) to keep the wiki
content synchronized with the git repo
