# Guidelines for aesthetic coding :fire:
When we say aesthetic coding, we mean code efficiently and maximise readability. This is done through basic compliance to Go principles and following software engineering rules of thumb. This is a non-exhaustive list of general things to keep in mind.

1) **DO USE** a text editor with good Golang support
    * Examples would be [Visual Studio Code](https://code.visualstudio.com/), [Atom](https://atom.io/) or even vim if that tickles your fancy
    * Make sure that the supported Golang plugins are installed
    * Please remember to `gofmt` your code before any commit, it ensures style consistency and ideally you should set it up to run on save
2) **DO READ** the [Effective Go](https://golang.org/doc/effective_go.html) guide for some opinionated views on how to write idiomatic Go, straight from Google
3) **DO HANDLE** errors gracefully and properly
    * Remember to log errors at the correct log level
    * We use [logrus](https://github.com/sirupsen/logrus) to handle our logging implementation
    * Follow logging standards, mainly that logging should be done at the service and the resolver level, but never at the database repository level
    * Always propagate error messages, and give appropriate GraphQL error messages at the resolver level