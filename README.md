B2SUM
=====

Utility to calculate BLAKE2 checksums written in Go.

Note: The original utility available from https://blake2.net is faster.


Installation
------------

From source, if you have Go installed:

	$ go get github.com/dchest/b2sum

TODO: binary downloads.

Usage
-----

	b2sum [-a HASH] [-s=SIZE] [filename1] [filename2] ...

Currently the only HASH supported is "blake2b".
Default SIZE is 64.

If no filenames specified, reads from the standard input.


Examples
--------

	$ echo -n "Hello world" | b2sum
	6ff843ba685842aa82031d3f53c48b66326df7639a63d128974c5c14f31a0f33343a8c65551134ed1ae0f2b0dd2bb495dc81039e3eeb0aa1bb0388bbeac29183

	$ echo -n "Hello world" | b2sum -s=20
	5ad31b81fc4dde5554e36af1e884d83ff5b24eb0

	$ b2sum -s=32 /bin/sh /etc/bashrc 
	BLAKE2b-32 (/bin/sh) = 376f70f4acc6e204ed9d098ce0e93798cb7ed1b047686872c7f496d02364c85c
	BLAKE2b-32 (/etc/bashrc) = 1572d4fe68a18bae127fe79fa5d009fdb2e3357c238f722109012fa739aaacb7

	$ time b2sum -s=20 FreeBSD-9.0-RELEASE-amd64-disc1.iso 
	BLAKE2b-20 (FreeBSD-9.0-RELEASE-amd64-disc1.iso) = 4174862a104245d26b61315b80af92892f4a45f2

	real	0m3.358s
	user	0m3.034s
	sys	0m0.266s