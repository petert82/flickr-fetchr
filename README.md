# Flickr Fetchr

## About

Flickr Fetchr is a command line utility for downloading a minimal set of information about a Flickr photostream to a local JSON file.

For all public photos in a given user's photostream it will save:

- The title
- The description
- Large and small thumbnail URLs
- Original-size image thumbnail

## Usage

You'll need a working install of [Go][1] to be able to build the `flickr-fetchr` binary.

Once your Go environment is up and running, you should be able to clone this repository and run `go get` from inside your local copy of it to build the app.

Once built, simply run the `flickr-fetchr` binary like so:

```sh
$ flickr-fetchr --api-key=[Your Flickr API key] --user-id=[Users flickr ID] --output-file=[Path to save to]
```

[1]: https://golang.org/