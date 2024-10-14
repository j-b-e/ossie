# ossie
A powerful Tool to manage Openstack Contexts

*Inspired by [kubie](https://github.com/sbstp/kubie)*

# Usage

- `[required]` parameter
- `<optional>` parameter

* `ossie rc` display a selectable menu of rc-files
* `ossie rc [rc]` selects the rc and spawns a new shell in this context
* `ossie rc -` switch back to the previous context
* `ossie rc [rc] -r <region>` spawns the shell with this region if available
* `ossie regions` display a selectable menu of regions
* `ossie edit [rc]` opens rc-file with default editor 
* `ossie api-version [service] <version>` sets api-version of the service or shows selectable menu
* `ossie export` prints the current configured rc-file to stdout
* `ossie info <rc>` shows current or selected rc in nice overview
* `ossie create <rc>` menu driven creation of rc from scratch or based on rc
