# zippim: the trivial package manager

## Vision

All Windows "package management" should do is downloading a zip file,
unzipping it somewhere and making the executables available in PATH.

Subsequently, all that "windows packaging" should entail is creating a
zip file with the application and uploading it somewhere.
If your application can work with this approach, Chocolatey is drastic overkill.

This is the disease; zippim is the cure.


## Implementation

Zippim is a Go application (read: single .exe with no dependencies).


- Fetch the archive to c:\pkg\\_download
- Unzip it to c:\pkg\$PACKAGE_NAME (where $PACKAGE_NAME
is derived from the file name or specified at command line)
- Scan the unpacked archive for .exe files and create .cmd launchers to c:\pkg\bin

Usage:

- Create c:\pkg\bin, add it to the path and download zippim.exe there
- Run the command (this downloads SciTE):

```shell
$ zippim get http://www.scintilla.org/wscite366.zip
```

If you don't like the "wscite366" directory name, specify the desired name on command line:

```shell

$ zippim get --name scite http://www.scintilla.org/wscite366.zip

Downloading http://www.scintilla.org/wscite366.zip
Launcher to: C:\pkg\scite\wscite\SciTE.exe
```

Status: "seems to work". I will make the binary available when I have used it for a
while.

License: MIT
