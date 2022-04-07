#!/bin/sh
name="ddns-go"

version=$(git describe --tags `git rev-list --tags --max-count=1`)

export CGO_ENABLED=0

export GOOS="linux"
export GOARCH="mips"
export GOMIPS="softfloat"

go build -ldflags "-s -w -extldflags -static" -o ./bin/linux_mips/$name ./main.go   # mips

GOARCH="mipsle"
go build -ldflags "-s -w -extldflags -static" -o ./bin/linux_mipsle/$name ./main.go   # mipsle

mkdir ./release

tar -cvf ./release/$name-linux_mips-$version.tar --transform s=./bin/linux_mips/== ./bin/linux_mips/$name
tar -cvf ./release/$name-linux_mipsle-$version.tar --transform s=./bin/linux_mipsle/== ./bin/linux_mipsle/$name


#!/bin/sh
name="ddns-go"
version=$(git describe --tags `git rev-list --tags --max-count=1`)

mkdir -p ./ipk/opt/bin

cat>./ipk/postinst<<EOF1
#!/bin/sh
[ ! -d "/etc/ddns-go" ] && mkdir /etc/ddns-go

[ ! -d "/etc/init.d" ] && mkdir /etc/init.d
cat>/etc/init.d/ddns-go<<EOF
#!/bin/sh
START=99
start() {
    echo "begin start"
	pid=\\\`ps | grep ddns-go | grep -v 'grep' | awk '{print \\\$1}' | head -n 1\\\`
	if [ -n "\\\$pid" ]; then
      echo "Already started!"
	else
	  /opt/bin/ddns-go -D
    fi
}
stop() {
    echo "begin stop"
	pid=\\\`ps | grep ddns-go | grep -v 'grep' | awk '{print \\\$1}' | head -n 1\\\`
    if [ -n "\\\$pid" ]; then
	  kill -9 \\\$pid
      echo "stopped"
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
    reload|restart|force-reload)
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
EOF1

cat>./ipk/prerm<<EOF1
#!/bin/sh
[ -d "/etc/ddns-go" ] && rm -rf /etc/ddns-go
[ -e "/etc/init.d/ddns-go" ] && rm -rf /etc/init.d/ddns-go
EOF1

chmod 755 ./ipk/postinst
chmod 755 ./ipk/prerm

echo "2.0" >./ipk/debian-binary


/bin/cp -f ./bin/linux_mips/$name ./ipk/opt/bin/

echo "Package: ${name}" >./ipk/control
echo "Version: ${version}" >>./ipk/control
echo "Section: lang" >>./ipk/control
echo "Auther: jeessy" >>./ipk/control
echo "feat by: D0raemon <labulac@88.com>" >>./ipk/control
echo "Architecture: all" >>./ipk/control
echo "Installed-Size: `stat -c "%s" ./ipk/opt/bin/$name`" >>./ipk/control
echo "Description:  简单好用的DDNS。自动更新域名解析到公网IP(支持阿里云、腾讯云dnspod、Cloudflare、华为云)" >>./ipk/control

tar -zcvf ./ipk/data.tar.gz --transform s=/ipk== ./ipk/opt
tar -zcvf ./ipk/control.tar.gz --transform s=/ipk== ./ipk/control ./ipk/postinst ./ipk/prerm
tar -zcvf ./ipk/${name}_${version}_mips.ipk --transform s=/ipk== ./ipk/data.tar.gz ./ipk/control.tar.gz ./ipk/debian-binary


rm -f ./ipk/data.tar.gz ./ipk/control.tar.gz
/bin/cp -f ./bin/linux_mipsle/$name ./ipk/opt/bin/

echo "Package: ${name}" >./ipk/control
echo "Version: ${version}" >>./ipk/control
echo "Section: lang" >>./ipk/control
echo "Auther: jeessy" >>./ipk/control
echo "Maintainer: D0raemon <labulac@88.com>" >>./ipk/control
echo "Architecture: all" >>./ipk/control
echo "Installed-Size: `stat -c "%s" ./ipk/opt/bin/$name`" >>./ipk/control
echo "Description:  简单好用的DDNS。自动更新域名解析到公网IP(支持阿里云、腾讯云dnspod、Cloudflare、华为云)" >>./ipk/control

tar -zcvf ./ipk/data.tar.gz --transform s=/ipk== ./ipk/opt
tar -zcvf ./ipk/control.tar.gz --transform s=/ipk== ./ipk/control ./ipk/postinst ./ipk/prerm
tar -zcvf ./ipk/${name}_${version}_mipsle.ipk --transform s=/ipk== ./ipk/data.tar.gz ./ipk/control.tar.gz ./ipk/debian-binary

rm -rf ./ipk/data.tar.gz ./ipk/control.tar.gz ./ipk/control ./ipk/postinst ./ipk/prerm ./ipk/opt ./ipk/debian-binary
rm -rf ./bin ./release
