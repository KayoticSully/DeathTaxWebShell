FROM mcr.microsoft.com/powershell:alpine-3.10

ADD server/server /usr/local/bin/

CMD ["server"]