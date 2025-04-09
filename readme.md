## crest
Fast, simple, and extensive testing for static pages.

---

### Installation
(1. Change `INSTALL_PATH` to the directory you want to install it in in the `Makefile`. This could be equated to PREFIX on other codebases.<br>
(2. run `make` to build and `make install` to install to the install path you specified earlier in the makefile. <br>

---

<br>Make sure to specify scheme in your URL. Only HTTP/HTTPS is recognized. Make sure that your host is localhost: if you try to crawl a site that is not localhost you will recieve an error, but one of your pages links to another site it will simply be skipped. <br>
The design philosophy of crest is permissive but contained. Meaning, it will by default crawl everything unless specificied otherwise, but it will make sure only to crawl your site. Most features that will be added to crest will follow that general idea.
<br>
