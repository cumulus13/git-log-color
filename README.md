# git-log-color

`git-log-color` is a wrapper around `git log` with:

- colorful graph output
- per-author sequential coloring, so each commit owner gets their own consistent color
- default emoji/icon rendering
- configurable styles and icons from a JSON config file
- structured output modes: `json`, `table`, `xml`, `yaml`, `toml`
- relative path support via `--path` or `-C`
- pass-through support for extra `git log` arguments

## Usage

```bash
git-log-color --path . --format text
git-log-color -C ../repo --all --oneline
git-log-color --format json -- --author=Hadi -- README.md
```

## Config

The tool looks for one of these files inside the repo:

- `.git-log-color.json`
- `git-log-color.json`
- `.config/git-log-color.json`

You can also pass a specific config path with `--config`.

Example:

```json
{
  "ui": {
    "color_enabled": true,
    "icons_enabled": true
  },
  "styles": {
    "commit": "bold #FFD700",
    "message": "white",
    "graph_main": "bold cyan",
    "author": "white bg:#550000 italic"
  },
  "icons": {
    "commit": "🔧 ",
    "version": "📦 ",
    "tag": "🏷️ "
  }
}
```

## Author Colors

[#author-colors](#author-colors)

In `text` format, each distinct commit author is colored using a fixed, sequential palette instead of a single flat `author` style: the first author seen gets the first color in the palette, the second author gets the second color, and so on. The same author always keeps the same color for the whole run, and the palette wraps around if there are more authors than colors.

This is on by default. To disable it (falling back to the single `author` style) or to customize the palette:

```json
{
  "ui": {
    "author_colors_enabled": true
  },
  "author_palette": [
    "bold #FF5555",
    "bold #55FF55",
    "bold #5599FF",
    "bold #FFD700"
  ]
}
```

Set `"author_colors_enabled": false` to disable per-author colors and use the static `author` style for everyone.

## Supported Style Keys

- `commit`
- `branch`
- `tag`
- `remote`
- `author`
- `date`
- `message`
- `graph_main`
- `graph_merge`
- `graph_branch`
- `graph_dim`
- `commit_dot`
- `merge_dot`
- `commit_graph`
- `username`
- `date_string`
- `date_number`
- `version_string`
- `version_number`
- `time_number`
- `microsecond`

Style syntax supports:

- named colors like `red`, `green`, `cyan`, `white`
- hex colors like `#FFAAFF`
- modifiers `bold`, `dim`, `italic`
- background colors with `bg:#AA0000`

## 📄 License

MIT

---

## 👤 Author
        
[Hadi Cahyadi](mailto:cumulus13@gmail.com)
    

[![Buy Me a Coffee](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/cumulus13)

[![Donate via Ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/cumulus13)
 
[Support me on Patreon](https://www.patreon.com/cumulus13)
