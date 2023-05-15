# Cusi
Cusi is command line tools for M5Stack MicroPython (UIFlow) system.
This command provides the functionality to read/write files to and from M5Stack device.

# Features
* multi-platform
* works with a single binary file
* list available serial ports
* list files or derectories on the device
* upload local files to the device
* download or view files on the device
* delete files on the device

# How to use
## Setup
1. Download a file named `cusi-vN.N.N.zip` (`N` is a digit. e.g. `cusi-v1.0.0.zip`) from the [release page](https://github.com/zuku/cusi/releases/latest).
2. Unarchive the downloaded ZIP file.
3. Find the command file (`cusi` or `cusi.exe`) for the platform you are using.

|Directory      |Platform         |
|---------------|-----------------|
|`darwin_amd64` |Intel Mac        |
|`darwin_arm64` |Apple Silicon Mac|
|`linux_amd64`  |Linux (x86_64)   |
|`windows_amd64`|Windows (x86_64) |

When downloading the ZIP file, your clever web browser may warn you that the file is dangerous.
If you are concerned about it, you can download the source code from GitHub and build it using your Go environment.

### macOS
On macOS, Gatekeeper blocks unidentified developer's program execution.
If you want to execute `cusi` command, follow the steps below.

1. Show the command file in Finder.
2. Control-click (or right-click) the command file, then choose _Open_ from the shortcut menu.
3. Terminal.app opens and the command is executed, then close the window.
4. Once the command has been executed, you can execute the command from your prefer terminal app.

## Usage

### Connect to the device

#### macOS
List available serial ports.

```
$ cusi -l
/dev/cu.Bluetooth-Incoming-Port
/dev/cu.usbserial-XXXXXXXXXX
/dev/cu.wlan-debug
/dev/tty.Bluetooth-Incoming-Port
/dev/tty.usbserial-XXXXXXXXXX
/dev/tty.wlan-debug
```
Use `/dev/tty.usbserial-XXXXXXXXXX` to connect to the M5Stack device.
`XXXXXXXXXX` part is a hexadecimal string associated with the device.

```
$ cusi /dev/tty.usbserial-XXXXXXXXXX
```

#### Windows
List available serial ports.

```
$ cusi.exe -l
COMX
```
`X` part of `COMX` is a digit. For example, `COM1`, `COM3` etc.

```
$ cusi.exe COMX
```

### Prompt
When the device is connected, a prompt will be displayed.

```
>
```
#### List files or directories
```
> ls
apps
blocks
boot.py
emojiImg
img
main.py
res
temp.py
test.py
update
```

#### Upload file
```
> put /path/to/my_app.py apps/my_app.py
Uploading...
1234 / 1234 bytes
```

#### List in directory
```
> ls apps
my_app.py
```

#### exit
```
> exit
```

#### help
Type `help` for more information,
```
> help
```

# License
The cusi is released under the MIT license. See [LICENSE](./LICENSE).
In addition, see [THIRD-PARTY-NOTICES.txt](./THIRD-PARTY-NOTICES.txt) for the licenses of the third-party llibraries or other resources used by the cusi.
