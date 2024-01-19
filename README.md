# 🧠 bloggor

A HATEOAS-first headless (brainless) content publishing engine.

## ⚠️ WORK IN PROGRESS

Don't try to use this yet. It probably doesn't even run.

### todo!

**Basically everything!**

## Motivation

Frankly, pretty fed up with REST APIs and had a crazy idea for a _headless_ HTML API. Want to see where it goes...

## Structure

```

├── cmd/
│   └── app/
│       └── main.go  # entrypoint
├── db/
│   └── migrations/  # database migrations
├── deploy/  # deployment config
├── internal/  # internals of app
│   ├── app/  # the fiber server and associated stuffs
│   └── pkg
│       ├── config/  # configuration
│       └── models/  # data models and db connection
├── pkg
│   └── api/  # the functions api to expose to templates
├── web
│   ├── static/  # static web assets
│   └── templates/  # customisable templates

```
