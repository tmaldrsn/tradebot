import os
from redis.asyncio import Redis


def get_redis_connection():
    redis_host = os.getenv("REDIS_HOST")
    if not redis_host:
        raise Exception("Environment variable `REDIS_HOST` is not set.")
    return Redis(host=redis_host, port=6379, decode_responses=True)