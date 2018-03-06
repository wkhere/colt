## colt

![](https://small.shootingsportsmedia.com/52544.jpg "BANG!")

`colt` is COLumn Transformer.

`colt +N command` transforms the Nth column of input with given command,
leaving other columns as they are.

`colt -N command` does the same for the Nth column counting from the end.

`colt command` transforms the last column (`-1` is the default).

`colt -s':' command` specifies the column separator
(the default is `;`, it must be 1 character).

`command` should accept input data as an argument and print the output.

NOTE: this program can be replaced by a several lines long AWK program.
The thing which would make it better over awk is for example considering
the semantics of quotes, doublequotes and backslashes - ignoring the separator
which is quoted this way.
