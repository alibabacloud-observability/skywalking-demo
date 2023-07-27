
import json
import requests
from skywalking import agent, config
from flask import Flask, request, jsonify
from flask_sqlalchemy import SQLAlchemy
import argparse

from httprequests.HttpRequest import send_post_request

config.init(
            agent_collector_backend_services='',
            agent_protocol='grpc',
            agent_authentication='',
            agent_name='',
            agent_meter_reporter_active=False,
            agent_log_reporter_active=False)

agent.start()

app = Flask(__name__)

# 配置数据库的连接信息
app.config['SQLALCHEMY_DATABASE_URI'] = 'mysql://root@127.0.0.1:3306/skywalking_test'
# 关闭动态追踪修改的警告信息
app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False
# 展示sql语句
app.config['SQLALCHEMY_ECHO'] = True
db = SQLAlchemy(app)


class User(db.Model):
    id = db.Column(db.Integer, primary_key=True, autoincrement=True)
    username = db.Column(db.String(50), nullable=False)
    password = db.Column(db.String(50), nullable=False)

    def to_dict(self):
        return dict(id=self.id, name=self.username, password=self.password)


# 定义路由
@app.route('/', methods=['GET'])
def index():
    return "hello"


@app.route('/test', methods=['GET'])
def test():
    url = 'http://127.0.0.1:10000/test'
    response = requests.get(url)

    return jsonify(response.text)
    # 处理响应


@app.route('/api', methods=['POST'])
def api():
    data = request.json
    response = {'message': 'Hello, {}!'.format(data.get('username'))}
    return jsonify(response)


@app.route('/user', methods=['GET'])
def get_all_users():
    users = User.query.all()
    return jsonify([user.to_dict() for user in users])


@app.route('/user/<int:user_id>', methods=['GET'])
def get_user(user_id):
    user = User.query.filter_by(id=user_id).first()
    if user:
        return jsonify(user.to_dict())
    else:
        return jsonify(error='User not found'), 404


@app.route('/user', methods=['POST'])
def add_user():
    data = request.get_json()
    username = data.get('username')
    password = data.get('password')
    if username and password:
        user = User(username=username, password=password)
        db.session.add(user)
        db.session.commit()
        return jsonify(user.to_dict())
    else:
        return jsonify(error='Bad Request'), 400


@app.route('/user/<int:user_id>', methods=['PUT'])
def update_user(user_id):
    user = User.query.filter_by(id=user_id).first()
    if user:
        data = request.get_json()
        name = data.get('username')
        password = data.get('password')
        if name:
            user.username = name
        if password:
            user.password = password
        db.session.commit()
        return jsonify(user.to_dict())
    else:
        return jsonify(error='User not found'), 404


@app.route('/user/<int:user_id>', methods=['DELETE'])
def delete_user(user_id):
    user = User.query.filter_by(id=user_id).first()
    if user:
        db.session.delete(user)
        db.session.commit()
        return '', 204
    else:
        return jsonify(error='User not found'), 404


@app.route('/forward/user', methods=['GET', 'POST'])
def forward():
    response = None
    print("wvwvw")
    url = 'http://127.0.0.1:10000/user'
    if request.method == 'GET':
        response = requests.get(url=url)
    elif request.method == 'POST':
        data = request.get_json()
        username = data.get('username')
        password = data.get('password')

        headers = {'Content-Type': 'application/json'}
        newdata = {'username': username, 'password': password}
        json_data = json.dumps(newdata)

        # response = httprequests.post(url=url, headers=headers, data=json_data)
        response = send_post_request(url=url, headers=headers, data=json_data)

    if response.status_code == 200:
        return jsonify(response.text)
    else:
        return '请求失败'


@app.route('/forward/user/<int:user_id>', methods=['PUT'])
def forward_update_user(user_id):
    response = None
    if request.method == 'PUT':
        url = 'http://127.0.0.1:10000/user' + '/' + str(user_id)
        headers = {'Content-Type': 'application/json'}
        response = requests.put(url=url, headers=headers, data=request.data)

    if response.status_code == 200:
        return jsonify(response.text)
    else:
        return '请求失败'


@app.route('/forward/user/<int:user_id>', methods=['DELETE'])
def forward_delete_user(user_id):
    response = None
    if request.method == 'PUT':
        url = 'http://127.0.0.1:10000/user' + '/' + str(user_id)
        headers = {'Content-Type': 'application/json'}
        response = requests.put(url=url, headers=headers, data=request.data)

    if response.status_code == 200:
        return jsonify(response.text)
    else:
        return '请求失败'


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument("--host", type=str, default="127.0.0.1")
    parser.add_argument("--port", type=int, default=9999)

    args = parser.parse_args()

    app.run(debug=False, host=args.host, port=args.port)
