<h1 align="center">gruppetto</h1>

<div align="center">

concerto turn server 🌈

[![Docker tests](https://github.com/concerto-app/gruppetto/actions/workflows/test-docker.yml/badge.svg)](https://github.com/concerto-app/gruppetto/actions/workflows/test-docker.yml)
[![Docs](https://github.com/concerto-app/gruppetto/actions/workflows/docs.yml/badge.svg)](https://github.com/concerto-app/gruppetto/actions/workflows/docs.yml)

</div>

---

This `README` provides info about the development process.

For more info about the package itself
see `gruppetto/README.md`
or [docs](https://concerto-app.github.io/gruppetto).

## Quickstart (on Ubuntu)

```shell
$ apt update && apt install curl git mercurial make binutils bison gcc build-essential
$ curl -sSL https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer -o gvm.sh
$ bash gvm.sh && exec bash
$ gvm install go1.4 -B && gvm use go1.4
$ GOROOT_BOOTSTRAP=$GOROOT gvm install go1.18
$ git clone https://github.com/concerto-app/gruppetto
$ cd gruppetto/gruppetto
$ go run cmd/main.go
```

## Testing

Just go to `gruppetto` directory and run:

```shell
go test ./...
```

## Building docs

We are using [`mkdocs`](https://www.mkdocs.org)
with [`material`](https://squidfunk.github.io/mkdocs-material)
for building the docs. It lets you write the docs in Markdown format and
creates a nice webpage for them.

Docs should be placed in `gruppetto/docs/docs`. They
are pretty straightforward to write.

To build the docs,
`cd` into `gruppetto/docs` and run:

```sh
mkdocs build
```

It will generate `site` directory with the webpage source.

## Continuous Integration

When you push changes to remote, different GitHub Actions run to ensure project
consistency. There are defined workflows for:

- deploying docs to GitHub Pages
- testing inside Docker container
- drafting release notes
- publishing Docker images

For more info see the files in `.github/workflows` directory and `Actions` tab
on GitHub.

Generally if you see a red mark next to your commit on GitHub or a failing
status on badges in `README`
it means the commit broke something (or workflows themselves are broken).

## Releases

Every time you merge a pull request into main, a draft release is automatically
updated, adding the pull request to changelog. Changes can be categorized by
using labels. You can configure that in `.github/release-drafter.yml` file.

Every time you publish a release the Docker image is built and uploaded to GitHub registry with tag taken from release
tag

## Docker

You can build a Docker image of the package (e.g. for deployment). The build
process is defined in `Dockerfile` and it's optimized to keep the size small.

To build the image, run from project root:

```sh
 docker build -t gruppetto .
```

To also run the container in one go, run:

```sh
docker build -t gruppetto . && docker run --rm -it gruppetto
```
