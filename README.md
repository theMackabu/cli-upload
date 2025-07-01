upload files to [themackabu/uploader](https://github.com/theMackabu/uploader)

```bash
# create config.json from template
# and edit it with appropriate values
go build -o upload
./upload <filename> [options]
```

```bash
# upload public file
./upload myfile.txt

# upload private file
./upload myfile.txt -p
./upload myfile.txt --private
```
