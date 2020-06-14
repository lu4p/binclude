package main

import (
	b "binclude"
	"time"
)

const (
	_binclude0	= `H4sIAAAAAAAA/youTUrJLEosLk4tMQQEAAD//zzpe5cMAAAA`
	_binclude1	= `H4sIAAAAAAAA/youTUrJLEosLk4tMQQEAAD//zzpe5cMAAAA`
	_binclude2	= `H4sIAAAAAAAA/0rLzEnVK6koAQQAAP//JRb34AgAAAA=`
	_binclude3	= `H4sIAAAAAAAA/0osLk4tMQQEAAD//zJdA+EGAAAA`
	_binclude4	= `H4sIAAAAAAAA/0osLk4tMQQEAAD//zJdA+EGAAAA`
)

var binFS = b.FileSystem{"assets/subdir/subdirasset1.txt": {Filename: "subdirasset1.txt", Mode: 420, ModTime: time.Unix(1592104027, 0), Content: b.StrToByte(_binclude0)}, "assets/subdir/subdirasset2.txt": {Filename: "subdirasset2.txt", Mode: 420, ModTime: time.Unix(1592104033, 0), Content: b.StrToByte(_binclude1)}, "file.txt": {Filename: "file.txt", Mode: 420, ModTime: time.Unix(1592104011, 0), Content: b.StrToByte(_binclude2)}, "assets": {Filename: "assets", Mode: 2147484141, ModTime: time.Unix(1592097652, 0), Content: nil}, "assets/asset1.txt": {Filename: "asset1.txt", Mode: 420, ModTime: time.Unix(1592103949, 0), Content: b.StrToByte(_binclude3)}, "assets/asset2.txt": {Filename: "asset2.txt", Mode: 420, ModTime: time.Unix(1592103958, 0), Content: b.StrToByte(_binclude4)}, "assets/subdir": {Filename: "subdir", Mode: 2147484141, ModTime: time.Unix(1592104033, 0), Content: nil}}
