FROM centos

MAINTAINER timwang

LABEL version="1.0" description="redis_orm_workbench" by="timwang"

RUN echo "mkdir redis_orm_workbench"
RUN mkdir /usr/local/redis_orm_workbench
RUN mkdir /usr/local/redis_orm_workbench/config
RUN mkdir /usr/local/redis_orm_workbench/views
RUN mkdir /usr/local/redis_orm_workbench/static

COPY redis_orm_workbench /usr/local/redis_orm_workbench/
COPY config/basic.conf /usr/local/redis_orm_workbench/config/
COPY build/local.conf /usr/local/redis_orm_workbench/config/
COPY views /usr/local/redis_orm_workbench/views
COPY static /usr/local/redis_orm_workbench/static
COPY build/start.sh /usr/local/redis_orm_workbench/

WORKDIR /usr/local/redis_orm_workbench
EXPOSE 8881

ENTRYPOINT ["./redis_orm_workbench"]


