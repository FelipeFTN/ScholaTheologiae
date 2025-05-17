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

# ---------- Final Container ----------
FROM frolvlad/alpine-glibc

# Install necessary dependencies
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
COPY nginx.conf /etc/nginx/nginx.conf

# Expose Heroku-compatible port
EXPOSE 80

# Set environment
ENV LD_LIBRARY_PATH=/libyaml/src/.libs
ENV SECRET_KEY_BASE ""
ENV RAILS_ENV "production"
ENV PORT "443"

# Let's build libyaml from scratch before installing dependencies

RUN apk add --no-cache \
  autoconf \
  automake \
  build-base \
  cmake \
  git \
  libtool \
  perl-dev \
  ruby-dev

RUN mkdir -p /usr/src/libyaml
WORKDIR /usr/src/libyaml
# I will just wget the latest release i found
RUN wget https://github.com/yaml/libyaml/releases/download/0.2.5/yaml-0.2.5.tar.gz
RUN tar -xzf yaml-0.2.5.tar.gz
WORKDIR /usr/src/libyaml/yaml-0.2.5
RUN ./configure --prefix=/usr
RUN make && make install
RUN rm -rf /usr/src/libyaml

WORKDIR /app

RUN cd /app/rails && bundle install --gemfile Gemfile

# Start all services
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]
