import 'dart:convert';

import 'package:firebase_auth/firebase_auth.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_app/controllers/loading_dialog.dart';
import 'package:flutter_app/controllers/shared_preferences.dart';
import 'package:http/http.dart' as http;

import 'in_room.dart';

class RoomPage extends StatefulWidget {
  @override
  _RoomPageState createState() => _RoomPageState();
}

class _RoomPageState extends State<RoomPage> {
  Future<List<Room>> _futureList;
  SharedPrefs _prefs;

  Future<void> _init() async {
    _futureList = fetchRooms();
    _prefs = await generateSharedPrefs();
  }

  Future<bool> _onRefresh() async {
    _futureList = fetchRooms();
    await _futureList;
    setState(() {});
    return true;
  }

  Future<void> enterRoom(BuildContext context, Room roomInfo) async {
    _prefs = await generateSharedPrefs();

    final _body = {
      'room_id': roomInfo.roomId,
      'user_id': await _prefs.getUserId(),
      'id_token': await FirebaseAuth.instance.currentUser.getIdToken(),
    };
    Uri uri = Uri.https('us-central1-online-study-room-f1f30.cloudfunctions.net', '/EnterRoom');

    final response = await http.post(
        uri,
        body: _body);
    if (response.statusCode == 200) {
      EnterRoomResponse enterRoomResp = EnterRoomResponse.fromJson(json.decode(utf8.decode(response.bodyBytes)));
      if (enterRoomResp.result == 'ok') {
        await _prefs.setCurrentRoomId(roomInfo.roomId);
        await _prefs.setCurrentRoomName(roomInfo.roomBody.name);
        Navigator.of(context).pushNamed(
            InRoom.routeName,
            arguments: InRoomArguments(roomInfo)
        );
      } else {
        Navigator.pop(context);
        throw Exception('Failed to enter room : ' + enterRoomResp.message);
      }
    } else {
      Navigator.pop(context);
      throw Exception('http request failed');
    }
  }

  @override
  void initState() {
    _init();
    super.initState();
  }

  void showEnterRoomDialog(BuildContext context, Room roomInfo) {
    showDialog(
        context: context,
        builder: (context) {
          return AlertDialog(
            title: Text(roomInfo.roomBody.name + 'の部屋に入りますか？'),
            content: null,
            actions: <Widget>[
              // ボタン領域
              FlatButton(
                child: Text("キャンセル"),
                onPressed: () => Navigator.pop(context),
              ),
              FlatButton(
                child: Text("OK"),
                onPressed: () {
                  LoadingDialog.show(context, title: '入室中');
                  enterRoom(context, roomInfo);
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
        title: Center(
            child: Text('カテゴリ一覧')
        ),
      ),
      body: RefreshIndicator(
          onRefresh: () {
            return _onRefresh();
          },
          child: FutureBuilder<List<Room>>(
            future: _futureList,
            builder: (context, snapshot) {
              if (snapshot.hasData) {
                final List<Room> rooms = snapshot.data;
                return ListView.separated(
                  padding: const EdgeInsets.all(8),
                  itemCount: rooms.length,
                  itemBuilder: (BuildContext context, int index) {
                    return Container(
                      child: ListTile(
                        title: Text(rooms[index].roomBody.name),
                        onTap: () {
                          showEnterRoomDialog(context, rooms[index]);
                        },
                      ),
                    );
                  },
                  separatorBuilder: (BuildContext context, int index) => const Divider(),
                );
              } else if (snapshot.hasError) {
                print(snapshot.error);
              } else {
              }
              return Center(child: CircularProgressIndicator());
            },
          )
      ),
    );
  }
}

Future<List<Room>> fetchRooms() async {
  print('fetchRooms()');
  const url = 'https://us-central1-online-study-room-f1f30.cloudfunctions.net/Rooms';
  final response = await http.get(url);
  if (response.statusCode == 200) {
    RoomsResponse roomsResponse = RoomsResponse.fromJson(json.decode(utf8.decode(response.bodyBytes)));
    if (roomsResponse.result == 'ok') {
      return roomsResponse.rooms;
    } else {
      throw Exception('Failed to load room list: ' + roomsResponse.message);
    }
  } else {
    throw Exception('http request failed');
  }
}

class RoomsResponse {
  final String result;
  final String message;
  final List<Room> rooms;

  RoomsResponse({this.result, this.message, this.rooms});

  factory RoomsResponse.fromJson(Map<String, dynamic> json) {
    return RoomsResponse(
      result: json['result'] as String,
      message: json['message'] as String,
      rooms: (json['rooms'] as List<dynamic>).map((i) => Room.fromJson(i)).toList(),
    );
  }
}

class Room {
  final String roomId;
  final RoomBody roomBody;

  Room({this.roomId, this.roomBody});

  factory Room.fromJson(Map<String, dynamic> json) {
    return Room(
      roomId: json['room_id'] as String,
      roomBody: RoomBody.fromJson(json['room_body']),
    );
  }
}

class RoomBody {
  final DateTime created;
  final String name;
  final String type;
  final List<dynamic> users; // List<String>だとエラーなる

  RoomBody({this.created, this.name, this.type, this.users});

  factory RoomBody.fromJson(Map<String, dynamic> json) {
    return RoomBody(
      created: DateTime.parse(json['created'] as String),
      name: json['name'] as String,
      type: json['type'] as String,
      users: json['users'] as List<dynamic>,
    );
  }
}

class EnterRoomResponse {
  final String result;
  final String message;

  EnterRoomResponse({this.result, this.message});

  factory EnterRoomResponse.fromJson(Map<String, dynamic> json) {
    return EnterRoomResponse(
      result: json['result'] as String,
      message: json['message'] as String,
    );
  }
}
