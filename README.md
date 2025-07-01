upload files to [themackabu/uploader](https://github.com/theMackabu/uploader)

```bash
# edit config/config.json
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
