# customization

this documentation covers how you can customize the look of your wiki by using a
custom color them and using custom assets

## custom themes

themes are simply short CSS files that define the colors that will be used in
the frontend application

to create a custom color theme, you can just copy the
[dark theme file](/static/themes/dark.css) and change the colors and save it as
your own theme file

for this example we'll assume the file is saved as `my_theme.css` and it's
located in the same directory as the compose file

in order to use this theme, you need to make sure it's accessible by teawiki, so
you'll need to mount the theme file to `static/themes` directory:

```yaml

---
volumes:
  - ./my_theme.css:/tw/static/themes/:ro
```

after mounting your theme file, you can now use your custom theme by using the
`TW_THEME` configuration option, for this example the configuration would look
like this:

```yaml

---
environment:
  TW_THEME: "my_theme"
```

## custom assets

assets are image files that is used by the frontend application

there are currently two different assets that you can modify:

- **logo (`TW_LOGO`)**: this is a 130x130 pixel image file that is displayed as
  the logo of your wiki in the sidebar of the frontend application

- **icon (`TW_ICON`)**: this is the website icon that is displayed by the
  browser, usually in the corner of the tab menu

you can create your own images to replace these assets, and mount them into
`static/assets` to make sure they are accessible by teawiki, for example:

```yaml

---
volumes:
  - ./my_logo.png:/tw/static/assets
  - ./my_icon.svg:/tw/static/assets
```

then you'll need to specify the name of icons using the configuration options of
the assets you would like to change, for this example the configuration would
look like this:

```yaml

---
environment:
  TW_LOGO: "my_logo.png"
  TW_ICON: "my_icon.svg"
```
