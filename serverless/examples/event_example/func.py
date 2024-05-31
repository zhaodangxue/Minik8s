import redis
from datetime import datetime

def main(params):
    # 更新redis
    redis_url = "redis://127.0.0.1:6379"
    func_log_key = "event-example-status"

    r = redis.StrictRedis.from_url(redis_url)
    res = r.get(func_log_key)
    if res:
        res = eval(res)
    else:
        res = {"invoke": 0, "last_invoke": ""}
    
    res["invoke"] += 1
    res["last_invoke"] = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    res = str(res)
    r.set(func_log_key, res)

if __name__ == "__main__":
    main({})
