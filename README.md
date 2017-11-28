## colt

`colt` is COLumn TRansformer.

`colt +N subcommand` transforms the Nth column of input with subcommand,
leaving other columns as they are.

`colt -N subcommand` does the same for the Nth column counting from the end.

`colt subcommand` transforms the last column (`-1` is the default).

`colt -d':' subcommand` specifies the column delimiter (the default is `;`).

