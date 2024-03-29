[![Workflow Status](https://github.com/nixpig/dunce/actions/workflows/validate.yml/badge.svg?branch=main)](https://github.com/nixpig/dunce/actions/workflows/validate.yml?query=branch%3Amain)
[![Coverage Status](https://coveralls.io/repos/github/nixpig/dunce/badge.svg?branch=main)](https://coveralls.io/github/nixpig/dunce?branch=main)

# 🧠 dunce

A HATEOAS-first headless (brainless) content publishing engine.

## TODO

> MOVE OFF OF FIBER

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

## Template functions

| Function               | Description |
| ---------------------- | ----------- |
| `GetUsers`             |             |
| `GetUserByUsername`    |             |
| `GetUserById`          |             |
| `GetLoggedInUser`      |             |
| `GetTags`              |             |
| `GetTypes`             |             |
| `GetArticles`          |             |
| `GetArticlesByAuthor`  |             |
| `GetArticlesByType`    |             |
| `GetArticlesByTagName` |             |
| `GetArticleBySlug`     |             |
| `GetArticleById`       |             |
| `SiteName`             |             |
| `SiteDescription`      |             |
| `SiteUrl`              |             |
| `SiteOwner`            |             |
| -                      | -           |
| `Login`                |             |
| `Logout`               |             |

## Models API

### Articles

| Implemented | Model Method  | Params                               | Returns             |
| ----------- | ------------- | ------------------------------------ | ------------------- |
| [ ]         | `Create`      | `article ArticleData`                | `*Article, error`   |
| [ ]         | `GetAll`      |                                      | `*[]Article, error` |
| [ ]         | `GetByType`   | `typeId int`                         | `*[]Article, error` |
| [ ]         | `GetByAuthor` | `authorId int`                       | `*[]Article, error` |
| [ ]         | `GetById`     | `articleId int`                      | `*Article, error`   |
| [ ]         | `GetBySlug`   | `articleSlug string`                 | `*Article, error`   |
| [ ]         | `GetByTag`    | `tagId int`                          | `*[]Article, error` |
| [ ]         | `UpdateById`  | `articleId int, article ArticleData` | `*Article, error`   |
| [ ]         | `DeleteById`  | `articleId int`                      | `error`             |

### `UserModel`

| Implemented | Model Method    | Params                               | Returns          |
| ----------- | --------------- | ------------------------------------ | ---------------- |
| [x]         | `Create`        | `newUser *UserData, password string` | `*User, error`   |
| [x]         | `UpdateById`    | `id int, user *UserData`             | `*User, error`   |
| [x]         | `GetAll`        |                                      | `*[]User, error` |
| [x]         | `GetById`       | `userId int`                         | `*User, error`   |
| [x]         | `GetByUsername` | `username string`                    | `*User, error`   |
| [x]         | `GetByEmail`    | `email string`                       | `*User, error`   |
| [x]         | `GetByRole`     | `role RoleName`                      | `*[]User, error` |
| [x]         | `DeleteById`    | `userId int`                         | `error`          |

### Tags

| Implemented | Model Method | Params                   | Returns         |
| ----------- | ------------ | ------------------------ | --------------- |
| [x]         | `Create`     | `tag TagData`            | `*Tag, error`   |
| [x]         | `GetAll`     |                          | `*[]Tag, error` |
| [x]         | `GetById`    | `tagId int`              | `*Tag, error`   |
| [x]         | `GetBySlug`  | `tagSlug string`         | `*Tag,error`    |
| [x]         | `UpdateById` | `tagId int, tag TagData` | `*Tag, error`   |
| [x]         | `DeleteById` | `tagId int`              | `error`         |

### Types

| Implemented | Model Method | Params                      | Returns          |
| ----------- | ------------ | --------------------------- | ---------------- |
| [ ]         | `Create`     | `type TypeData`             | `*Type, error`   |
| [ ]         | `GetAll`     |                             | `*[]Type, error` |
| [ ]         | `GetById`    | `typeId int`                | `*Type, error`   |
| [ ]         | `GetBySlug`  | `typeSlug string`           | `*Type, error`   |
| [ ]         | `UpdateById` | `typeId int, type TypeData` | `*Type, error`   |
| [ ]         | `DeleteById` | `typeId int`                | `error`          |

### Site

| Implemented | Model Method | Params      | Returns |
| ----------- | ------------ | ----------- | ------- |
| [ ]         | `Update`     | `site Site` | `Site`  |
