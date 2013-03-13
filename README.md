Go (golang) bindings for libstemmer
===================================

This simple library provides Go (golang) bindings for the snowball libstemmer library including the popular porter and porter2 algorithms.

Requirements
------------

You'll need the development package of libstemmer, usually this is simply a matter of:

    sudo apt-get install libstemmer-dev

... or you might need to [install it from source](http://snowball.tartarus.org/).

Installation
------------

First, ensure you have your GOPATH env variable set to the root of your Go project:

    export GOPATH=`pwd`
    export PATH=$PATH:$GOPATH/bin

Then this cute statement should do the trick:

    go get github.com/rjohnsondev/golibstemmer

Usage
-----

Basic usage:

    package main

    import "github.com/rjohnsondev/golibstemmer"
    import "fmt"
    import "os"

    func main() {
        s, err := stemmer.NewStemmer("english")
        defer s.Close()
        if err != nil {
            fmt.Println("Error creating stemmer: "+err.Error())
            os.Exit(1)
        }
        word := s.StemWord("happy")
        fmt.Println(word)
    }

To get a list of supported stemming algorithms:

    list := stemmer.GetSupportedLanguages()

Testing
-------

You can execute the basic included tests with:

    go test

If you have issues, double check you've installed the libstemmer development library.

If you still have issues, let me know!
