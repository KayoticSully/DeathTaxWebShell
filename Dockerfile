FROM mcr.microsoft.com/powershell:alpine-3.10

ADD server/server /usr/local/bin/
ADD DeathTax/DeathTax /usr/local/bin/DeathTax

CMD ["server"]