import 'dart:convert';

import 'package:firebase_auth/firebase_auth.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_app/controllers/shared_preferences.dart';
import 'package:http/http.dart' as http;


class InRoom extends StatefulWidget {
  @override
  _InRoomState createState() => _InRoomState();
}

class _InRoomState extends State<InRoom> {
  SharedPrefs _prefs;
  String _roomId;
  String _roomName = '';

  Future _init() async {
    _prefs = await generateSharedPrefs();
    _roomId = await _prefs.getCurrentRoomId();
    _roomName = await _prefs.getCurrentRoomName();
  }
  @override
  void initState() {
    _init();
    super.initState();
  }

  void showExitRoomDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (_) {
        return AlertDialog(
          title: Text("部屋を退出しますか？"),
          content: null,
          actions: <Widget>[
            // ボタン領域
            FlatButton(
              child: Text("キャンセル"),
              onPressed: () => Navigator.pop(context),
            ),
            FlatButton(
              child: Text("OK"),
              onPressed: () async {
                await exitRoom();
              },
            ),
          ],
        );
      }
    );
  }

  Future<void> exitRoom() async {
    final _body = {
      'room_id': _roomId,
      'user_id': await _prefs.getUserId(),
      'id_token': await FirebaseAuth.instance.currentUser.getIdToken(),
    };
    Uri uri = Uri.https('us-central1-online-study-room-f1f30.cloudfunctions.net', '/ExitRoom');

    final response = await http.post(
        uri,
        body: _body);
    if (response.statusCode == 200) {
      ExitRoomResponse exitRoomResp = ExitRoomResponse.fromJson(json.decode(utf8.decode(response.bodyBytes)));
      if (exitRoomResp.result == 'ok') {
        Navigator.popUntil(context, ModalRoute.withName('/home'));
      } else {
        Navigator.pop(context);
        // todo show alert
        throw Exception('Failed to enter room : ' + exitRoomResp.message);
      }
    } else {
      throw Exception('http request failed');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(_roomName + 'の部屋'),
        leading: FlatButton(
          onPressed: () {
            showExitRoomDialog(context);
          },
          child: Icon(
            Icons.close
          ),
        ),
      ),
    );
  }
}

class ExitRoomResponse {
  final String result;
  final String message;

  ExitRoomResponse({this.result, this.message});

  factory ExitRoomResponse.fromJson(Map<String, dynamic> json) {
    return ExitRoomResponse(
      result: json['result'] as String,
      message: json['message'] as String,
    );
  }
}
