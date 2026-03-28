#### Maybe make it update on any change even non html


## CLI build
- [x] Checks if provided go file exists
- [x] Create a ./temp folder
- [x] Builds the go file to the temp folder
- [x] Writes the JS script to file in temp folder with WS address
- [ ] Start WS server in a separate goroutine and ping to check for health
- [ ] Start go binary
- [ ] Check for file changes and rebuild only the go file. Keep same script and keep WS server
- [ ] On stopping the tool, remove temp folder