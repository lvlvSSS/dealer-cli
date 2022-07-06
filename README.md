# dealer-cli

## Introduction

a simple cli automatic tool.

depends on the cli tools - github.com/urfave/cli/v2

## Installatioin

    go get github.com/lvlvSSS/dealer-cli.git    

## Usage

Flag:

- mode `global flag`
    - define the current mode of dealer-cli. values : debug, info, warn, info
- log `global flag`
    - define the logging's location of dealer-cli
- load-yaml `global flag`
    - define the yaml configuration file.

commands :

- file ```parent command, it should have subcommands.```
    - extract ```child command of file. it could extract some message from logging.```

- schedule  ```parent command, it should have subcommands.```
    - http ```child command of schedule. it could send http messages to server.```