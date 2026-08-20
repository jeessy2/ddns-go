# DDNS-GO For iOS

Add support for iOS 12 and above systems

## English

### Usage

1. Download the latest version
2. Jailbreak the device (both Rootless and Rootful)
3. Download the latest version and place it in `/var/mobile/ddns-go`
4. Install NewTerm or other terminals
5. Enter `su` to obtain permission
6. Enter:

```bash
cat > /tmp/ent.xml << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>platform-application</key>
    <true/>
    <key>com.apple.security.cs.disable-executable-page-protection</key>
    <true/>
    <key>com.apple.security.cs.allow-unsigned-executable-memory</key>
    <true/>
    <key>com.apple.security.cs.disable-library-validation</key>
    <true/>
    <key>network-client</key>
    <true/>
</dict>
</plist>
EOF

ldid -S/tmp/ent.xml /rootfs/private/var/mobile/ddns-go/ddns-go

chmod +x /rootfs/private/var/mobile/ddns-go/ddns-go
/rootfs/private/var/mobile/ddns-go/ddns-go -l :9876 -f 600
```

## 中文

### 使用步骤

1. 下载最新版本 如果下载压缩包 要执行解压 得到一个名为ddns-go的可执行文件
2. 越狱设备 我测试时使用的是RootHide Dopemine
3. 使用Filza或者其它工具把可执行文件放到 `/var/mobile/ddns-go/`
4. 下载一个终端插件 如NewTerm2
5. 输入 `su` 获得权限(一般是alpine)
6. 完整复制下方代码到终端:

```bash
cat > /tmp/ent.xml << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>platform-application</key>
    <true/>
    <key>com.apple.security.cs.disable-executable-page-protection</key>
    <true/>
    <key>com.apple.security.cs.allow-unsigned-executable-memory</key>
    <true/>
    <key>com.apple.security.cs.disable-library-validation</key>
    <true/>
    <key>network-client</key>
    <true/>
</dict>
</plist>
EOF

ldid -S/tmp/ent.xml /rootfs/private/var/mobile/ddns-go/ddns-go

chmod +x /rootfs/private/var/mobile/ddns-go/ddns-go
/rootfs/private/var/mobile/ddns-go/ddns-go -l :9876 -f 600
```
