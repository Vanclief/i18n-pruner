# i18n-pruner

**i18n-pruner** its intended to be used for JSON files that are used for i18n.
It receives as input two json files: `source` and `destination`.

It will:

- Format both files and order the keys alphabetically
- Find any duplicate values in the `destination` file
- It will make sure that the `destination` file is not missing any key from
  the `source` file, if it does it will add them with an empty string as value

Optionally:

- Translate missing values using ChatGPT

## Usage

Simple run:

`i18n-pruner -s en.json -d es.json`

Translate missing values to spanish:

`i18n-pruner -s en.json -d es.json -t Spanish`

Arguments

- `-s | --source` is the path to the source file
- `-d | --destination` is the path to the destination file
- `-r | --read-only` is a flag that will prevent the program from saving changes and will only list the duplicate values
- `-t | --translate` the language to translate the missing values to

## Instalation

Run:

`curl -L https://github.com/Vanclief/i18n-pruner/raw/master/install.sh | bash`

You can also manually download the binaries from `/bin/`
