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

### Tags

| DB  | Service | REST | Proto | Model Method | Params           | Returns         |
| --- | ------- | ---- | ----- | ------------ | ---------------- | --------------- |
| [ ] | [ ]     | [ ]  | [ ]   | `Create`     | `tag TagData`    | `*Tag, error`   |
| [ ] | [ ]     | [ ]  | [ ]   | `GetAll`     |                  | `*[]Tag, error` |
| [ ] | [ ]     | [ ]  | [ ]   | `GetById`    | `tagId int`      | `*Tag, error`   |
| [ ] | [ ]     | [ ]  | [ ]   | `GetBySlug`  | `tagSlug string` | `*Tag,error`    |
| [ ] | [ ]     | [ ]  | [ ]   | `Update`     | `tag TagData`    | `*Tag, error`   |
| [ ] | [ ]     | [ ]  | [ ]   | `DeleteById` | `tagId int`      | `error`         |

### Articles

| DB  | Service | REST | Proto | Model Method      | Params                 | Returns             |
| --- | ------- | ---- | ----- | ----------------- | ---------------------- | ------------------- |
| [ ] | [ ]     | [ ]  | [ ]   | `Create`          | `article ArticleData`  | `*Article, error`   |
| [ ] | [ ]     | [ ]  | [ ]   | `GetAll`          |                        | `*[]Article, error` |
| [ ] | [ ]     | [ ]  | [ ]   | `GetByTypeName`   | `typeName string`      | `*[]Article, error` |
| [ ] | [ ]     | [ ]  | [ ]   | `GetByAuthorName` | `authorName string`    | `*[]Article, error` |
| [ ] | [ ]     | [ ]  | [ ]   | `GetById`         | `articleId int`        | `*Article, error`   |
| [ ] | [ ]     | [ ]  | [ ]   | `GetBySlug`       | `articleSlug string`   | `*Article, error`   |
| [ ] | [ ]     | [ ]  | [ ]   | `GetByTagName`    | `tagName string`       | `*[]Article, error` |
| [ ] | [ ]     | [ ]  | [ ]   | `Update`          | `article *ArticleData` | `*Article, error`   |
| [ ] | [ ]     | [ ]  | [ ]   | `DeleteById`      | `articleId int`        | `error`             |

### `UserModel`

| DB  | Service | REST | Proto | Model Method    | Params                               | Returns          |
| --- | ------- | ---- | ----- | --------------- | ------------------------------------ | ---------------- |
| [ ] | [ ]     | [ ]  | [ ]   | `Create`        | `newUser *UserData, password string` | `*User, error`   |
| [ ] | [ ]     | [ ]  | [ ]   | `Update`        | `user *UserData`                     | `*User, error`   |
| [ ] | [ ]     | [ ]  | [ ]   | `GetAll`        |                                      | `*[]User, error` |
| [ ] | [ ]     | [ ]  | [ ]   | `GetById`       | `userId int`                         | `*User, error`   |
| [ ] | [ ]     | [ ]  | [ ]   | `GetByUsername` | `username string`                    | `*User, error`   |
| [ ] | [ ]     | [ ]  | [ ]   | `GetByEmail`    | `email string`                       | `*User, error`   |
| [ ] | [ ]     | [ ]  | [ ]   | `GetByRole`     | `role RoleName`                      | `*[]User, error` |
| [ ] | [ ]     | [ ]  | [ ]   | `DeleteById`    | `userId int`                         | `error`          |

### Types

| DB  | Service | REST | Proto | Model Method | Params            | Returns          |
| --- | ------- | ---- | ----- | ------------ | ----------------- | ---------------- |
| [ ] | [ ]     | [ ]  | [ ]   | `Create`     | `type TypeData`   | `*Type, error`   |
| [ ] | [ ]     | [ ]  | [ ]   | `GetAll`     |                   | `*[]Type, error` |
| [ ] | [ ]     | [ ]  | [ ]   | `GetById`    | `typeId int`      | `*Type, error`   |
| [ ] | [ ]     | [ ]  | [ ]   | `GetBySlug`  | `typeSlug string` | `*Type, error`   |
| [ ] | [ ]     | [ ]  | [ ]   | `Update`     | `type TypeData`   | `*Type, error`   |
| [ ] | [ ]     | [ ]  | [ ]   | `DeleteById` | `typeId int`      | `error`          |

### Site

| DB  | Service | REST | Proto | Model Method | Params      | Returns |
| --- | ------- | ---- | ----- | ------------ | ----------- | ------- |
| [ ] | [ ]     | [ ]  | [ ]   | `Get`        |             | `*Site` |
| [ ] | [ ]     | [ ]  | [ ]   | `Update`     | `site Site` | `*Site` |
