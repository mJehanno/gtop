# Todo

## V1

- [ ] Change some filenames so that linux targeted file only get build on linux (`[filename]_linux.go`), this might need some refacto
- [ ] Fix cpu usage => progress bar is always empty, freq might be the metric we need to set it as it shows current freq and not total freq
- [ ] Process manager should sort process based on cpu/mem usage
- [ ] Process manager should be able to send signals to process (ex: SIGTERM)
- [ ] Process manager should be able to filter on user input 
- [ ] Adds network usage tab (might want to group all network rellated stuff ?)
- [ ] Adds goreleaser (go binary, apt, aur, snap ?, homebrew, chocolatey) - only linux based for V1 rest for V2
- [ ] Adds a license and a readme with some gif made with [VHS](https://github.com/charmbracelet/vhs)


## Somewhere between V1 and V2

- [ ] Improve error handling
- [ ] Adds "responsiveness"


## V2

- [ ] Adds MacOS support (`[filename]_darwin.go`)
- [ ] Adds Windobe suppport (`[filename]_windows.go`)
- [ ] Separate code to create a lib (`gtop-core` ?) used by the tui app and unit tested (Especially unmarshalling)
- [ ] Might need to pass all 64b type to their flexible equivalent (int64 => int) to ensure compilation for 32bit computer

## V3

- [ ] [wails](https://wails.io/) or [fyne](https://fyne.io/) Gui based on lib


> note.txt contains some helpfull knowledge to get the job done
