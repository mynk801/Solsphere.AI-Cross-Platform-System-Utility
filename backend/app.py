# app.py

import os
import json
from flask import Flask, request, jsonify
from flask_sqlalchemy import SQLAlchemy
from datetime import datetime
from flask_cors import CORS 

from config import SECRET_API_KEY


app = Flask(__name__)
CORS(app)

app.config.from_mapping(
    SECRET_KEY='dev',
    SQLALCHEMY_DATABASE_URI='sqlite:///' + os.path.join(app.instance_path, 'devices.db'),
    SQLALCHEMY_TRACK_MODIFICATIONS=False
)

db = SQLAlchemy(app)


class Device(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    machine_id = db.Column(db.String(100), unique=True, nullable=False)
    os = db.Column(db.String(50), nullable=False)
    last_check_in = db.Column(db.DateTime, default=datetime.utcnow, onupdate=datetime.utcnow)
    latest_data = db.Column(db.Text)

    def to_dict(self):
        data = self.latest_data
        if isinstance(data, str) and data:
            try:
                data = json.loads(data)
            except json.JSONDecodeError:
                data = {"error": "invalid json data"}
        return {
            'machineId': self.machine_id,
            'os': self.os,
            'lastCheckIn': self.last_check_in.isoformat() + 'Z',
            'latestData': data
        }


@app.route('/api/report', methods=['POST'])
def report_data():
    api_key = request.headers.get('X-API-Key')
    if not api_key or api_key != SECRET_API_KEY:
        return jsonify({"error": "Unauthorized"}), 401
    data = request.json
    if not data or 'machineId' not in data:
        return jsonify({"error": "Invalid data payload"}), 400
    machine_id = data.get('machineId')
    check_results = data.get('checkResults')
    latest_data_str = json.dumps(check_results)
    device = Device.query.filter_by(machine_id=machine_id).first()
    if device:
        device.os = data.get('os', device.os)
        device.latest_data = latest_data_str
    else:
        device = Device(
            machine_id=machine_id,
            os=data.get('os'),
            latest_data=latest_data_str
        )
        db.session.add(device)
    db.session.commit()
    return jsonify({"message": "Data received successfully"}), 200

@app.route('/api/devices', methods=['GET'])
def get_devices():
    all_devices = Device.query.all()
    device_list = [device.to_dict() for device in all_devices]
    return jsonify(device_list)


@app.cli.command("init-db")
def init_db_command():
    try:
        os.makedirs(app.instance_path)
    except OSError:
        pass
    db.create_all()
    print("Initialized the database.")