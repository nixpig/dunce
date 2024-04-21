[![Workflow Status](https://github.com/nixpig/dunce/actions/workflows/validate.yml/badge.svg?branch=main)](https://github.com/nixpig/dunce/actions/workflows/validate.yml?query=branch%3Amain)
[![Coverage Status](https://coveralls.io/repos/github/nixpig/dunce/badge.svg?branch=main)](https://coveralls.io/github/nixpig/dunce?branch=main)

# 🧠 dunce

Super-simple platform used to build my personal site.

## ⚠️ WORK IN PROGRESS

Don't try to use this. Not only does it probably not run yet, it's specific to _my_ use case.

## Models API

### Tags

| DB  | Service | Controller | Server | Model Method | Params           | Returns         |
| --- | ------- | ---------- | ------ | ------------ | ---------------- | --------------- |
| [x] | [x]     | [ ]        | [ ]    | `Create`     | `tag TagData`    | `*Tag, error`   |
| [x] | [x]     | [ ]        | [ ]    | `GetAll`     |                  | `*[]Tag, error` |
| [x] | [x]     | [ ]        | [ ]    | `DeleteById` | `tagId int`      | `error`         |
| [x] | [x]     | [ ]        | [ ]    | `GetBySlug`  | `tagSlug string` | `*Tag,error`    |
| [x] | [x]     | [ ]        | [ ]    | `Update`     | `tag TagData`    | `*Tag, error`   |
| [ ] | [x]     | [ ]        | [ ]    | `GetById`    | `tagId int`      | `*Tag, error`   |

### Articles

| DB  | Service | REST | Proto | Model Method   | Params                 | Returns             |
| --- | ------- | ---- | ----- | -------------- | ---------------------- | ------------------- |
| [ ] | [ ]     | [ ]  | [ ]   | `Create`       | `article ArticleData`  | `*Article, error`   |
| [ ] | [ ]     | [ ]  | [ ]   | `GetAll`       |                        | `*[]Article, error` |
| [ ] | [ ]     | [ ]  | [ ]   | `GetById`      | `articleId int`        | `*Article, error`   |
| [ ] | [ ]     | [ ]  | [ ]   | `GetBySlug`    | `articleSlug string`   | `*Article, error`   |
| [ ] | [ ]     | [ ]  | [ ]   | `GetByTagName` | `tagName string`       | `*[]Article, error` |
| [ ] | [ ]     | [ ]  | [ ]   | `Update`       | `article *ArticleData` | `*Article, error`   |
| [ ] | [ ]     | [ ]  | [ ]   | `DeleteById`   | `articleId int`        | `error`             |

### Site

| DB  | Service | REST | Proto | Model Method | Params      | Returns |
| --- | ------- | ---- | ----- | ------------ | ----------- | ------- |
| [ ] | [ ]     | [ ]  | [ ]   | `Get`        |             | `*Site` |
| [ ] | [ ]     | [ ]  | [ ]   | `Update`     | `site Site` | `*Site` |
