# Entware mirror client

### A spider client for clone the official entware package site


**A golang project**

# Unoffical mirror (Host @ China Mainland)
## Entware-ng

This is software repository for network attached storages, routers and other embedded devices.

Browse through ~1800 packages for different platforms:

* armv5 - http://entware.mirrors.ligux.com/binaries/armv5/Packages.html
* armv7 - http://entware.mirrors.ligux.com/binaries/armv7/Packages.html
* mipsel - http://entware.mirrors.ligux.com/binaries/mipsel/Packages.html

## Usage

Take armv7 as the example
Modify /opt/etc/opkg.conf:

<pre>
src/gz packages http://entware.mirrors.ligux.com/binaries/armv7
dest root /
dest ram /opt/tmp
lists_dir ext /opt/var/opkg-lists
option tmp_dir /opt/tmp
</pre>

Save the file and execute `opkg update`

Done!

#HOW-TO

### Step 1

`
go get github.com/flzyup/entwaremirror
`

### Step 2
`
go build github.com/flzyup/entwaremirror
`

### Step 3

Copy the 

`
$GOPATH/src/github.com/flzyup/entwaremirror/config.toml 
`

file with the same location as entwaremirror binary

### Step 4

Modify config.toml and Run binaries, HAVE FUN!