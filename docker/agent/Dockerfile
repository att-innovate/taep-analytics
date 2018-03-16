FROM python:2.7.14-alpine

COPY src /

RUN python setup.py install 

# To share Unix socket with Kapacitor
VOLUME /var/lib/kapacitor

ENTRYPOINT ["python", "agent.py"]