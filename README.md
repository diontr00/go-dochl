[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit)](https://github.com/pre-commit/pre-commit)
![ci workflow](https://github.com/diontr00/go-dochl/actions/workflows/ci.yml/badge.svg)

# Intallation

```
go install github.com/diontr00/go-dochl@latest
```

# Usage

Default Keywords: [todo, fixme, bug, hack ]
Any comment start with either of them (case insensitive) will be highlighted and extracted into and print toward stdout.

![show case](https://i.imgur.com/uFcpj6K.png)
![show case](https://i.imgur.com/5hR47vq.png)

Using with pre-commit:

```
repos:
  - repo: https://github.com/diontr00/go-dochl
    rev: latest
    hooks:
      - id: go-dochl
        args: [--keys="TODO,FIX,HACK"] # optional
```
