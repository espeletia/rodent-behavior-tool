FROM ultralytics/ultralytics:8.3.100

WORKDIR /app
ADD ./ratt/requirements.txt /app/requirements.txt
RUN pip3 install -r requirements.txt

ADD ./ratt/ /app/

ENTRYPOINT ["python", "main.py"]

