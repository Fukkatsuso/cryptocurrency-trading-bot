# https://hub.docker.com/layers/mysql/library/mysql/5.7/images/sha256-95a3b244b2e990ce5633dd449326cec6ce20e9e366ba01fdf11ebc15a9cd38d3

FROM mysql:5.7

ADD my.cnf etc/mysql/conf.d/my.cnf
RUN chmod 644 /etc/mysql/conf.d/my.cnf

EXPOSE 3306
