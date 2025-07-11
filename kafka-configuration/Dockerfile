FROM apache/kafka:3.9.1

COPY server.properties /opt/kafka/config/kraft/server.properties
COPY entrypoint.sh /entrypoint.sh


EXPOSE 29092
ENTRYPOINT ["/entrypoint.sh"]