FROM alpine:latest
LABEL developers="Polo"
LABEL code="v1"
LABEL lead="Polo <sespolo@gmail.com>"
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
#

COPY ./harpo ./
COPY ./harpo.yml ./harpo.yml

EXPOSE 8080
RUN chmod +x /harpo
CMD [ "/harpo" ]