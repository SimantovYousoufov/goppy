Goppy
=============

Goppy is a command-line clipboard watcher. It'll keep track of your clipboard history in a terminal window and manage storing it behind the scenes for you. Goppy will hydrate clipboard history from storage.

## Usage:

* Start goppy with default options: 

    `$ goppy`
    
* Start goppy and specify the location of the history file

    `$ goppy -f /some/path/to/file`

## Available Storage Types

1. File Store - goppy will use a simple JSON file to store your clipboard history
    - NOTE: Your clipboard history is stored in plaintext
