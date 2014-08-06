#panic-attack has reached v0.1.1. 
This means it might possibily work, but you would be a fool to run it automatically.

when run panic-attack will transform code like this:

    func critical_path() {
        data, _ := GetData()
        //TODO handle error
        return data
    }

to

    func critical_path() {
        data, err := GetData()
    if err != nil {
        panic(err)
    }
        //TODO handle error
        return data
    }

I have found myself typing that bit of code too often, so I made this tool.

For vim users:

    $ go install github.com/placeybordeaux/panic-attack
    $ echo 'command Panic :%!panic-attack % | goimports' >> ~/.vimrc


# TODO

* Get real flags
  * have a linting flag
* Identify unchecked errors in local packages
* Handle multiple files or even packages
* Find edge cases where an error is not detected (parser.ParseFile comes to mind)
* Handle := vs = properly
* Verify that tool doesn't need multiple passes
* Move this list into github issues
* Tests maybe
* Handle indentions properly
