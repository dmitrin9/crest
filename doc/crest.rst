=====
Crest
=====

This is the documentation for the crest commandline interface.

Usage
=====

Usage: ``crest [options] url``

-v, --verbose        Print in verbose mode.
-q, --quiet          Print in quiet mode.
-f, --follow-robots  Follow robots.txt policy.
-t, --test-http      Test HTTP.

Notes
=====

Make sure to specify scheme in your URL. Only HTTP/HTTPS is recognized. Make sure that your host is localhost: if you try to crawl a site that is not localhost you will recieve an error, but one of your pages links to another site it will simply be skipped.
The design philosophy of crest is permissive but contained. Meaning, it will by default crawl everything unless specificied otherwise, but it will make sure only to crawl your site. Most features that will be added to crest will follow that general idea.
