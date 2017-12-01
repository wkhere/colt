## colt

![](https://www.colt.com/portals/0/productimages/2015/P4750_left.png "BANG!")

`colt` is COLumn Transformer.

`colt +N command` transforms the Nth column of input with given command,
leaving other columns as they are.

`colt -N command` does the same for the Nth column counting from the end.

`colt command` transforms the last column (`-1` is the default).

`colt -d':' command` specifies the column delimiter (the default is `;`, it must be 1 character).
