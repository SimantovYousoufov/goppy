Goppy
=============

Goppy is a command-line clipboard watcher. It'll keep track of your clipboard history in a terminal window and manage storing it behind the scenes for you. Goppy will hydrate clipboard history from storage.

## Usage:

* Options:

```
  -e	Use the encrypted file storage format.
  -f string
    	Path to history file (default "$PWD/goppy_history.dat")
  -h	Get usage information
  -n int
    	Number of items to keep in history. (default 50)
  -s	Use the NullStore to prevent clipboard history from being written to disk.
```

* Start goppy with default options: 

    `$ goppy`
    
* Start goppy and specify the location of the history file and file encryption

    `$ goppy -f /some/path/to/file -e` 

## Available Storage Types

1. File Store - goppy will use a simple JSON file to store your clipboard history
    - NOTE: Your clipboard history is stored in plaintext.
2. Encrypted File Store - goppy will use an encrypted format to store your clipboard history to a file with a password you choose
    - NOTE: goppy is a project built for personal use, no actual security is guaranteed and it has not been audited
    - NOTE: If you're using terminal as a view (default) and you have any form of logging enabled, your clipboard history will be stored there too.
    
## TODO

1. Implement way to clear history
1. Implement some form of verbosity and debug logging
1. Ignore new copy events if the string already exists in history beyond the most recent event
