from flask import jsonify, Flask
from flask_sqlalchemy import SQLAlchemy
from sqlalchemy import text
from sqlalchemy.testing import db
from sqlalchemy.sql import quoted_name

app = Flask(__name__)
app.config['SQLALCHEMY_DATABASE_URI'] = 'mysql://root@127.0.0.1:3306/skywalking_test'  # 根据自己数据库实际情况进行修改
app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False
app.config['SQLALCHEMY_ECHO'] = True
db = SQLAlchemy(app)


@app.route('/test', methods=['GET'])
def test_connection():
    result = db.session.execute(text("select * from " + quoted_name("user", True)))
    # result = db.session.execute(text("show tables"))
    # return result
    # result = db.session.execute(text("show tables"))
    tables = [row[1] for row in result.fetchall()]
    return jsonify(tables)


if __name__ == '__main__':
    app.run(debug=False, host='127.0.0.1', port=9000)
