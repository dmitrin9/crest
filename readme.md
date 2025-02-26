## crest
Crawl all anchor tags on your site to test if any URLs return status errors when requested. 

---

### Installation
(1. Change `INSTALL_PATH` to the directory you want to install it in. <br>
(2. run `make` and `make install`. <br>
(3. If your `INSTALL_PATH` is in your PATH, then you can use the `crest` command immediately.
---

### Usage
``` bash
crest -t/--test-http http://localhost:8080 # Crawl the anchors on your site.
crest -v/--verbose -t/--test-http http://localhost:8080 # The verbose flag will modify your testing flag, so it is not a flag that has any standalone functionoality. It will print verbose logging messages for all stages in the crawling process.
```
<br>
Make sure to include the protocol in your url (for example: "http://localhost:8080" as opposed to just "localhost:8080"). Also make your url is on localhost. For now I'm only supporting crawling localhost and ignoring all anchors that lead to another site. I might support something like this in the future, but for the mean time I'm only planning on supporting a feature that has a clear usecase to me and supporting it as a feature, and not an "oh well that was nice" incidental functionality. So until I have a clear usecase, or I am told a clear usecase for allowing anchors or crawling for non-localhost sites, I will flat out not support it. This goes, as I said earlier, for anything else; nothing is an incidental "that's nice functionality", if it doesn't have a clear use I will not support it. Rant aside, the design philosophy of this tule is to have no implicit knowledge, you must define everything you are doing.
