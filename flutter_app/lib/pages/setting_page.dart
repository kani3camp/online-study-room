import 'dart:convert';

import 'package:firebase_auth/firebase_auth.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_app/controllers/api_links.dart';
import 'package:flutter_app/controllers/custom_dialog.dart';
import 'package:flutter_app/controllers/shared_preferences.dart';
import 'package:flutter_app/login_page.dart';
import 'package:google_sign_in/google_sign_in.dart';
import 'package:http/http.dart' as http;

import '../home_page.dart';

class SettingPage extends StatefulWidget {
  static const String pageTitle = '設定';

  @override
  State<StatefulWidget> createState() => SettingPageState();
}

class SettingPageState extends State<SettingPage> {
  SharedPrefs _prefs;

  String _displayName = '';
  String _quickWord = '';
  String _accountType = '';
  String _mailAddress = '';
  int _totalStudyTime = 0;
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
    _prefs = await generateSharedPrefs();
    if (FirebaseAuth.instance?.currentUser != null) {
      _displayName = FirebaseAuth.instance.currentUser.displayName;
      _quickWord = await _prefs.getQuickWord();
      _accountType =
          FirebaseAuth.instance.currentUser.providerData[0].providerId;
      _mailAddress = FirebaseAuth.instance.currentUser.email;
      _totalStudyTime = _prefs.getTotalStudyTime();
      _registrationDate = _prefs.getRegistrationDate();

      _displayNameController.text = _displayName;
      _quickWordController.text = _quickWord;

      _displayNameController.addListener(updateButtonState);
      _quickWordController.addListener(updateButtonState);
    } else {
      GoogleSignIn().signOut();
      Navigator.of(context).pushReplacementNamed(LoginPage.routeName);
    }

    if (this.mounted) {
      setState(() {});
    }
  }

  Future<void> _fetchPreferences() async {
    if (_prefs == null) {
      _prefs = new SharedPrefs();
      await _prefs.init();
    }
    Map<String, String> queryParams = {
      'user_id': FirebaseAuth.instance.currentUser.uid
    };
    Uri uri = Uri.https(ApiLinks.Authority, ApiLinks.UserStatus, queryParams);
    final response = await http.get(uri);
    if (response.statusCode == 200) {
      UserStatusResponse userStatusResp = UserStatusResponse.fromJson(json.decode(utf8.decode(response.bodyBytes)));
      if (userStatusResp.result == 'ok') {
        UserBody user = userStatusResp.userStatus.userBody;
        await _prefs.setQuickWord(user.status);
        await _prefs.setTotalStudyTime(user.totalStudyTime);
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
              && _quickWord == _quickWordController.text
              || _displayNameController.text == ''
              || _quickWordController.text == '';
    });
  }

  Future<void> saveNewValues() async {
    CustomDialog.showLoadingDialog(context, title: '保存中');
    setState(() {
      _isButtonDisabled = true;
    });
    final _body = json.encode({
      'display_name': _displayNameController.text,
      'status_message': _quickWordController.text,
      'user_id': FirebaseAuth.instance.currentUser.uid,
      'id_token': await FirebaseAuth.instance.currentUser.getIdToken(),
    });
    Uri uri = Uri.https(ApiLinks.Authority, ApiLinks.ChangeUserInfo);

    final response = await http.post(
        uri,
        body: _body);
    if (response.statusCode == 200) {
      ChangeUserInfoResponse changeUserInfoResp = ChangeUserInfoResponse.fromJson(json.decode(utf8.decode(response.bodyBytes)));
      if (changeUserInfoResp.result == 'ok') {
        CustomDialog.showAlertDialog(
          context,
          '設定が完了しました。',
          onOkPressed: () {
            Navigator.pop(context);
          },
        );

        await _prefs.setQuickWord(_quickWordController.text);
        await _initPreferences();
      } else {
        CustomDialog.showAlertDialog(context,
            '問題が発生しました。\n' + changeUserInfoResp.message,
            onOkPressed: () {
              Navigator.popUntil(context, ModalRoute.withName(MyHomePage.routeName));
            }
        );
        await _initPreferences();
        setState(() {
          _isButtonDisabled = false;
        });
      }
    } else {
      CustomDialog.showAlertDialog(context,
          '通信が失敗しました。\n',
          onOkPressed: () {
            Navigator.popUntil(context, ModalRoute.withName(MyHomePage.routeName));
          }
      );
      await _initPreferences();
      setState(() {
        _isButtonDisabled = false;
      });
    }
    Navigator.popUntil(context, ModalRoute.withName(MyHomePage.routeName));
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
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
            title: Text(
              (_totalStudyTime ~/ 3600).toString() + '時間' + ((_totalStudyTime % 3600) ~/ 60).toString() + '分',
            ),
          ),
          Text('登録日：'),
          ListTile(
            title: Text(_registrationDate.toString()),
          ),
          Container(
            padding: EdgeInsets.all(10.0),
            child: RaisedButton(
              color: Colors.blueAccent,
              child: Text('保存'),
              onPressed: _isButtonDisabled ? null : saveNewValues,
            ),
          ),
          Container(
            padding: EdgeInsets.all(10.0),
            child: RaisedButton(
              color: Colors.yellow,
              child: Text('ログアウト'),
              onPressed: () {
                setState(() {
                  FirebaseAuth.instance.signOut();
                  GoogleSignIn().signOut();
                  Navigator.of(context).pushReplacementNamed(LoginPage.routeName);
                });
              },
            ),
          )
        ],
      ),
    );
  }

  // 画面遷移するだけでもdisposeされるため、結局毎回リロードしなければならなくなるためなくて良い
  // @override
  // void dispose() {
  //   // Clean up the controller when the widget is removed from the
  //   // widget tree.
  //   // _displayNameController.dispose();
  //   // _quickWordController.dispose();
  //   super.dispose();
  // }
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
  final String displayName;
  final UserBody userBody;

  UserStatus({
    this.userId,
    this.displayName,
    this.userBody
  });

  factory UserStatus.fromJson(Map<String, dynamic> json) {
    return UserStatus(
      userId: json['user_id'] as String,
      displayName: json['display_name'] as String,
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
  final bool online;
  final String status;
  final DateTime registrationDate;
  final int totalStudyTime;
  final int totalBreakTime;
  
  UserBody({
    this.inRoom,
    this.lastAccess,
    this.lastEntered,
    this.lastExited,
    this.lastStudied,
    this.online,
    this.status,
    this.registrationDate,
    this.totalStudyTime,
    this.totalBreakTime,
  });
  
  factory UserBody.fromJson(Map<String, dynamic> json) {
    return UserBody(
      inRoom: json['in'] as String,
      lastAccess: DateTime.parse(json['last_access']).toLocal(),
      lastEntered: DateTime.parse(json['last_entered']).toLocal(),
      lastExited: DateTime.parse(json['last_exited']).toLocal(),
      lastStudied: DateTime.parse(json['last_studied']).toLocal(),
      online: json['online'] as bool,
      status: json['status'] as String,
      registrationDate: DateTime.parse(json['registration_date']),
      totalStudyTime: json['total_study_time'] as int,
      totalBreakTime: json['total_break_time'] as int,
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
