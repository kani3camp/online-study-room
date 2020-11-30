import 'dart:async';
import 'dart:convert';

import 'package:firebase_auth/firebase_auth.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_app/controllers/api_links.dart';
import 'package:flutter_app/controllers/custom_dialog.dart';
import 'package:flutter_app/controllers/shared_preferences.dart';
import 'package:flutter_app/pages/room_page.dart';
import 'package:flutter_app/pages/setting_page.dart';
import 'package:http/http.dart' as http;
import 'package:web_socket_channel/io.dart';

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
  String _userId;
  String _roomId;
  String _roomName = '';
  Room _roomInfo;
  DateTime _enteredTime = DateTime.now();
  List<UserStatus> _roomOtherUsers = [];

  IOWebSocketChannel _channel;

  bool _isButtonDisabled = true;

  Future _init() async {
    _prefs = await generateSharedPrefs();
    _userId = FirebaseAuth.instance.currentUser.uid;
    _roomId = await _prefs.getCurrentRoomId();
    _roomName = await _prefs.getCurrentRoomName();
    await updateMyUserData();
    setState(() {});
    await updateRoomInfo();

    startStudying();
    stayStudying();
  }

  Future updateRoomInfo() async {
    Map<String, String> queryParams = {'room_id': _roomId};
    Uri uri = Uri.https(ApiLinks.Authority, ApiLinks.RoomStatus, queryParams);
    final response = await http.get(uri);
    if (response.statusCode == 200) {
      RoomStatusResponse roomStatus = RoomStatusResponse.fromJson(json.decode(utf8.decode(response.bodyBytes)));
      List<UserStatus> users = roomStatus.users;

      if (this.mounted) {
        setState(() {
          users.removeWhere((element) => element.userId == FirebaseAuth.instance.currentUser.uid);
          _roomOtherUsers = users;
        });
      }
    }
  }

  Future updateMyUserData() async {
    Map<String, String> queryParams = {'user_id': FirebaseAuth.instance.currentUser.uid};
    Uri uri = Uri.https(ApiLinks.Authority, ApiLinks.UserStatus, queryParams);
    final response = await http.get(uri);
    if (response.statusCode == 200) {
      UserStatusResponse userStatusResp = UserStatusResponse.fromJson(
          json.decode(utf8.decode(response.bodyBytes)));
      if (userStatusResp.result == 'ok') {
        UserBody user = userStatusResp.userStatus.userBody;
        _enteredTime = user.lastEntered;
      }
    }
  }

  void startStudying() async {
    _channel = IOWebSocketChannel.connect(
        'wss://0ieer51ju9.execute-api.ap-northeast-1.amazonaws.com/production');
    _channel.stream.listen((event) {
      if (this.mounted) {
        StayStudyingResponse resp = StayStudyingResponse.fromJson(json.decode(event.toString()));
        if (resp.isOk) {
          print('stay studying : OK (' + DateTime.now().toString() + ')');
          List<UserStatus> users = resp.users;

          for (UserStatus user in users) {
            if (user.userId == FirebaseAuth.instance.currentUser.uid) {
              _enteredTime = user.userBody.lastEntered;
            }
          }
          users.removeWhere((element) => element.userId == FirebaseAuth.instance.currentUser.uid);
          this.setState(() {
            _roomOtherUsers = users;
          });
        } else {
          print('stay studying : NG');
          print('message: ' + resp.message);
          CustomDialog.showAlertDialog(context,
            '問題が発生しました。ルームを出ます。\n' + resp.message,
            onOkPressed: () {
              Navigator.popUntil(context, ModalRoute.withName(MyHomePage.routeName));
            }
          );
        }
      }
    });
  }

  void stayStudying() async {
    if (this.mounted) {
      final _params = json.encode({
        'user_id': _userId,
        'id_token': await FirebaseAuth.instance.currentUser.getIdToken(),
        'room_id': _roomId,
        'device_type': 'mobile'
      });
      _channel?.sink?.add(_params);
    }
    Timer(Duration(seconds: 5), stayStudying);
  }

  Future<void> exitRoom(BuildContext context, Room roomInfo) async {
    _prefs = await generateSharedPrefs();

    final _body = json.encode({
      'room_id': roomInfo.roomId,
      'user_id': FirebaseAuth.instance.currentUser.uid,
      'id_token': await FirebaseAuth.instance.currentUser.getIdToken(),
    });
    Uri uri = Uri.https(ApiLinks.Authority, ApiLinks.ExitRoom);

    final response = await http.post(uri, body: _body);
    if (response.statusCode == 200) {
      ExitRoomResponse exitRoomResp = ExitRoomResponse.fromJson(
          json.decode(utf8.decode(response.bodyBytes)));
      if (exitRoomResp.result == 'ok') {
        Navigator.popUntil(context, ModalRoute.withName(MyHomePage.routeName));
      } else {
        CustomDialog.showAlertDialog(context,
          '問題が発生しました。\n' + exitRoomResp.message,
          onOkPressed: () {
            Navigator.popUntil(context, ModalRoute.withName(MyHomePage.routeName));
          }
        );
      }
    } else {
      CustomDialog.showAlertDialog(context,
        '通信が失敗しました。\n',
        onOkPressed: () {
          Navigator.popUntil(context, ModalRoute.withName(MyHomePage.routeName));
        }
      );
    }
  }

  void showExitRoomDialog(BuildContext context, Room roomInfo) {
    showDialog(
        context: context,
        builder: (context) {
          return AlertDialog(
            title: Text("部屋を出ますか？"),
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
                  CustomDialog.showLoadingDialog(context, title: '退室中');
                  exitRoom(context, roomInfo);
                },
              ),
            ],
          );
        });
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

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(
          backgroundColor: Colors.white,
          title: Text(
            _roomName + 'の部屋',
            style: TextStyle(
                color: Theme.of(context).primaryColor
            ),
          ),
          leading: FlatButton(
            onPressed: () => _isButtonDisabled
                ? null
                : showExitRoomDialog(context, _roomInfo),
            child: Icon(
              Icons.close,
              color: _isButtonDisabled ? Theme.of(context).backgroundColor : Theme.of(context).primaryColor,
            ),
          ),
        ),
        body: Column(children: [
          Padding(
            padding: EdgeInsets.all(5.0),
            child: Opacity(
              opacity: 0.5,
              child: Text('入室時刻：'
                  + _enteredTime.hour.toString()
                  + '時'
                  + _enteredTime.minute.toString()
                  + '分',
              ),
            ),
          ),
          Divider(),
          Align(
            child: Padding(
              padding: EdgeInsets.all(10.0),
              child: Text('同じ部屋の他のユーザー'),
            ), 
            alignment: Alignment.centerLeft,
          ),
          Flexible(
            child: GridView.builder(
              itemCount: _roomOtherUsers.length,
              gridDelegate:
                  SliverGridDelegateWithFixedCrossAxisCount(crossAxisCount: 2),
              itemBuilder: (BuildContext context, int index) {
                // 滞在時間を計算
                final int nowSecondsSinceEpoch = ((DateTime.now().millisecondsSinceEpoch) / 1000.0).floor();
                final int lastEnteredSecondsSinceEpoch = ((_roomOtherUsers[index].userBody.lastEntered.millisecondsSinceEpoch) / 1000.0).floor();
                final int timeStudySeconds = nowSecondsSinceEpoch - lastEnteredSecondsSinceEpoch;
                final int timeStudyMinutes = (timeStudySeconds / 60.0).floor();
                print('timeStudySeconds: ' + timeStudySeconds.toString());
                print('timeStudyMinutes: ' + timeStudyMinutes.toString());

                return Container(
                  padding: EdgeInsets.all(10.0),
                  child: Container(
                    padding: EdgeInsets.all(10.0),
                    decoration: BoxDecoration(
                      border: Border.all(
                        width: 1,
                      ),
                    ),
                    child: GridTile(
                      header: Padding(
                        padding: EdgeInsets.all(10.0),
                        child: Icon(
                          Icons.account_circle,
                          size: 50,
                        ),
                      ),
                      child: Center(
                        child: Text(
                          _roomOtherUsers[index].displayName.length < 10
                              ? _roomOtherUsers[index].displayName
                              : _roomOtherUsers[index].displayName.substring(0, 9) + '…',
                          style: TextStyle(
                            fontSize: 20,
                          ),
                        ),
                      ),
                      footer: Center(child: Text(timeStudyMinutes.toString() + '分')),
                    ),
                  ),
                );
              },
            ),
          )
        ]));
  }

  @override
  void dispose() {
    _channel?.sink?.close();
    super.dispose();
  }
}

class RoomStatusResponse {
  final String result;
  final String message;
  final RoomStatus roomStatus;
  final List<UserStatus> users;

  RoomStatusResponse({this.result, this.message, this.roomStatus, this.users});

  factory RoomStatusResponse.fromJson(Map<String, dynamic> json) {
    return RoomStatusResponse(
      result: json['result'] as String,
      message: json['message'] as String,
      roomStatus: RoomStatus.fromJson(json['room_status']),
      users: (json['users'] as List).length > 0 ? (json['users'] as List<dynamic>)
          .map((i) => UserStatus.fromJson(i))
          .toList() : [],
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
        roomBody: RoomBody.fromJson(json['room_body']));
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

class StayStudyingResponse {
  final bool isOk;
  final String message;
  final List<dynamic> users;

  StayStudyingResponse({this.isOk, this.message, this.users});

  factory StayStudyingResponse.fromJson(Map<String, dynamic> json) {
    return StayStudyingResponse(
      isOk: json['is_ok'] as bool,
      message: json['message'] as String,
      users: json['users'] is List
          ? (json['users'] as List<dynamic>).map((i) => UserStatus.fromJson(i)).toList()
          : [],
    );
  }
}
