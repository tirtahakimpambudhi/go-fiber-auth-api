FROM golang:1.22.10-alpine AS builder

ARG username="appuser"
ARG unique_id=3125

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && \
    adduser \
            --disabled-password \
            --gecos "" \
            --home "/nonexistent" \
            --shell "/sbin/nologin" \
            --no-create-home \
            --uid "${unique_id}" \
            "${username}"

COPY --chown=${username}:${username} . .

RUN apk --no-cache add ca-certificates make \
    && make build APP_NAME="main" \
    && chmod a+x /app/build/main \
    && chown -R ${username}:${username} /app

FROM gcr.io/distroless/static-debian12 AS final

ARG username

# Set environment variables (can be overridden during runtime)
ENV PORT="" \
    CACHE_DB_NAME=0 \
    CACHE_DB_HOST="" \
    CACHE_DB_PORT=6379 \
    CACHE_DB_USER="" \
    CACHE_DB_PASS="" \
    CACHE_DB_MAX_CON=100 \
    CACHE_DB_MIN_CON=10 \
    CACHE_DB_MAX_TIME=10 \
    CACHE_DB_MIN_TIME=2 \
    MODEL_PATH=resource/model \
    FIBER_HOST="" \
    FIBER_PORT=0 \
    FIBER_PREFORK=false \
    FIBER_STRICT_ROUTING=true \
    FIBER_CASE_SENSITIVE=true \
    FIBER_BODY_LIMIT=0 \
    FIBER_READ_TIMEOUT=0 \
    FIBER_WRITE_TIMEOUT=0 \
    FIBER_REDUCE_MEMU=true \
    FIBER_JSON=json \
    HASH_SALT=12 \
    LOG_PATH=resource/logs \
    LOG_MAX_SIZE=10 \
    LOG_MAX_BACKUP=5 \
    LOG_MAX_SIZE_ROTATE=20 \
    LOG_TIME_FORMAT="2006-01-02" \
    LOG_COLOR_OUTPUT=false \
    LOG_QUOTE_STR=false \
    LOG_END_WITH_MESSAGE=false \
    DB_DRIVER="postgres" \
    DB_PROTOCOL="postgresql" \
    DB_NAME="" \
    DB_HOST="" \
    DB_PORT=5432 \
    DB_USER="" \
    DB_PASS="" \
    DB_MAX_CON=100 \
    DB_MIN_CON=10 \
    DB_MAX_TIME=30 \
    DB_MIN_TIME=5 \
    TOKEN_NAME="RESTful_API_AUTH" \
    SECRET_KEY_ACCESS_TOKEN="" \
    SECRET_KEY_REFRESH_TOKEN="" \
    SECRET_KEY_FP_TOKEN="" \
    SECRET_KEY_CSRF="" \
    SECRET_TEST_CLIENT="" \
    CACHE_TIMEOUT=8 \
    DB_TIMEOUT=20 \
    DOWN_STREAM_TIMEOUT=30 \
    CORS_ALLOW_METHODS="" \
    CORS_ALLOW_HEADERS="" \
    CORS_EXPOSE_HEADERS="" \
    CORS_ALLOW_ORIGINS="*" \
    CORS_ALLOW_CREDENTIALS=""

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

WORKDIR /app

USER root

# Copy application files with correct ownership
COPY --chown=${username}:${username} --from=builder /app/build/main .
COPY --chown=${username}:${username} --from=builder /app/docs /app/docs
COPY --chown=${username}:${username} --from=builder /app/resource /app/resource

# Expose the necessary port
EXPOSE ${PORT}

# Use an unprivileged user.
USER ${username}:${username}

# Command to run the application
ENTRYPOINT ["/app/main"]
