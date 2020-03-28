FROM mcr.microsoft.com/powershell:alpine-3.10

ADD server/server /usr/local/bin/
ADD DeathTax/DeathTax /usr/local/share/deathtax/
Add site /usr/local/share/deathtax/web

CMD ["server"]