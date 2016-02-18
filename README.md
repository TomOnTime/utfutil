# utfutil

Utilities to make it easier to use golang.org/x/text/encoding/unicode

## Dealing with UTF-16 files from Windows.

Ever have code that worked for years until you received a file from
a MS-Windows system that just didn't work at all?  Looking at a hex
dump you realize every other byte is \0.  WTF?

No, UTF.  UTF-16LE with an optional BOM.

What does all that mean?  Well, first you should read
["The Absolute Minimum Every Software Developer Absolutely, Positively Must Know About Unicode and Character Sets (No Excuses!)"](http://www.joelonsoftware.com/articles/Unicode.html) by Joel Spolsky.

However if you use Golang, you can take the easy way out.
Just change ioutil.ReadFile() to utfutil.ReadFile() and things
should just work.

### utfutil.ReadFile() is the equivalent of ioutil.ReadFile()

OLD: Works on Mac/Linux:

```
		data, err := ioutil.ReadFile(filename)
```

NEW: Works if someone gives you a Windows UTF-16LE file:

```
		data, err := utfutil.ReadFile(filename, utfutil.UTF8)
```

### utfutil.OpenFile() is the equivalent of os.Open().

OLD: Works on Mac/Linux:

```
		data, err := os.Open(filename)
```

NEW: Works if someone gives you a file with a BOM:

```
		data, err := utfutil.OpenFile(filename, utfutil.HTML5)
```

### utfutil.NewScanner() is for reading files line-by-line

It works like os.Open():

```
		s, err := utfutil.NewScanner(filename, utfutil.HTML5)
```


## Encoding hints:

What's that second argument all about?

Since it is impossible to guess 100% correctly if there is no BOM,
the functions take a 2nd parameter of type "EncodingHint" where you
specify the default encoding for BOM-less files.

```
UTF8        No BOM?  Assume UTF-8
UTF16LE     No BOM?  Assume UTF 16 Little Endian
UTF16BE     No BOM?  Assume UTF 16 Big Endian
WINDOWS = UTF16LE   (i.e. a good assumption if file is from MS-Windows)
POSIX   = UTF8      (i.e. a good assumption if file is from Unix or Unix-like systems)
HTML5   = UTF8      (i.e. a good assumption if file is from the web)
```

## Future Directions

In the future I'd like to add a default type "AUTO" which
makes some educated guesses about the default type, similar
to the uchardet command.  Pull requests gladly accepted!
