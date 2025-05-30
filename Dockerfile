# ---------- Build API ----------
FROM golang:latest as api-build

WORKDIR /app/api
COPY ./api /app/api
COPY ./books/* /app/books/

RUN make build

# ---------- Build Rails ----------
FROM ruby:3.2 as rails-build

WORKDIR /app/rails
COPY ./app/Gemfile ./app/Gemfile.lock ./
RUN gem install bundler && bundle install

COPY ./app .

# Move assets subfolders into assets/*
RUN find /app/rails/app/assets -type f -exec cp {} /app/rails/public/ \;

RUN /app/rails/bin/rails assets:precompile

# ---------- Final Container ----------
FROM frolvlad/alpine-glibc

# Install necessary dependencies
RUN apk add \
  bash \
  nginx \
  supervisor \
  ruby \
  ruby-bundler \
  ruby-json \
  ruby-irb \
  libstdc++ \
  libffi-dev \
  build-base \
  libressl \
  linux-headers \
  sqlite-libs \
  tzdata

# Set working directory
WORKDIR /app

RUN mkdir -p /app/api

# Copy API binary & Data
COPY --from=api-build /app/api/bin/schola-theologiae-api ./api/schola-theologiae-api 
COPY --from=api-build /app/api/data ./data
RUN ls /app/api

# Copy Rails app
COPY --from=rails-build /app/rails /app/rails

# Copy supervisor config
COPY supervisord.conf /etc/supervisord.conf

# Copy Nginx config
COPY ./nginx.conf /etc/nginx/nginx.conf

# Expose Heroku-compatible port
ARG PORT
ENV ENV_PORT $PORT
EXPOSE $PORT

# Set environment
ENV LD_LIBRARY_PATH=/libyaml/src/.libs
ENV SECRET_KEY_BASE ""
ENV RAILS_ENV "production"

RUN apk add \
  libtool \
  perl-dev \
  ruby-dev

WORKDIR /app

# Copy compiled libyaml
COPY libyaml.tar.gz /tmp/libyaml.tar.gz
RUN tar -xzf /tmp/libyaml.tar.gz -C /usr && rm /tmp/libyaml.tar.gz

RUN cd /app/rails && bundle install --gemfile Gemfile

# Replace nginx port
RUN sed -i "s/\${PORT}/$PORT/g" /etc/nginx/nginx.conf

# Start all services
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]
