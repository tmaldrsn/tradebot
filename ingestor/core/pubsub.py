import json


async def publish_event(rdb, channel, event):
    await rdb.publish(channel, json.dumps(event))