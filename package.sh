#!/bin/bash

#gen docker image
echo $#
if [ $# -lt 1 ] ;then
    echo "please  input version V1.2.* ...."
    exit 1
fi


echo $1
#make
sudo docker build -t registry.cn-hangzhou.aliyuncs.com/meetwhale/whale-traj-api:$1 .
sudo docker login --username=xiuboye@whale registry.cn-hangzhou.aliyuncs.com
tag=`sudo docker images | grep $1 | awk '{print $3}'`
echo "tagid: " $tag
sudo docker tag $tag registry.cn-hangzhou.aliyuncs.com/meetwhale/whale-traj-api:$1
sudo docker push registry.cn-hangzhou.aliyuncs.com/meetwhale/whale-traj-api:$1
