FROM postgres:latest

COPY ./deploy/pg/init.sh /usr/local/bin/init.sh
RUN chmod +x /usr/local/bin/init.sh

EXPOSE 5432 

CMD [ "postgres" ]