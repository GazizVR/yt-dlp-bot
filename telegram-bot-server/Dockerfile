FROM alpine:latest
ARG BUILD_JOBS
WORKDIR /app
RUN apk update && apk add alpine-sdk linux-headers zlib-dev openssl-dev gperf cmake git
RUN git clone --recursive https://github.com/tdlib/telegram-bot-api.git
WORKDIR /app/telegram-bot-api
RUN rm -rf build && mkdir build
WORKDIR /app/telegram-bot-api/build
RUN cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX:PATH=/usr/local ..
RUN cmake --build . -j$BUILD_JOBS --target install
CMD ["sh", "-c", "/usr/local/bin/telegram-bot-api --local --http-port $HTTP_PORT"]
