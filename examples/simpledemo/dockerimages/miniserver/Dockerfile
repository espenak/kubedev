FROM python:3.6-alpine

ENV APP_ROOT_PATH /usr/src/app
RUN mkdir -p $APP_ROOT_PATH

WORKDIR $APP_ROOT_PATH
COPY server.py .

EXPOSE 8000

CMD [ "python", "server.py" ]
