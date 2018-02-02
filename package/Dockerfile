FROM ubuntu:16.04

COPY rds-broker /usr/bin/

ADD https://kubernetes-charts.storage.googleapis.com/mysql-0.3.4.tgz /mysql-0.3.4.tgz
ADD https://kubernetes-charts.storage.googleapis.com/mariadb-2.1.3.tgz /mariadb-2.1.3.tgz
ADD https://kubernetes-charts.storage.googleapis.com/postgresql-0.8.8.tgz /postgresql-0.8.8.tgz

ENTRYPOINT ["rds-broker"]
