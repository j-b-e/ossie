# ossie

A powerful Tool to manage Openstack Environments

## Usage

* `ossie rc` Displays a selectable menu of Openstack environments
* `ossie rc <rc>` Spawns a new shell with the selected  Openstack environment
* `ossie info <rc>` Shows information about the Openstack environment


## Features

* Spawn Shell with selected Openstack environment from autodetected clouds.yaml or RC Files
* Configurable Prompt
* Protects OS_ env against accidental changes
* Includes Quality-of-Life enhancements
    * `osenv` - Prints `OS_*` Environment
    * `o`and `os`aliases for `openstack`
* Shell Support
    * Bash

## Configuration

Have a look in `ossie.toml.example` for available Settings. Run ossie with `-c <FILE>` flag to specify the configuration file or copy it to `.config/openstack/ossie.toml` for auto-loading.

## Installation

Download the latest **ossie** Binary from [Github Releases](https://github.com/j-b-e/ossie/releases) and place the file in your PATH.



---
_**Ossie** is inspired by [Kubie](https://github.com/sbstp/kubie)_
