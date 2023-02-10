# 通过命令获取 IP

## 在 Linux 系统中获取网卡 eth1 的 IPv6 地址
```sh
ip -6 addr show eth1
```

## 在 Linux 系统中获取网卡 eth1 的 IPv4 地址
```sh
ip -4 addr show eth1
```

## Linux get IPv6 prefix （移动开头2409）
```sh
Linux get IPv6 prefix （移动开头2409）
ip -6 route | awk '{print $1}' | awk '/2409:?/' | awk -F::/ '{print $1 "any other suffix of other mac"}'
```

更多请见 [#531](https://github.com/jeessy2/ddns-go/issues/531)