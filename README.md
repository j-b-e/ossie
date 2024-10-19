# ossie

A powerful Tool to manage Openstack Environments

## Usage

* `ossie rc` Displays a selectable menu of Openstack environments
* `ossie rc <rc>` Spawns a new shell with the selected  Openstack environment
* `ossie rc -` Switch back to the previous environment
* `ossie export` Exports active or selected menu-driven environment in specifed format (rc or clouds.yaml)
* `ossie export <rc>` Export pre-selected environment
* `ossie info <rc>` Shows current or selected environment in an overview


## Features

* Spawn Shell with selected Openstack environment from autodetected clouds.yaml or RC Files
* Configurable Prompt
* Protects OS_ env against accidental reset
* Includes Quality-of-Life utilities
    * `osenv` - Prints `OS_*` Environment
* Shell Support
    * Bash


## Installation

Download the latest **ossie** Binary from [Github Releases](https://github.com/j-b-e/ossie/releases) and place the file in your PATH.


---
_**Ossie** is inspired by [Kubie](https://github.com/sbstp/kubie)_
