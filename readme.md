## crest
Very basic static page testing for small personal sites.

---

### Installation
#### (1. Change `INSTALL_PATH` to the directory you want to install it in in the `Makefile`. This could be equated to PREFIX on other codebases.<br>
#### (2. run `make` to build and `make install` to install to the install path you specified earlier in the makefile. <br>
---

### Usage
``` bash
crest -t/--test-http http://localhost:8080 # Crawl the anchors on your site.
crest -v/--verbose -t/--test-http http://localhost:8080 # The verbose flag will modify your testing flag, so it is not a flag that has any standalone functionoality. It will print verbose logging messages for all stages in the crawling process.
crest -q/--quiet -t/--test-http http://localhost:8080 # The quiet flag will print messages sent to stderr.
crest -f/--follow-robots -t/--test-http http://localhost:8080 # Follow the rules from the robots.txt.

# Flags can be combined together.
crest -tfv http://localhost:8080 # Test-http follow-robots verbose.
```
<br>
Make sure to specify scheme in your URL. Only HTTP/HTTPS is recognized. Make sure that your host is localhost: if you try to crawl a site that is not localhost you will recieve an error, but one of your pages links to another site it will simply be skipped. <br>
The design philosophy of crest is permissive but contained. Meaning, it will by default crawl everything unless specificied otherwise, but it will make sure only to crawl your site. Most features that will be added to crest will follow that general idea.
<br>
