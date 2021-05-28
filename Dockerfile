FROM alpine:latest
RUN echo -e "https://mirror.tuna.tsinghua.edu.cn/alpine/latest-stable/main\n\
https://mirror.tuna.tsinghua.edu.cn/alpine/latest-stable/community" > /etc/apk/repositories

# RUN apt-get update
# RUN apt-get install -y net-tools
# RUN apt-get install -y iputils-ping
# RUN apt-get install -y ca-certificates

RUN apk update && apk --no-cache add tzdata ca-certificates wget \
&& cp -r -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

#RUN /bin/cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' >/etc/timezone

RUN mkdir -p /dawn/
RUN mkdir -p /dawn/config
ADD ./dawn_api /dawn

WORKDIR /dawn

EXPOSE 3000
CMD ./dawn_api
