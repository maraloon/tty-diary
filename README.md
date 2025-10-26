# TTY diary

- run it, select date with [datepicker](https://github.com/maraloon/datepicker), write something in your `$EDITOR`, exit
- if you add time (ex: `01:15 text`) at start of string, get notification


![demo gif](readme/demo.gif) 

## Install

### manualy:

```bash
git clone git@github.com:maraloon/tty-diary.git
cd tty-diary
go install ./...
```


## Usage

- `diary` for calling TUI

- add `* * * * * diary-notify` cron job to get notifications

## Config file

### Paths

- `$HOME/ttydiary.yaml`
- `$HOME/.config/ttydiary.yaml`

### Format

```sh
file-color: "#b1a286"
file-format: "md"
diary-dir: "~/Documents/diary"
monday: true
```

`file-color`: if file `{diary-dir}/yyyy/mm/dd.{file-format}` is existed and not emtpy (= you have some notes on this date), highlight this date on calendar with this color

`monday: true` for monday as week start or `sunday: true` for sunday

---

It depends on [datepicker](https://github.com/maraloon/datepicker) bubble, which you can use in your go apps too
Also see [pickdate](https://github.com/maraloon/pickdate) as cli-wrapper of datepicker
