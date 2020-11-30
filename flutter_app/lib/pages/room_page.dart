import 'dart:convert';

import 'package:firebase_auth/firebase_auth.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_app/controllers/api_links.dart';
import 'package:flutter_app/controllers/custom_dialog.dart';
import 'package:flutter_app/controllers/shared_preferences.dart';
import 'package:http/http.dart' as http;

import '../home_page.dart';
import '../main.dart';
import 'in_room.dart';

class RoomPage extends StatefulWidget {
  static const String pageTitle = 'ルーム一覧';

  @override
  _RoomPageState createState() => _RoomPageState();
}

class _RoomPageState extends State<RoomPage> {
  Future<List<Room>> _futureList;
  SharedPrefs _prefs;

  DateTime _lastLoaded = DateTime.now().toLocal();

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

  Future<List<Room>> fetchRooms() async {
    print('fetchRooms()');
    Uri uri = Uri.https(ApiLinks.Authority, ApiLinks.Rooms);
    final response = await http.get(uri);
    if (response.statusCode == 200) {
      RoomsResponse roomsResponse = RoomsResponse.fromJson(
          json.decode(utf8.decode(response.bodyBytes)));
      if (roomsResponse.result == 'ok') {
        _lastLoaded = DateTime.now().toLocal();
        return roomsResponse.rooms;
      } else {
        CustomDialog.showAlertDialog(context,
          '読み込みに失敗しました。\n' + roomsResponse.message,
          onOkPressed: () {
            Navigator.popUntil(context, ModalRoute.withName(MyHomePage.routeName));
          },
        );
      }
    } else {
      CustomDialog.showAlertDialog(context,
        '通信に失敗しました。',
        onOkPressed: () {
          Navigator.popUntil(context, ModalRoute.withName(MyHomePage.routeName));
        },
      );
    }
  }

  Future<void> enterRoom(BuildContext context, Room roomInfo) async {
    _prefs = await generateSharedPrefs();

    final _body = json.encode({
      'room_id': roomInfo.roomId,
      'user_id': FirebaseAuth.instance.currentUser.uid,
      'id_token': await FirebaseAuth.instance.currentUser.getIdToken(),
    });
    Uri uri = Uri.https(ApiLinks.Authority, ApiLinks.EnterRoom);

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

  void showEnterRoomDialog(BuildContext context, Room roomInfo) {
    CustomDialog.showAlertDialog(context,
      roomInfo.roomBody.name + 'の部屋に入りますか？',
      onCancelPressed: () {
        Navigator.pop(context);
      },
      onOkPressed: () {
        CustomDialog.showLoadingDialog(context, title: '入室中');
        enterRoom(context, roomInfo);
      }
    );
  }

  @override
  void initState() {
    _init();
    super.initState();
  }


  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: RefreshIndicator(
          onRefresh: () {
            return _onRefresh();
          },
          child: FutureBuilder<List<Room>>(
            future: _futureList,
            builder: (context, snapshot) {
              if (snapshot.hasData) {
                final List<Room> rooms = snapshot.data;
                return Column(
                  children: [
                    Padding(
                      padding: EdgeInsets.all(7.0),
                      child: Align(
                        alignment: Alignment.centerRight,
                        child: Opacity(
                          opacity: 0.5,
                          child: Text(
                            '最終更新：'
                                + _lastLoaded.hour.toString() + '時 '
                                + _lastLoaded.minute.toString() + '分',
                            style: TextStyle(),
                          ),
                        ),
                      ),
                    ),
                    Divider(),
                    Flexible(
                      child: ListView.builder(
                        shrinkWrap: true,
                        padding: const EdgeInsets.all(8),
                        itemCount: rooms.length,
                        itemBuilder: (BuildContext context, int index) {
                          return Card(
                            child: ListTile(
                              title: Text(
                                rooms[index].roomBody.name,
                                style: TextStyle(
                                  fontSize: 20,
                                ),
                              ),
                              subtitle: Text(
                                rooms[index].roomBody.users.length.toString() + '人',
                                textAlign: TextAlign.right,
                              ),
                              onTap: () {
                                showEnterRoomDialog(context, rooms[index]);
                              },
                            ),
                          );
                        },
                        // separatorBuilder: (BuildContext context, int index) => const Divider(),
                      ),
                    ),
                  ],
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


class RoomsResponse {
  final String result;
  final String message;
  final List<Room> rooms;

  RoomsResponse({this.result, this.message, this.rooms});

  factory RoomsResponse.fromJson(Map<String, dynamic> json) {
    return RoomsResponse(
      result: json['result'] ?? '',
      message: json['message'] ?? '',
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
      roomId: json['room_id'] ?? '',
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
      name: json['name'] ?? '',
      type: json['type'] ?? '',
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
      result: json['result'] ?? '',
      message: json['message'] ?? '',
    );
  }
}
