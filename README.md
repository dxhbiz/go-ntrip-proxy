### ntrip-proxy
An ntrip forwarding tool, the main purpose is to be able to support multiple devices at the same time through an ntrip account to obtain rtk correction data.

### Usage
1. Build
```bash
make
```
If all goes well, it will generate the ntrip-proxy executable in the bin folder of the current directory.

2. Configuration
The default configuration file is in config/config.json in the same level as the executable. A reference example is as follows
```json
{
  "server": {
    "host": "0.0.0.0",
    "port": 2101
  },
  "casters": [
    {
      "name": "SWIFT",
      "host": "na.l1l2.skylark.swiftnav.com",
      "port": 2101,
      "username": "Your-Username",
      "password": "Your-Password",
      "mountpoint": "RTK-RTCM31"
    }
  ],
  "log": {
    "development": true,
    "level": "info",
    "filename": "./logs/ntrip-proxy.log",
    "maxSize": 16,
    "maxBackups": 30,
    "maxAge": 7,
    "compress": false
  }
}
```

3. Run
```bash
./ntrip-proxy

# or
./ntrip-proxy --config the/path/of/config.json
```

4. Help
```bash
./ntrip-proxy -h
```

### Workflow diagram
![ntrip-proxy](./images/ntrip-proxy.png)

