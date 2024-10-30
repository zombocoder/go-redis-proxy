from locust import HttpUser, task, between
import redis
import random
import string


# Define Redis command tasks
class RedisTaskSet(HttpUser):
    wait_time = between(1, 3)

    def on_start(self):
        # Connect to Redis on startup
        self.client = redis.Redis(
            host="localhost", port=6380
        )

    @task
    def set_command(self):
        key = "".join(random.choices(string.ascii_lowercase, k=10))
        value = "".join(random.choices(string.ascii_lowercase, k=20))
        self.client.set(key, value)

    @task
    def get_command(self):
        key = "sample_key"
        self.client.get(key)

    @task
    def incr_command(self):
        key = "counter"
        self.client.incr(key)

    @task
    def lpush_command(self):
        key = "mylist"
        value = "".join(random.choices(string.ascii_lowercase, k=5))
        self.client.lpush(key, value)

    @task
    def lrange_command(self):
        key = "mylist"
        result = self.client.lrange(key, 0, -1)
        [item.decode() for item in result]  # Decode each item in the list

    @task
    def sadd_command(self):
        key = "myset"
        value = "".join(random.choices(string.ascii_lowercase, k=5))
        self.client.sadd(key, value)

    @task
    def smembers_command(self):
        key = "myset"
        result = self.client.smembers(key)
        {item.decode() for item in result}  # Decode each item in the set

    @task
    def hset_command(self):
        key = "myhash"
        field = "".join(random.choices(string.ascii_lowercase, k=5))
        value = "".join(random.choices(string.ascii_lowercase, k=10))
        self.client.hset(key, field, value)

    @task
    def hgetall_command(self):
        key = "myhash"
        result = self.client.hgetall(key)
        {
            k.decode(): v.decode() for k, v in result.items()
        }

    @task
    def zadd_command(self):
        key = "myzset"
        value = "".join(random.choices(string.ascii_lowercase, k=5))
        score = random.randint(1, 100)
        self.client.zadd(key, {value: score})

    @task
    def zrange_command(self):
        key = "myzset"
        result = self.client.zrange(key, 0, -1)
        [item.decode() for item in result]  # Decode each item in the sorted set

    @task
    def ttl_command(self):
        key = "sample_key"
        self.client.ttl(key)  # TTL returns an integer, no decoding needed

    @task
    def del_command(self):
        key = "sample_key"
        self.client.delete(key)
