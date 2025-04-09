==========
Crestfiles
==========

Crest is not only a commandline tool, it is also a language. Crestfile is a compiled declarative language which allows you to preconfigure everything you want crest to do when you run it using a simple and easy to understand language. By that nature, Crestfile has support for more useful features than the traditional commandline interface would.

There are two main features you get to use when using crestfile as opposed to the traditional commandline interface. The first is exclusions, where you can exclude a path from being crawled. The second is hooks, a hook allows you to attatch a bash process to the crawling of a certain path. Say for example you wanted to to echo 'hello' to the terminal every time you crawled '/foo', in the Crestfile you would write 'hook "/foo" \`echo hello\`'

Usage
=====

How to run: ``crest run path/to/Crestfile``

How to write a crestfile ::

    anotherWorkingSite = "/AnotherWorkingSite"

    url          http://localhost:8080
    type         testHTTP
    verbose      true

    hook "/" `echo Hello`
    hook {anotherWorkingSite} `echo Foo`

In this code, the following happens:
#. There is a variable for the path "/AnotherWorkingSite" defined as anotherWorkingSite.
#. the URL to crawl is set to http://localhost:8080
#. the type of testing is set to testHTTP.
#. verbose is set to true, outputting the code in verbose mode. You can also use "quiet" for a quiet print.
#. A hook is created which echo's "Hello" when "/" is crawled.
#. Another hook is created which echos Foo when "/AnotherWorkingSite" is crawled.

Notes
=====

Crestfiles do not support single quotes for string literals.
