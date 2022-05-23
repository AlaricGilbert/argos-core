# argos-core
The core sniff & master node implementation of a cryptocurrency public chian monitor system called `argos`.

It should be noticed that due to the support situation of `netpoll`, code nodes cannot run on Windows systems.
## Getting Started

### Insturctions to deploy Master Node
* Setting up your Database environment
* Run `build.sh` or manually copy `master/config/config.example` to `master/config/config.go`
* Modify your database Data Source Name and your web service listen address (`:8080` default)
* Build master node and build your master node images.
* Deploy it by just execute it.
### Insturctions to deploy Sniffer Node 
* Run `build.sh` or manually build sniffer node.
* Create a `config.json` like:
```jsonc
{
    "master_address": "127.0.0.1:4222",     // Master IP:4222 (4222 is default RPC port)
    "identifier": "hubei-SIp7m1Lkc4"        // [Prefix]-[Random Unique ID]
}
```
* Build your sniffer node images (executable + json).
* Deploy it by just execute it.

### Instructions to build & develop the project
#### Install kitex compiler
* Ensure `GOPATH` environment variable is defined properly (for example `export GOPATH=~/go`), then add `$GOPATH/bin` to `PATH` environment variable (for example export `PATH=$GOPATH/bin:$PATH`). Make sure `GOPATH` is accessible.
* Install kitex: 
```sh
go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
```
* Install thriftgo: 
```sh
go install github.com/cloudwego/thriftgo@latest
```
#### Generate RPC codes
Generated codes are not included in the repository, you should run the `kitexgen.sh` in the repository folder to generate RPC codes.
```sh
$ pwd
your-code-folder/argos-core
$ sh kitexgen.sh
```
#### Build master and sniffer
```sh
$ cd sniffer
$ bash build.sh
$ ...
$ cd ../master
$ bash build.sh
```

Also you can just run `build.sh` in top folder, it would update all kitex generated codes and build both `master` and `sniffer`
## Project Structure
> //TODO
## Contribution
Pull Requests are welcomed after this repo become public.
## License

```
MIT License

Copyright (c) 2022 Alaric

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```