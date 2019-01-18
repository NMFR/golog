[![Build Status](https://travis-ci.org/mlimaloureiro/golog.svg)](https://travis-ci.org/mlimaloureiro/golog)
[![Coverage](http://gocover.io/_badge/github.com/mlimaloureiro/golog?0)](http://gocover.io/github.com/mlimaloureiro/golog)
[![License](http://img.shields.io/:license-apache-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0.html)

# golog
`golog` is an easy and lightweight CLI tool to time track your tasks. The goal is to enable to track concurrent from small to big tasks.

![](http://i.imgur.com/o2F0JbW.gif?1)

# Overview
I work in a very fast paced company, and I'm always receiving requests, plus *a lot* of small requests and I've struggled to find a tool that fit my needs. We do use other tools to track the time spent on a task, but sometimes it gets so overwhelming that it's just not worth to create a bunch of small tasks and track them **but you do want to track them**. If you have your terminal always opened like me, `golog` is perfect for this environments, you can log multiple tasks at the same time without going to your browser/proj management tool to improve productiveness.

# Installation
Make sure you have a working Go environment (go 1.1+ is *required*). [See the install instructions](http://golang.org/doc/install.html).

To get `golog`, run:
```
$ go get github.com/mlimaloureiro/golog
```

To install it in your path so that `golog`can be easily used:

```
$ cd $GOPATH/src/github.com/mlimaloureiro/golog
$ GOBIN="/usr/local/bin" go install
```

#### Enabling autocomplete

Copy `autocomplete/bash_autocomplete` into `/etc/bash_completion.d/golog`.
Don't forget to source the file to make it active in the current shell.

```
   sudo cp autocomplete/bash_autocomplete /etc/bash_completion.d/golog
   source /etc/bash_completion.d/golog
```

Alternatively, you can just source `autocomplete/bash_autocomplete` in your bash configuration with `$PROG` set to golog.

```
PROG=golog source "$GOPATH/src/github.com/mlimaloureiro/golog/autocomplete/bash_autocomplete"
```

If using `zsh` use `zsh_autocomplete`

```
PROG=golog source "$GOPATH/src/github.com/mlimaloureiro/golog/autocomplete/zsh_autocomplete"
```

## Getting Started

The **start** command will start tracking time for a given taskname. **Note that a taskname cannot have white spaces**, they serve as identifiers.

```
$ golog start {taskname}
```

To stop tracking use the **stop** command, if you want to **resume** a stopped task just golog start {taskname} again.

```
$ golog stop {taskname}
```

With the **list** command you can see all your tasks and see which of them are active.

```
$ golog list
0h:1m:10s    create-readme (running)
0h:0m:44s    do-some-task
```

If you only want to check the status of one task, use the **status** command.

```
$ golog status create-readme
0h:3m:55s    create-readme (running)
```

The lifetime of the info I usually need is very short (actually is just a workday), in the next day it's unlikely that i'll need previous info. This is one case where **clear** command is handy.

```
$ golog clear
All tasks were deleted.
```

You can use the **export** command to export all the tasks to a file. 

The available formats are: 

- csv (.csv)
- ical (.ics)

``` sh
golog export [csv | ical] [file_path]
```

# Contribution Guidelines
@TODO
If you have any questions feel free to link @mlimaloureiro to the issue in question and we can review it together.
