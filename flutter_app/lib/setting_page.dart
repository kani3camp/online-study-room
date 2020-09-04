import 'dart:convert';

import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:http/http.dart' as http;

class SettingPage extends StatefulWidget {
  @override
  State<StatefulWidget> createState() => SettingPageState();
}

class SettingPageState extends State<SettingPage> {
  SharedPreferences _prefs;

  static const DISPLAY_NAME = 'display-name';
  static const QUICK_WORD = 'quick-word';
  static const ACCOUNT_TYPE = 'account-type';
  static const MAIL_ADDRESS = 'mail-address';
  static const SUM_STUDY_TIME = 'sum-study-time';
  static const REGISTRATION_DATE = 'registration-date';

  String _displayName = '';
  String _quickWord = '';
  String _accountType = '';
  String _mailAddress = '';
  Duration _sumStudyTime = new Duration();
  DateTime _registrationDate = DateTime.now();


  @override
  void initState() {
    _initPreferences();
    super.initState();
    _fetchPreferences();
  }

  Future<void> _initPreferences() async {
    _prefs = await SharedPreferences.getInstance();
    setState(() {
      _displayName = _prefs.getString(DISPLAY_NAME) ?? '';
      _quickWord = _prefs.getString(QUICK_WORD) ?? '';
      _accountType = _prefs.getString(ACCOUNT_TYPE) ?? '';
      _mailAddress = _prefs.getString(MAIL_ADDRESS) ?? '';
      _sumStudyTime = new Duration(minutes: _prefs.getInt(SUM_STUDY_TIME) ?? 0);
      _registrationDate = DateTime.fromMillisecondsSinceEpoch(_prefs.getInt(REGISTRATION_DATE) ?? 0);
    });
  }

  Future<void> _fetchPreferences() async {
    Map<String, String> queryParams = {
      'user_id': 'test01'
    };
    Uri uri = Uri.https('us-central1-online-study-room-f1f30.cloudfunctions.net', '/UserStatus', queryParams);
    final response = await http.get(uri);
    if (response.statusCode == 200) {
      UserStatusResponse userStatusResp = UserStatusResponse.fromJson(json.decode(utf8.decode(response.bodyBytes)));
      if (userStatusResp.result == 'ok') {
        UserBody user = userStatusResp.userStatus.userBody;
        _prefs.setString(DISPLAY_NAME, user.name);
        _prefs.setString(QUICK_WORD, user.status);
        // todo _prefs.setInt(SUM_STUDY_TIME, )
        _prefs.setInt(REGISTRATION_DATE, user.registrationDate.millisecondsSinceEpoch);

        await _initPreferences();
      } else {
        throw Exception('Failed to load user status: ' + userStatusResp.message);
      }
    } else {
      throw Exception('http request failed');
    }

  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Center(
            child: Text('設定')
        ),
      ),
      body: ListView(
        children: [
          Container(
            child: Row(
              children: [
                Text('表示名：'),
                Text(_displayName),
              ],
            ),
          ),
          Container(
            child: Row(
              children: [
                Text('ひとこと：'),
                Text(_quickWord),
              ],
            ),
          ),
          Divider(),
          Container(
            child: Row(
              children: [
                Text('ログイン中のアカウント：'),
                Text(_accountType),
              ],
            ),
          ),
          Container(
            child: Row(
              children: [
                Text('メールアドレス：'),
                Text(_mailAddress),
              ],
            ),
          ),
          Container(
            child: Row(
              children: [
                Text('合計学習時間：'),
                Text(_sumStudyTime.toString() + '分'),
              ],
            ),
          ),
          Container(
            child: Row(
              children: [
                Text('登録日：'),
                Text(_registrationDate.toString()), // todo
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class UserStatusResponse {
  final String result;
  final String message;
  final UserStatus userStatus;

  UserStatusResponse({this.result, this.message, this.userStatus});

  factory UserStatusResponse.fromJson(Map<String, dynamic> json) {
    return UserStatusResponse(
        result: json['result'] as String,
        message: json['message'] as String,
        userStatus: UserStatus.fromJson(json['user_status'])
    );
  }
}

class UserStatus {
  final String userId;
  final UserBody userBody;

  UserStatus({this.userId, this.userBody});

  factory UserStatus.fromJson(Map<String, dynamic> json) {
    return UserStatus(
      userId: json['user_id'] as String,
      userBody: UserBody.fromJson(json['user_body'])
    );
  }
}

class UserBody {
  final String inRoom;
  final DateTime lastAccess;
  final DateTime lastEntered;
  final DateTime lastExited;
  final DateTime lastStudied;
  final String name;
  final bool online;
  final String status;
  final DateTime registrationDate;
  
  UserBody({this.inRoom, this.lastAccess, this.lastEntered, this.lastExited, this.lastStudied, this.name, this.online, this.status, this.registrationDate});
  
  factory UserBody.fromJson(Map<String, dynamic> json) {
    return UserBody(
      inRoom: json['in'] as String,
      lastAccess: DateTime.parse(json['last_access']),
      lastEntered: DateTime.parse(json['last_entered']),
      lastExited: DateTime.parse(json['last_exited']),
      lastStudied: DateTime.parse(json['last_studied']),
      name: json['name'] as String,
      online: json['online'] as bool,
      status: json['status'] as String,
      registrationDate: DateTime.parse(json['registration_date'])
    );
  }
}