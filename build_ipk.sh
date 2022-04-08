#!/bin/sh
name="ddns-go"

version=$(git describe --tags `git rev-list --tags --max-count=1`)
echo ${version}

export CGO_ENABLED=0
export GOOS="linux"
export GOARCH="arm64"
go build -ldflags="-X main.version=$version -X 'main.buildTime=`date`' -extldflags -static" -o ./bin/linux_arm64/$name ./main.go   # arm64

GOARCH="mips"
export GOMIPS="softfloat"

go build -ldflags="-X main.version=$version -X 'main.buildTime=`date`' -extldflags -static" -o ./bin/linux_mips/$name ./main.go   # mips

GOARCH="mipsle"
go build -ldflags="-X main.version=$version -X 'main.buildTime=`date`' -extldflags -static" -o ./bin/linux_mipsle/$name ./main.go   # mipsle



#!/bin/sh
name="ddns-go"
version=$(git describe --tags `git rev-list --tags --max-count=1`)
echo ${version}

mkdir -p ./ipk/opt/bin

cat>./ipk/postinst<<EOF1
#!/bin/sh
[ ! -d "/etc/init.d" ] && mkdir /etc/init.d
cat>/etc/init.d/ddns-go<<EOF
#!/bin/sh
START=99
start() {
    echo "开始启动ddns-go"
	pid=\\\`ps | grep ddns-go | grep -v 'grep' | grep -v 'start'| grep -v 'restart'|awk '{print \\\$1}' | head -n 1\\\`
    echo "pid: \\\$pid"
    pidd=\\\`ps | grep ddns-go | grep -v 'grep'\\\`
    echo "pidd: \\\$pidd"
	if [ -n "\\\$pid" ]; then
      echo "Already started!"
	else
      nohup /opt/bin/ddns-go >/dev/null 2>&1 &
      echo "Started!"
    fi
}
stop() {
    echo "开始停止ddns-go"
	pid=\\\`ps | grep ddns-go | grep -v 'grep' | grep -v 'stop'| grep -v 'restart'|awk '{print \\\$1}' | head -n 1\\\`
    echo "pid: \\\$pid"
    pidd=\\\`ps | grep ddns-go | grep -v 'grep'\\\`
    echo "pidd: \\\$pidd"
	if [ -n "\\\$pid" ]; then
        kill -9 \\\$pid
      echo "Stopped"
    else
      echo "Error! not started!" 1>&2
    fi
}
case "\\\$1" in
    start)
        start
        exit 0
    ;;
    stop)
        stop
        exit 0
    ;;
    restart)
        stop
        start
        exit 0
    ;;
    **)
        echo "Usage: \\\$0 {start|stop|reload}" 1>&2
        exit 1
    ;;
esac
EOF
chmod 755 /etc/init.d/ddns-go
chmod 777 /opt/bin/ddns-go
nohup /opt/bin/ddns-go >/dev/null 2>&1 &
EOF1

cat>./ipk/prerm<<EOF1
#!/bin/sh
[ -d "/etc/ddns-go" ] && rm -rf /etc/ddns-go
[ -e "/etc/init.d/ddns-go" ] && rm -rf /etc/init.d/ddns-go
EOF1

chmod 755 ./ipk/postinst
chmod 755 ./ipk/prerm

echo "2.0" >./ipk/debian-binary

# create mips package
/bin/cp -f ./bin/linux_mips/$name ./ipk/opt/bin/

echo "Package: ${name}" >./ipk/control
echo "Version: ${version}" >>./ipk/control
echo "Section: lang" >>./ipk/control
echo "Author: jeessy" >>./ipk/control
echo "feat by: D0raemon <labulac@88.com>" >>./ipk/control
echo "Architecture: all" >>./ipk/control
echo "Installed-Size: `stat -c "%s" ./ipk/opt/bin/$name`" >>./ipk/control
echo "Description:  简单好用的DDNS。自动更新域名解析到公网IP(支持阿里云、腾讯云dnspod、Cloudflare、华为云)" >>./ipk/control

tar -zcvf ./ipk/data.tar.gz --transform s=/ipk== ./ipk/opt
tar -zcvf ./ipk/control.tar.gz --transform s=/ipk== ./ipk/control ./ipk/postinst ./ipk/prerm
tar -zcvf ./ipk/${name}_${version}_mips.ipk --transform s=/ipk== ./ipk/data.tar.gz ./ipk/control.tar.gz ./ipk/debian-binary



# Create arm64 package
rm -f ./ipk/data.tar.gz ./ipk/control.tar.gz
/bin/cp -f ./bin/linux_arm64/$name ./ipk/opt/bin/

echo "Package: ${name}" >./ipk/control
echo "Version: ${version}" >>./ipk/control
echo "Section: lang" >>./ipk/control
echo "Author: jeessy" >>./ipk/control
echo "feat by: D0raemon <labulac@88.com>" >>./ipk/control
echo "Architecture: all" >>./ipk/control
echo "Installed-Size: `stat -c "%s" ./ipk/opt/bin/$name`" >>./ipk/control
echo "Description:  简单好用的DDNS。自动更新域名解析到公网IP(支持阿里云、腾讯云dnspod、Cloudflare、华为云)" >>./ipk/control

tar -zcvf ./ipk/data.tar.gz --transform s=/ipk== ./ipk/opt
tar -zcvf ./ipk/control.tar.gz --transform s=/ipk== ./ipk/control ./ipk/postinst ./ipk/prerm
tar -zcvf ./ipk/${name}_${version}_arm64.ipk --transform s=/ipk== ./ipk/data.tar.gz ./ipk/control.tar.gz ./ipk/debian-binary


# Create mips64le package
rm -f ./ipk/data.tar.gz ./ipk/control.tar.gz
/bin/cp -f ./bin/linux_mipsle/$name ./ipk/opt/bin/

echo "Package: ${name}" >./ipk/control
echo "Version: ${version}" >>./ipk/control
echo "Section: lang" >>./ipk/control
echo "Author: jeessy" >>./ipk/control
echo "feat by: D0raemon <labulac@88.com>" >>./ipk/control
echo "Architecture: all" >>./ipk/control
echo "Installed-Size: `stat -c "%s" ./ipk/opt/bin/$name`" >>./ipk/control
echo "Description:  简单好用的DDNS。自动更新域名解析到公网IP(支持阿里云、腾讯云dnspod、Cloudflare、华为云)" >>./ipk/control

tar -zcvf ./ipk/data.tar.gz --transform s=/ipk== ./ipk/opt
tar -zcvf ./ipk/control.tar.gz --transform s=/ipk== ./ipk/control ./ipk/postinst ./ipk/prerm
tar -zcvf ./ipk/${name}_${version}_mipsle.ipk --transform s=/ipk== ./ipk/data.tar.gz ./ipk/control.tar.gz ./ipk/debian-binary

rm -rf ./ipk/data.tar.gz ./ipk/control.tar.gz ./ipk/control ./ipk/postinst ./ipk/prerm ./ipk/opt ./ipk/debian-binary


