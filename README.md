# Mongo

[![Mongo Test](https://github.com/akshaybabloo/mongo/workflows/Mongo%20Test/badge.svg)](https://github.com/akshaybabloo/mongo/actions)

A simple wrapper for Go's Mongo Driver

## Instillation

Use Go modules

```
go get github.com/akshaybabloo/mongo
```

## Usage

> Unlike the MongoDB driver, this library depends on `id` and NOT `_id`. That means you will have to create an index for `id` field.

See `example_test.go`
