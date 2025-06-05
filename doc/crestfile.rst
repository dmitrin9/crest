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
    depth        2
    exclude {excludeSomething}

Explanation of the above file:
==============================

excludeSomething    This is a string variable, the link in it will be excluded.
url                 keyword ``url`` is required. It defines which url will be crawled.
type                keyword ``type`` is required. It defines how you wanna test your website.
followRobots        setting this to true will obey the robots.txt policy of your website.
verbose             setting this to true will print everything happening. There is also a ``quiet`` keyword that will print in quiet mode.
depth               depth allows you to define to what depth you want to crawl.
exclude             exclude will allow you to exclude a specific path from being crawled.

Notes
=====

Crestfiles do not support single quotes for string literals.

Variables are used by wrapping the variable name in curly braces.

The difference between verbose and quiet mode: Verbose mode will print everything that is happening at each stage of the test. Quiet mode will only print errors. Crest will by default print in an inbetween state where it prints messages but not detailed ones.
