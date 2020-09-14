import 'dart:convert';

import 'package:firebase_auth/firebase_auth.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_app/shared_preferences.dart';
import 'package:google_sign_in/google_sign_in.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:http/http.dart' as http;

class SettingPage extends StatefulWidget {
  @override
  State<StatefulWidget> createState() => SettingPageState();
}

class SettingPageState extends State<SettingPage> {
  SharedPrefs _prefs;

  String _displayName = '';
  String _quickWord = '';
  String _accountType = '';
  String _mailAddress = '';
  Duration _sumStudyTime = new Duration();
  DateTime _registrationDate = DateTime.now();

  final _displayNameController = TextEditingController();
  final _quickWordController = TextEditingController();

  bool _isButtonDisabled = true;

  @override
  void initState() {
    _initPreferences();
    _fetchPreferences();
    super.initState();
  }

  Future<void> _initPreferences() async {
    _prefs = await SharedPrefs.create();
    _displayName = await _prefs.getDisplayName();
    _quickWord = await _prefs.getQuickWord();
    _accountType = await _prefs.getAccountType();
    _mailAddress = await _prefs.getMailAddress();
    _sumStudyTime = _prefs.getSumStudyTime();
    _registrationDate = _prefs.getRegistrationDate();

    _displayNameController.text = _displayName;
    _quickWordController.text = _quickWord;

    _displayNameController.addListener(updateButtonState);
    _quickWordController.addListener(updateButtonState);

    setState(() {});
  }

  Future<void> _fetchPreferences() async {
    if (_prefs == null) {
      _prefs = await SharedPrefs.create();
    }
    String userId = await _prefs.getUserId();
    Map<String, String> queryParams = {
      'user_id': await _prefs.getUserId()
    };
    Uri uri = Uri.https('us-central1-online-study-room-f1f30.cloudfunctions.net', '/UserStatus', queryParams);
    final response = await http.get(uri);
    if (response.statusCode == 200) {
      UserStatusResponse userStatusResp = UserStatusResponse.fromJson(json.decode(utf8.decode(response.bodyBytes)));
      if (userStatusResp.result == 'ok') {
        UserBody user = userStatusResp.userStatus.userBody;
        await _prefs.setDisplayName(user.name);
        await _prefs.setQuickWord(user.status);
        // await _prefs.setSumStudyTime(user.); todo
        await _prefs.setRegistrationDate(user.registrationDate);

        await _initPreferences();
      } else {
        throw Exception('Failed to load user status: ' + userStatusResp.message);
      }
    } else {
      throw Exception('http request failed');
    }
  }

  void updateButtonState() {
    setState(() {
      _isButtonDisabled =
          _displayName == _displayNameController.text
          && _quickWord == _quickWordController.text;
    });
  }

  Future<void> saveNewValues() async {
    setState(() {
      _isButtonDisabled = true;
    });
    final _body = {
      'display_name': _displayNameController.text,
      'status_message': _quickWordController.text,
      'user_id': await _prefs.getUserId(),
      'id_token': await FirebaseAuth.instance.currentUser.getIdToken(),
    };
    Uri uri = Uri.https('us-central1-online-study-room-f1f30.cloudfunctions.net', '/ChangeUserInfo');

    final response = await http.post(
        uri,
        body: _body);
    if (response.statusCode == 200) {
      ChangeUserInfoResponse changeUserInfoResp = ChangeUserInfoResponse.fromJson(json.decode(utf8.decode(response.bodyBytes)));
      if (changeUserInfoResp.result == 'ok') {
        await _prefs.setDisplayName(_displayNameController.text);
        await _prefs.setQuickWord(_quickWordController.text);

        await _initPreferences();
      } else {
        throw Exception('Failed to change user info: ' + changeUserInfoResp.message);
      }
    } else {
      setState(() {
        _isButtonDisabled = false;
      });
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
          Text('表示名：'),
          ListTile(
            title: TextField(
              controller: _displayNameController,
              decoration: InputDecoration(
                hintText: '表示名'
              ),
            ),
          ),
          Text('ひとこと：'),
          ListTile(
            title: TextField(
              controller: _quickWordController,
              decoration: InputDecoration(
                hintText: 'ひとこと'
              ),
            ),
          ),
          Divider(),
          Text('ログイン中のアカウントの種類：'),
          ListTile(
            title: Text(_accountType),
          ),
          Text('メールアドレス：'),
          ListTile(
            title: Text(_mailAddress),
          ),
          Text('合計学習時間：'),
          ListTile(
            title: Text(_sumStudyTime.toString() + '分'),
          ),
          Text('登録日：'),
          ListTile(
            title: Text(_registrationDate.toString()),
          ),
          Container(
            child: RaisedButton(
              child: Text('保存'),
              onPressed: _isButtonDisabled ? null : saveNewValues,
            ),
          ),
          Container(
            child: RaisedButton(
              child: Text('ログアウト'),
              onPressed: () {
                setState(() {
                  FirebaseAuth.instance.signOut();
                  GoogleSignIn().signOut();
                  Navigator.of(context).pushReplacementNamed('/login');
                });
              },
            ),
          )
        ],
      ),
    );
  }

  @override
  void dispose() {
    // Clean up the controller when the widget is removed from the
    // widget tree.
    _displayNameController.dispose();
    super.dispose();
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

class ChangeUserInfoResponse {
  final String result;
  final String message;

  ChangeUserInfoResponse({this.result, this.message});

  factory ChangeUserInfoResponse.fromJson(Map<String, dynamic> json) {
    return ChangeUserInfoResponse(
        result: json['result'] as String,
        message: json['message'] as String,
    );
  }
}
