import json
from flask import jsonify, Flask
import requests

app = Flask(__name__)


def test1():
    url = 'http://127.0.0.1:9999/forward/user'
    # import urllib3.contrib.pyopenssl
    username = 'LiuHan'
    password = 'lh232r425'
    headers = {'Content-Type': 'application/json'}
    data = json.dumps({'username': username, 'password': password})
    # urllib3.contrib.pyopenssl.inject_into_urllib3()
    response = requests.post(url, headers=headers, data=data)
    print(json.loads(response.text))


def test2():
    url = 'http://127.0.0.1:9999/forward/user'
    response = requests.get(url)
    # with app.app_context():
    print(json.loads(response.text))


def test3():
    url = 'http://127.0.0.1:9999/forward/user/1'
    headers = {'Content-Type': 'application/json'}
    data = json.dumps({'password': "urvbyertbtrbn"})

    response = requests.put(url=url, headers=headers, data=data)

    print(json.loads(response.text))


if __name__ == "__main__":
    test2()
    test3()



