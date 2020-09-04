import 'dart:convert';

import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

class RoomPage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    Future<List> futureList = fetchRooms();

    return Scaffold(
      appBar: AppBar(
        title: Center(
            child: Text('カテゴリ一覧')
        ),
      ),
      body: FutureBuilder<List>(
        future: futureList,
        builder: (context, snapshot) {
          if (snapshot.hasData) {
            return ListView.separated(
              padding: const EdgeInsets.all(8),
              itemCount: snapshot.data.length,
              itemBuilder: (BuildContext context, int index) {
                return Container(
                  child: ListTile(
                    title: Text(snapshot.data[index].roomBody.name),
                    onTap: () {},
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
      ),
    );
  }
}

Future<List> fetchRooms() async {
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
  final String users;

  RoomBody({this.created, this.name, this.users});

  factory RoomBody.fromJson(Map<String, dynamic> json) {
    return RoomBody(
      created: DateTime.parse(json['created'] as String),
      name: json['name'] as String,
      users: json['users'] as String,
    );
  }
}