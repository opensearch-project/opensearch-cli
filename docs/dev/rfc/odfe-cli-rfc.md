# RFC: ODFE Command Line Interface



## 1. Overview

The ODFE Command Line Interface `(odfe-cli)` is an open source tool that enables 
users to manage their odfe clusters using through CLI or command-line-shell.
Additionally, command line interface will allow users to configure and manage 
odfe plugins configurations like anomaly detections, alerting, sql, and 
access Elasticsearch features from command line with prominence on 
automation. 

This cli enables users to start running commands that implement 
Kibana-equivalent functionality from the command prompt in their favorite 
terminal programs with a minimal configuration. As a programmer or an admin, 
users constantly want to perform ad-hoc operations that Kibana may not support. 
odfe-cli is better suited when the user wants to quickly combine a couple of commands 
to perform a task. These commands are powerful and concise. Since shell 
commands can be combined into script files, one can build sequence of commands to automate 
operations that they do often. This also provides direct access to the odfe cluster. 
The user can engage with different plugins using odfe-cli, and develop shell 
scripts to manage their cluster.

## 2. Tenets

1. Unified CLI for Open Distro for Elasticsearch.
2. Available for installation on all platforms (Mac , Windows and Linux).
3. odfe-cli will be distributed as Docker image just like any other image to enable possibilities for automation CI / CD.
4. Provide profiles to handle multiple clusters with different user preferences such as Configure security credentials, 
default input/output file format at once and apply across multiple plugins / cluster.
5. Flexibility to work with multiple odfe clusters.
6. Scaling your workflow allows you to automate and replicate and make life .
7. It also supports Elasticsearch core rest apis from [here](https://www.elastic.co/guide/en/elasticsearch/reference/current/rest-apis.html) 
as commands eg: cat nodes, index create, cluster health, etc...
8. Provide extra functions which are hard to implement on Kibana, like batch action, long running tasks, etc.



## 3. Design and Architecture

![](../images/design.png)

### 3.1 Commands

Command is the central point of the application. Every interaction that cli 
supports will be represented as command. A command can have sub-commands, 
arguments and flags . Usually command represent an action, arguments are the 
things and flags are ways to modify the behavior of the commands.

odfe-cli will be the root command and it contains three type of sub-commands

* es core commands as index get, index create, cluster health, etc ...
* plugins as commands like ad, monitor, sql, etc...
* connections, which helps set-up profiles that will be used to interact with 
clusters.

### 3.2 Input Handler

Command supports input from both command line and as file. The file could be 
either json or yaml. Input handler is responsible to understand this input 
and convert it into an entity which can be passed to controller. This can be 
extended by supporting more file types like csv or text in the future.

### 3.3 Controller

Controller implements actual business logic for a command. Every plugin should 
have at least one controller. and command should depend only on its controller. 
A controller can depend on another controller and gateway. 
Also, Plugin controller will implement method for every command.
For ex: AnomalyDetection Plugin will have a controller 
**AnomalyDetectorController** which implements methods like Create, Start, 
Stop, Delete, Update and so on. Also, AD Controller needs profile information 
like connection url, user name and password, hence, will depend on **ProfileController** 
which implements methods like Get, Create, Delete, etc...

### 3.4 Gateway

Gateway acts like an interface between application and cluster. 
Every plugin will have a gateway which implements the REST api provided by 
corresponding plugin. This Gateway is independent of commands and controller. 
Gateway requires implementing REST APIs provided by a plugin and return 
response as it is. Also, Gateway doesn't know about the type of cluster. It only uses 
Rest client with security details to communicate via REST protocol.

### 3.5 Model (Entities)

Model or Entities represents the request and response structure required to 
pass data to and from clusters and from commands to controller.

### 3.6 Output Handler

Every plugin has different styles to display output, for example, sql-cli 
might want to display output in different way than monitor plugin. Hence, 
individual commands  can customize the output based on output handler 
before displaying it on the output peripheral.

### 3.8 Rest Client

This client provides a HTTP client interface with automatic retries and 
exponential backoff on top of net/http package. Since elasticsearch exposes 
REST API, every gateway uses this client to communicate to cluster.



## 4. Install & Set-up

ODFE Command line interface can be installed in following ways

1. Install `odfe-cli` Using Docker Image
2. Install `odfe-cli` on Linux
3. Install `odfe-cli` on macOS
4. Install `odfe-cli` on Windows

## Synopsis

```
odfe-cli [options] <command> <sub-command> [parameters]
```

## options

*`--profile`  [profile name]*
Overrides commands default profile with given profile. Default profile name 
is `default` or, `$ODFE_PROFILE`

`--debug`
displays more messages about CLI operations, especially for debugging 
purposes. If you like to report a github issue, share output generated 
with `--debug` when submitting it. 

`--verbose`
displays useful information about operations performed by commands.

`--input-format` (string)
Overrides default formatting style for file input/output with given format. 
The available file formats are:

* json (default)
* yaml

``--version`` (string)
Display the version of this tool.

## Usage

```
# Start detector job from Anomaly Detection
$>odfe-cli ad start-detectors --id MuIjkjLEFWhuelQmmN-

# Sync monitors from remote to local from Alerting
$>odfe-cli monitor sync

# Returns information about one or more indices
$>odfe-cli es index get my-index-000001
```

You can get help on the command line to see the supported plugins

```
$>odfe-cli help
```

the operations for plugins

```
$>odfe-cli ad help
```

the parameters for a plugins operation

```
$>odfe-cli es cat help
```



## Request for Comments:

We would like comments and feedback on the proposal for building ODFE CLI 
tool [here](https://github.com/opendistro-for-elasticsearch/es-cli/issues/2). 
Some specific questions we’re seeking feedback include


* Do you prefer odfe-cli as root command name? 
* What function you want for the CLI tool? What’s your use case?
* Would you like to see every feature available in kibana to be available in 
CLI or would you prefer only features which are difficult and impossible to 
do in UI to be available.
* Do you prefer new release whenever a new odfe version releases, or an 
independent release based on cli features?

