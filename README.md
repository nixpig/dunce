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

| Implemented | Model Method  | Params                               | Returns               |
| ----------- | ------------- | ------------------------------------ | --------------------- |
| [ ]         | `Create`      | `article ArticleData`                | `( *Article, error )` |
| [ ]         | `GetAll`      |                                      | `[]Article`           |
| [ ]         | `GetByType`   | `typeId int`                         | `[]Article`           |
| [ ]         | `GetByAuthor` | `authorId int`                       | `[]Article`           |
| [ ]         | `GetById`     | `articleId int`                      | `Article`             |
| [ ]         | `GetBySlug`   | `articleSlug string`                 | `Article`             |
| [ ]         | `GetByTag`    | `tagId int`                          | `[]Article`           |
| [ ]         | `UpdateById`  | `articleId int, article ArticleData` | `Article`             |
| [ ]         | `DeleteById`  | `articleId int`                      | `bool`                |

### `UserModel`

| Implemented | Model Method    | Params                               | Returns          |
| ----------- | --------------- | ------------------------------------ | ---------------- |
| [x]         | `Create`        | `newUser *UserData, password string` | `*User, error`   |
| [x]         | `UpdateById`    | `id int, user *UserData`             | `*User, error`   |
| [x]         | `GetAll`        |                                      | `*[]User, error` |
| [x]         | `GetById`       | `userId int`                         | `*User, error`   |
| [x]         | `GetByUsername` | `username string`                    | `*User, error`   |
| [x]         | `GetByEmail`    | `email string`                       | `*User, error`   |
| [ ]         | `GetByRole`     | `role RoleName`                      | `*[]User, error` |
| [ ]         | `DeleteById`    | `userId int`                         | `error`          |

### Tags

| Implemented | Model Method | Params                   | Returns |
| ----------- | ------------ | ------------------------ | ------- |
| [ ]         | `Create`     | `tag TagData`            | `Tag`   |
| [ ]         | `GetAll`     |                          | `[]Tag` |
| [ ]         | `GetById`    | `tagId int`              | `Tag`   |
| [ ]         | `GetBySlug`  | `tagSlug string`         | `Tag`   |
| [ ]         | `UpdateById` | `tagId int, tag TagData` | `Tag`   |
| [ ]         | `DeleteById` | `tagId int`              | `bool`  |

### Types

| Implemented | Model Method | Params                      | Returns  |
| ----------- | ------------ | --------------------------- | -------- |
| [ ]         | `Create`     | `type TypeData`             | `Type`   |
| [ ]         | `GetAll`     |                             | `[]Type` |
| [ ]         | `GetById`    | `typeId int`                | `Type`   |
| [ ]         | `GetBySlug`  | `typeSlug string`           | `Type`   |
| [ ]         | `UpdateById` | `typeId int, type TypeData` | `Type`   |
| [ ]         | `DeleteById` | `typeId int`                | `bool`   |

### Site

| Implemented | Model Method | Params      | Returns |
| ----------- | ------------ | ----------- | ------- |
| [ ]         | `Update`     | `site Site` | `Site`  |
