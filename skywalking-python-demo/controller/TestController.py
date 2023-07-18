import argparse

from flask import Flask, request, jsonify
import requests
from flask_sqlalchemy import SQLAlchemy
from skywalking import config, agent


config.init(
            agent_collector_backend_services='',
            agent_protocol='grpc',
            agent_authentication='',
            agent_name='',
            agent_meter_reporter_active=False,
            agent_log_reporter_active=False)


app = Flask(__name__)


# 配置数据库的连接信息，需要替换成自己的数据库连接信息
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


agent.start()


@app.route('/user', methods=['GET'])
def get_all_users():
    users = User.query.all()
    return jsonify([user.to_dict() for user in users])


@app.route('/user', methods=['POST'])
def add_user():
    data = request.get_json()
    username = data.get('username')
    password = data.get('password')
    print('erve brb')
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


@app.route('/test', methods=['GET'])
def test():
    return 'Hello, World!'


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument("--host", type=str, default="127.0.0.1")
    parser.add_argument("--port", type=int, default=10000)

    args = parser.parse_args()

    app.run(debug=False, host=args.host, port=args.port)

