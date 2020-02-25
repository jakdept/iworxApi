# iworxApi
InterWorx API bindings for Golang


Uses XMLRPC:

http://xmlrpc.com/spec.md

InterWorx authentication information:

http://docs.interworx.com/interworx/api/index-Introduction.php#toc-Section--1
http://docs.interworx.com/interworx/api/index-Using-the-API.php#toc-Section-2.1

BaseURL for XMLRPC is:
```
https://%%SERVERNAME%%:2443/xmlrpc
```

InterWorx NodeWorx and SiteWorx APIs appear to not do any session persistance:

https://www.php.net/manual/en/soapserver.setpersistence.php

So we have to send credentials every time