package cli

const defaultCfg = `# Configuration for martian, tool for Stationeers localization.
# See https://github.com/st-l10n/martian for reference.

# List of available languages.
languages:
  - code: RU
    name: Russian
    font: russian

  - code: EN
    name: English
    font: english

  - code: FR
    name: French
    font: extended

  - code: DE
    name: German
    font: extended

  - code: IT
    name: Italian
    font: extended

  - code: KN
    name: Japanese
    font: cjk
    locale: ja

  - code: KO
    name: Korean
    font: hangul

  - code: PL
    name: Polish
    font: extended

  - code: PT
    name: Portuges
    prefix: portuguese

  - code: CN
    name: Simplified Chinese
    font: cjk
    locale: zh-CN

  - code: TW
    name: Traditional Chinese
    prefix: traditional-chinese
    font: cjk
    locale: zh-TW

  # Finnish
  - code: FI
    name: Suomi

  - code: SK
    name: Slovak
    font: russian
`
