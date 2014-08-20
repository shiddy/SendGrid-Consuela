SendGrid-Consuela
=================
Consuela is a tool that purges your list of all Unsubscribes, Bounces, Invalids, Blocks, and Spam Reports using the SendGrid web API.

Version
----
1.0

Installation
--------------
I am assuming you have a working verion of GO installed on your machine. Otherwise check out:
https://golang.org/doc/install

This tool was designed to be modular and pluged into other processes. However if you intend to use it as a standalone script you can do so by changing the "package consuela" to "package main" on line 1 of consulea.go

```sh
git clone git@github.com:shiddy/SendGrid-Consuela.git
cd SendGrid-Consuela
```

You will also need to change the hardcoded SendGrid username and password on lines 135 and 156 to your SendGrid API credentials

You then can run `go run consuela.go`

Testing
----
I've included an Example .csv included in base directory.

You can add emails with the following API call

`"https://api.sendgrid.com/api/unsubscribes.add.json?api_user=&api_key=&email="`

There is also a standard testing _test.go file included as well. This pulls from the "testingResources" directory for example emails. you can run it with 

`go test`

however if you changed yoru class from consuela to main this will not work

Update!!!
----
SendGird has implemented this tool at https://sendgrid.com/listassist

You can do all the fun things there with a pretty UI.
