Goppy
=============

Goppy is a command-line clipboard watcher. It'll keep track of your clipboard history in a terminal window and manage storing it behind the scenes for you. Goppy will hydrate clipboard history from storage.

## Warnings

- While Goppy keeps all data 100% local, it has access to your clipboard, this means Goppy can read any sensitive information like passwords, API keys, etc that you copy. Please make sure you read the source and are comfortable with this before using Goppy.
- Prefer to not store your history to disk, but if you choose to do so, use the Encrypted Store `$ goppy -e`

## Usage:

* Options:

```
  -e	Use the encrypted file storage format.
  -f string
    	Path to history file. (default "/usr/local/etc/goppy/goppy_history.dat")
  -h	Get usage information.
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

1. Implement some form of verbosity and debug logging
1. Ignore new copy events if the string already exists in history beyond the most recent event
1. Create a nice GUI for goppy as an alternative to using the terminal view
1. Implement dependency manager
1. Implement an easy installer (i.e. brew)
