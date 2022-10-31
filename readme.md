# Containers using only chroot

This is demo-ware and probably shouldn't be used.

Requirements:

 * Linux OS

## Usage

* Build: `go build chrun.go`
* Pull Image: `chrun pull <imagename>`
* Run Image: `chrun run <imagename> <entry-point>`

## Example

Pull image
```
> ./chrun pull redis
Pulling image redis
export image 16b87aa63c8f3a1e14a50feb94cba39eaa5d19bec64d90ff76c3ded058ad09c8
```

Run Redis server with chroot:
```
> chrun run redis "/usr/local/bin/redis-server"

Running /usr/local/bin/redis-server in /tmp/_assets_redis_tar_gz4234401501
4360:C 31 Oct 2022 16:07:57.253 # oO0OoO0OoO0Oo Redis is starting oO0OoO0OoO0Oo
4360:C 31 Oct 2022 16:07:57.253 # Redis version=7.0.5, bits=64, commit=00000000, modified=0, pid=4360, just started
4360:C 31 Oct 2022 16:07:57.253 # Warning: no config file specified, using the default config. In order to specify a config file use /usr/local/bin/redis-server /path/to/redis.conf
4360:M 31 Oct 2022 16:07:57.256 * Increased maximum number of open files to 10032 (it was originally set to 1024).
4360:M 31 Oct 2022 16:07:57.256 * monotonic clock: POSIX clock_gettime
                _._                                                  
           _.-``__ ''-._                                             
      _.-``    `.  `_.  ''-._           Redis 7.0.5 (00000000/0) 64 bit
  .-`` .-```.  ```\/    _.,_ ''-._                                  
 (    '      ,       .-`  | `,    )     Running in standalone mode
 |`-._`-...-` __...-.``-._|'` _.-'|     Port: 6379
 |    `-._   `._    /     _.-'    |     PID: 4360
  `-._    `-._  `-./  _.-'    _.-'                                   
 |`-._`-._    `-.__.-'    _.-'_.-'|                                  
 |    `-._`-._        _.-'_.-'    |           https://redis.io       
  `-._    `-._`-.__.-'_.-'    _.-'                                   
 |`-._`-._    `-.__.-'    _.-'_.-'|                                  
 |    `-._`-._        _.-'_.-'    |                                  
  `-._    `-._`-.__.-'_.-'    _.-'                                   
      `-._    `-.__.-'    _.-'                                       
          `-._        _.-'                                           
              `-.__.-'                                               

4360:M 31 Oct 2022 16:07:57.260 # Server initialized
4360:M 31 Oct 2022 16:07:57.265 * Ready to accept connections
```

While that's running connect to it in another chroot:
```
> chrun run redis "/usr/local/bin/redis-cli"

Running /usr/local/bin/redis-cli in /tmp/_assets_redis_tar_gz1366317376
127.0.0.1:6379> SET mykey "Hello\nWorld"
OK
127.0.0.1:6379> GET mykey
"Hello\nWorld"
127.0.0.1:6379> 
127.0.0.1:6379> exit
```

And there you go, containers using only chroot.