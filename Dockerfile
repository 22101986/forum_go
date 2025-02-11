FROM ubuntu

WORKDIR /forum

COPY . .

ENV PATH="$PATH:/forum"

CMD ["fishrum"]