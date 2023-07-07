FROM alpine:3.10

ENV PATH .:$PATH
ENV APP_NAME user

COPY usercenter /app/
COPY config/* /app/config/


WORKDIR /app

EXPOSE 15001
#ENTRYPOINT ["ls"]
CMD ["usercenter"]
#CMD ["tail -f usercenter"]