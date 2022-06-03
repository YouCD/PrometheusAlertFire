#!/bin/bash

workdir=${PWD}
download_url=`curl  -s  https://api.github.com/repos/youcd/PrometheusAlertFireFront/releases/latest|grep browser_download_url|awk -F"\"" '{print $4}'`
wget -q $download_url -O  dist.txz


tar Jxf dist.txz -C ./
rm -rf dist.txz
cd $workdir

#wget https://github.com/upx/upx/releases/download/v3.96/upx-3.96-amd64_linux.tar.xz
#tar xf upx-3.96-amd64_linux.tar.xz
#mv upx-3.96-amd64_linux/upx ./