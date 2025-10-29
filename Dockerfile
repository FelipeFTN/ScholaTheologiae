# ---------- Build API ----------
FROM golang:latest AS api-build

WORKDIR /app/api
COPY ./api /app/api
COPY ./books/* /app/books/

# Install sqlite3
RUN apt-get update && apt-get install -y libsqlite3-dev

RUN make build

# DEBUG: This creates a dummy binary that prints "Hello World!" to skip the API build
# RUN mkdir bin/ && echo "#!/bin/bash\necho 'Hello World!'" > bin/schola-theologiae-api && chmod +x bin/schola-theologiae-api

# ---------- Final Container ----------
FROM ubuntu:latest

# Install necessary dependencies
RUN apt update && apt install -y \
  bash \
  nginx \
  supervisor \
  ruby \
  ruby-dev \
  build-essential \
  sqlite3 \
  libtool \
  perl-base \
  ruby-dev \
  openssl \
  make \
  libyaml-dev

# Set working directory
WORKDIR /app

# Copy API binary & Data
RUN mkdir -p /app/api
COPY --from=api-build /app/api/bin/schola-theologiae-api ./api/schola-theologiae-api 
COPY --from=api-build /app/api/data ./data

# Copy Rails app
RUN mkdir -p /app/rails
COPY ./app /app/rails

WORKDIR /app/rails

# Install Rails and Bundler stuff
ENV RAILS_ENV="production"
ARG RAILS_ENV="production"
RUN gem install rails bundler

# Install gems
RUN bundle config set without development test
RUN bundle install

# Generate secret key base
RUN EDITOR="echo --wait" bin/rails credentials:edit

# Move assets subfolders into assets/*
RUN find /app/rails/app/assets -type f -exec cp {} /app/rails/public/ \;

# Don't know if this is necessary, so i will comment it out for now
# Update: It is necessary.
RUN /app/rails/bin/rails assets:precompile --trace

WORKDIR /app

# Copy supervisor config
COPY supervisord.conf /etc/supervisord.conf

# Copy Nginx config
COPY ./nginx.conf /etc/nginx/nginx.conf

# Expose Heroku-compatible port
ARG PORT=8000
ENV PORT=8000
EXPOSE $PORT

# Start all services
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]
