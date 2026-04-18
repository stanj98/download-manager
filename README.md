# Concurrent File Downloader in Go

A simple concurrent file downloader written in Go that splits a file into multiple sections, downloads them in parallel, and merges them into a final output file.

---

## Features

- Downloads files using HTTP range requests
- Splits file into multiple sections
- Concurrent downloads using goroutines
- Merges all parts into a single file
- Simple and minimal dependencies

---

## How It Works

1. Sends a `HEAD` request to get the file size
2. Splits the file into equal sections
3. Downloads each section concurrently using goroutines
4. Saves each section as a temporary file
5. Merges all sections into the final output file
