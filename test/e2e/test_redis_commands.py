import redis
import pytest
import random
import string


# Connect to Redis
@pytest.fixture(scope="module")
def redis_client():
    client = redis.Redis(host="localhost", port=6380)
    yield client
    client.flushdb()  # Clean up after tests


# Helper function to generate random strings
def random_string(length=8):
    return "".join(random.choices(string.ascii_lowercase, k=length))


def test_set_and_get(redis_client):
    key, value = random_string(), random_string()
    redis_client.set(key, value)
    assert redis_client.get(key).decode() == value


def test_incr(redis_client):
    key = "counter"
    redis_client.set(key, 0)
    redis_client.incr(key)
    assert int(redis_client.get(key)) == 1


def test_lpush_and_lrange(redis_client):
    key = "mylist"
    values = [random_string() for _ in range(3)]
    for value in values:
        redis_client.lpush(key, value)
    list_values = [item.decode() for item in redis_client.lrange(key, 0, -1)]
    assert list_values == values[::-1]  # LPUSH adds to start, so we reverse


def test_sadd_and_smembers(redis_client):
    key = "myset"
    values = {random_string() for _ in range(3)}
    for value in values:
        redis_client.sadd(key, value)
    assert set(redis_client.smembers(key)) == {v.encode() for v in values}


def test_hset_and_hgetall(redis_client):
    key = "myhash"
    field_value_pairs = {f"field_{i}": random_string() for i in range(3)}
    for field, value in field_value_pairs.items():
        redis_client.hset(key, field, value)
    hgetall_result = {
        k.decode(): v.decode() for k, v in redis_client.hgetall(key).items()
    }
    assert hgetall_result == field_value_pairs


def test_zadd_and_zrange(redis_client):
    key = "myzset"
    values = {random_string(): random.randint(1, 100) for _ in range(3)}
    redis_client.zadd(key, values)
    zset_values = [item.decode() for item in redis_client.zrange(key, 0, -1)]
    assert set(zset_values) == set(values.keys())


def test_ttl(redis_client):
    key, value = random_string(), random_string()
    redis_client.set(key, value, ex=10)  # Set with expiration of 10 seconds
    ttl = redis_client.ttl(key)
    assert ttl > 0 and ttl <= 10


def test_del(redis_client):
    key, value = random_string(), random_string()
    redis_client.set(key, value)
    assert redis_client.get(key).decode() == value
    redis_client.delete(key)
    assert redis_client.get(key) is None
