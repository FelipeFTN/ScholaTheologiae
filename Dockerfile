# ---------- Build API ----------
FROM golang:1.21-alpine as api-build

WORKDIR /app/api

# Copy Go module files first for better caching
COPY ./api/go.mod ./api/go.sum ./
RUN go mod download

# Copy the rest of the API source code
COPY ./api .
COPY ./books/* /app/books/

RUN make build

# ---------- Build Rails ----------
FROM ruby:3.2-alpine as rails-build

WORKDIR /app/rails

# Install build dependencies
RUN apk add --no-cache build-base nodejs yarn sqlite-dev

# Copy Gemfile and Gemfile.lock first for better caching
COPY ./app/Gemfile ./app/Gemfile.lock ./
RUN gem install bundler && bundle install --without development test

# Copy the rest of the Rails application
COPY ./app .

# Move assets subfolders into assets/*
RUN find /app/rails/app/assets -type f -exec cp {} /app/rails/public/ \; 2>/dev/null || true

# Precompile assets
RUN RAILS_ENV=production /app/rails/bin/rails assets:precompile

# ---------- Final Container ----------
FROM frolvlad/alpine-glibc

# Install necessary dependencies in a single layer
RUN apk add --no-cache \
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
  libtool \
  perl-dev \
  ruby-dev \
  tzdata \
  && rm -rf /var/cache/apk/*

# Set working directory
WORKDIR /app

# Copy API binary & Data
RUN mkdir -p /app/api
COPY --from=api-build /app/api/bin/schola-theologiae-api ./api/schola-theologiae-api 
COPY --from=api-build /app/api/data ./data

# Copy Rails app
COPY --from=rails-build /app/rails /app/rails

# Copy supervisor config
COPY supervisord.conf /etc/supervisord.conf

# Copy Nginx config
COPY ./nginx.conf /etc/nginx/nginx.conf

# Expose Heroku-compatible port
ARG PORT 8000
ENV PORT 8000
ENV ENV_PORT $PORT
EXPOSE $PORT

# Set environment
ENV LD_LIBRARY_PATH=/libyaml/src/.libs
ENV SECRET_KEY_BASE ""
ENV RAILS_ENV "production"

# Copy compiled libyaml
COPY libyaml.tar.gz /tmp/libyaml.tar.gz
RUN tar -xzf /tmp/libyaml.tar.gz -C /usr && rm /tmp/libyaml.tar.gz

RUN cd /app/rails && bundle install --gemfile Gemfile

# Replace nginx port
# RUN sed -i "s/\${PORT}/$PORT/g" /etc/nginx/nginx.conf

# Start all services
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]
