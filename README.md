# TTY diary

just date selector with `selected date file preview` and editing it in `$EDITOR`


## TODO:
- [x] config file
- [ ] rewrite my `notes-alarm` bash script to go as `tty-diary-daemon`
- [ ] write README


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
