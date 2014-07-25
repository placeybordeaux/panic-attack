panic-attack eventually supposed to be a tool that will run through your go code and transform lines such as

    data, _ := GetData()
    //TODO handle error

to
    data, err := GetData()
    if err != nil {
        panic(err)
    }

I have found myself typing that bit of code too often.

right now it can find functions that assign anything to _ . Next is to parse discovery the function signature of that given function.
