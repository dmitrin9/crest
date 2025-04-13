==========
Crestfiles
==========

Crest is not only a commandline tool, it is also a language. Crestfile is a compiled declarative language which allows you to preconfigure everything you want crest to do when you run it using a simple and easy to understand language. By that nature, Crestfile has support for more useful features than the traditional commandline interface would.

Usage
=====

How to run: ``crest run path/to/Crestfile``

How to write a crestfile ::

    excludeSomething = "/foo"
    url          http://localhost:8080
    type         testHTTP
    followRobots true
    verbose      true
    exclude {foo}

Notes
=====

Crestfiles do not support single quotes for string literals.

Variables are used by wrapping the variable name in curly braces.
