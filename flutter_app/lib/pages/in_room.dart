import 'dart:convert';

import 'package:firebase_auth/firebase_auth.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_app/controllers/loading_dialog.dart';
import 'package:flutter_app/controllers/shared_preferences.dart';
import 'package:flutter_app/pages/room_page.dart';
import 'package:flutter_app/pages/setting_page.dart';
import 'package:http/http.dart' as http;

import '../home_page.dart';


class InRoomArguments {
  final Room roomInfo;
  InRoomArguments(this.roomInfo);
}

class InRoom extends StatefulWidget {
  static const routeName = '/in_room';

  @override
  _InRoomState createState() => _InRoomState();
}

class _InRoomState extends State<InRoom> {
  SharedPrefs _prefs;
  String _roomId;
  String _roomName = '';
  Room _roomInfo;
  String _enteredTime = '　時　分';
  List<UserStatus> _roomUsers = [];

  bool _isButtonDisabled = true;

  Future _init() async {
    _prefs = await generateSharedPrefs();
    _roomId = await _prefs.getCurrentRoomId();
    _roomName = await _prefs.getCurrentRoomName();
    await updateUserData();
    setState(() {});
    await updateRoomInfo();
  }

  Future updateRoomInfo() async {
    Map<String, String> queryParams = {
      'room_id': _roomId
    };
    Uri uri = Uri.https('us-central1-online-study-room-f1f30.cloudfunctions.net', '/RoomStatus', queryParams);
    final response = await http.get(uri);
    if (response.statusCode == 200) {
      RoomStatusResponse roomStatusResponse = RoomStatusResponse.fromJson(json.decode(utf8.decode(response.bodyBytes)));
      List<UserStatus> users = roomStatusResponse.users;
      print(users);

      setState(() {
        _roomUsers = users;
      });
    }
  }

  Future updateUserData() async {
    Map<String, String> queryParams = {
      'user_id': await _prefs.getUserId()
    };
    Uri uri = Uri.https('us-central1-online-study-room-f1f30.cloudfunctions.net', '/UserStatus', queryParams);
    final response = await http.get(uri);
    if (response.statusCode == 200) {
      UserStatusResponse userStatusResp = UserStatusResponse.fromJson(
          json.decode(utf8.decode(response.bodyBytes)));
      if (userStatusResp.result == 'ok') {
        UserBody user = userStatusResp.userStatus.userBody;
        _enteredTime = user.lastEntered.hour.toString() + '時' + user.lastEntered.minute.toString() + '分';
      }
    }
  }

  Future<void> exitRoom(BuildContext context, Room roomInfo) async {
    _prefs = await generateSharedPrefs();

    final _body = {
      'room_id': roomInfo.roomId,
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
        Navigator.popUntil(context, ModalRoute.withName(MyHomePage.routeName));
      } else {
        Navigator.pop(context);
        // todo show alert
        throw Exception('Failed to enter room : ' + exitRoomResp.message);
      }
    } else {
      // todo show alert
      Navigator.pop(context);
      throw Exception('http request failed');
    }
  }

  @override
  void initState() {
    super.initState();
    _init();
  }

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    final InRoomArguments arg = ModalRoute.of(context).settings.arguments;
    _roomInfo = arg.roomInfo;
    setState(() {
      _isButtonDisabled = false;
    });
  }

  void showExitRoomDialog(BuildContext context, Room roomInfo) {
    showDialog(
      context: context,
      builder: (context) {
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
                LoadingDialog.show(context, title: '退室中');
                exitRoom(context, roomInfo);
              },
            ),
          ],
        );
      }
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(_roomName + 'の部屋'),
        leading: FlatButton(
          onPressed: () => _isButtonDisabled ? null : showExitRoomDialog(context, _roomInfo),
          child: Icon(
            Icons.close,
            color: _isButtonDisabled ? Colors.black : Colors.white,
          ),
        ),
      ),
      body: Column(
        children: [
          Text('入室時刻：' + _enteredTime),
          Divider(),
          Align(
              child: Text('入室中の他のユーザー'),
              alignment: Alignment.centerLeft
          ),
          Flexible(
            child: GridView.builder(
              itemCount: _roomUsers.length,
              gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
                crossAxisCount: 2
              ),
              itemBuilder: (BuildContext context, int index) {
                if (_roomUsers[index].userId != FirebaseAuth.instance.currentUser.uid) {
                  return GridTile(
                    child: Text(_roomUsers[index].displayName),
                  );
                } else {
                  return null;
                }
              },
            ),
          )
        ]
      )
    );
  }
}

class RoomStatusResponse {
  final String result;
  final String message;
  final RoomStatus roomStatus;
  final List<UserStatus> users;

  RoomStatusResponse({
    this.result,
    this.message,
    this.roomStatus,
    this.users
  });

  factory RoomStatusResponse.fromJson(Map<String, dynamic> json) {
    return RoomStatusResponse(
        result: json['result'] as String,
        message: json['message'] as String,
        roomStatus: RoomStatus.fromJson(json['room_status']),
        users: (json['users'] as List<dynamic>).map((i) => UserStatus.fromJson(i)).toList()
    );
  }
}

class RoomStatus {
  final String roomId;
  final RoomBody roomBody;

  RoomStatus({this.roomId, this.roomBody});

  factory RoomStatus.fromJson(Map<String, dynamic> json) {
    return RoomStatus(
        roomId: json['room_id'] as String,
        roomBody: RoomBody.fromJson(json['room_body'])
    );
  }
}

class RoomBody {
  final DateTime created;
  final String name;
  final List<dynamic> users;

  RoomBody({this.created, this.name, this.users});

  factory RoomBody.fromJson(Map<String, dynamic> json) {
    return RoomBody(
        created: DateTime.parse(json['created']).toLocal(),
        name: json['name'] as String,
        users: json['users'] as List<dynamic>,
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
