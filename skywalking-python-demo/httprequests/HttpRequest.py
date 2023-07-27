import requests


def send_post_request(url, headers, data):
    response = requests.post(url=url, headers=headers, data=data)
    return response